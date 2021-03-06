package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func ParallelParse(inputPath string, outputPath string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Processing %s -> %s\n", inputPath, outputPath)

	if err := Parse(inputPath, outputPath, false); err != nil {
		panic(err)
	}
}

func Parse(inputFile, outputFile string, execute bool) error {
	test, err := os.Stat(inputFile)
	if err != nil {
		return err
	}

	if test.IsDir() {
		parsedDirName := "parsed"

		files, err := ioutil.ReadDir(inputFile)
		if err != nil {
			return err
		}

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		outputPath := filepath.Join(wd, inputFile, parsedDirName)
		os.Mkdir(outputPath, 0755)

		group := sync.WaitGroup{}
		group.Add(len(files))

		for _, f := range files {
			input := wd + filepath.FromSlash(fmt.Sprintf("/%s/%s", inputFile, f.Name()))
			output := wd + filepath.FromSlash(fmt.Sprintf("/%s/%s/%s.go", inputFile, parsedDirName, f.Name()))

			go ParallelParse(input, output, &group)
		}

		group.Wait()

		return nil
	}

	input, err := readFile(inputFile)
	if err != nil {
		return err
	}

	tmpls, cleaned, err := parseTemplates(input)
	if err != nil {
		return err
	}

	out, err := buildFuncData(cleaned, tmpls)
	if err != nil {
		return err
	}

	if err := writeFile(outputFile, out); err != nil {
		return err
	}

	if execute {
		cmd := exec.Command("go", "run", outputFile)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func buildFuncData(data string, tmpls *Templates) ([]string, error) {
	pattern, err := regexp.Compile(CALLER_BODY)
	if err != nil {
		return nil, err
	}

	parseCaller := func(full string) (*Caller, error) {
		full = strings.TrimSpace(full)

		caller := Caller{
			Generator: NewGenerator(),
		}

		for _, tmpl := range tmpls.Funcs {
			if strings.HasPrefix(full, tmpl.Name) {
				definedType := strings.TrimSpace(strings.Split(full, tmpl.Name)[1])
				if strings.HasPrefix(definedType, "<") {
					funcCallTypes := strings.TrimSpace(strings.Split(strings.Split(full, "<")[1], ">")[0])
					funcCallTypeDefinitions := strings.Split(funcCallTypes, ",")

					if len(tmpl.TemplateArgs) == len(funcCallTypeDefinitions) {
						caller.LinkedTemplate = tmpl

						funcCallArgs := strings.TrimSpace(strings.Split(strings.Split(full, "(")[1], ")")[0])
						funcCallArgDefinitions := strings.Split(funcCallArgs, ",")

						for i, call := range funcCallArgDefinitions {
							funcArg := FuncArg{
								Position: i,
								Name:     &call,
							}

							if len(funcCallTypeDefinitions) == 1 {
								funcArg.Type = funcCallTypeDefinitions[0]
							} else {
								funcArg.Type = funcCallTypeDefinitions[i]
							}

							caller.TemplateReplacementArgs = append(caller.TemplateReplacementArgs, &funcArg)
						}
					}
				}
			}
		}

		return &caller, nil
	}

	var callers []*Caller
	lines := strings.Split(strings.TrimSpace(data), "\n")

	if matching := pattern.MatchString(data); matching {
		matches := pattern.FindAllString(data, -1)
		if len(matches) > 0 {
			for _, match := range matches {
				parsed, err := parseCaller(match)
				if err != nil {
					return nil, err
				}

				match = strings.TrimSpace(match)
				newName := fmt.Sprintf("%s_gonerics_%s", parsed.LinkedTemplate.Name, parsed.Generator)

				for idx, line := range lines {
					line = strings.TrimSpace(line)
					if IsMatchingGenericCall(line, parsed) || match == line {
						lines[idx] = strings.ReplaceAll(lines[idx], parsed.LinkedTemplate.Name, newName)
						lines[idx] = ReplaceGenericCalls(lines[idx])
					}
				}

				var params []string
				for _, param := range parsed.LinkedTemplate.FuncArgs {
					name := strings.TrimSpace(*param.Name)

					var namedType string
					if param.IsBuiltIn {
						namedType = strings.TrimSpace(param.Type)
					} else {
						namedType = strings.TrimSpace(parsed.TemplateReplacementArgs[param.Position].Type)
					}

					params = append(params, fmt.Sprintf("%s %s", name, namedType))
				}

				var args []string
				for _, arg := range parsed.LinkedTemplate.ReturnArgs {
					replace := "%s %s"

					var argType string
					if arg.IsBuiltIn {
						argType = strings.TrimSpace(arg.Type)
					} else {
						argType = strings.TrimSpace(parsed.TemplateReplacementArgs[arg.Position].Type)
					}

					var newArg string
					if arg.Name != nil {
						newArg = strings.TrimSpace(fmt.Sprintf(replace, *arg.Name, argType))
					} else {
						newArg = strings.TrimSpace(fmt.Sprintf(replace, "", argType))
					}

					args = append(args, newArg)
				}

				splat := strings.Join(args, ",")
				if len(splat) > 1 {
					splat = fmt.Sprintf("(%s)", splat)
				}

				gen := fmt.Sprintf(`
					func %s(%s) %s {
						%s
					}
				`, newName, strings.Join(params, ","), splat, parsed.LinkedTemplate.Body)

				lines = append(lines, strings.TrimSpace(gen))
				callers = append(callers, parsed)
			}
		}
	}

	debugging := true
	if debugging {
		templatesOutput, err := json.MarshalIndent(tmpls, "", " ")
		if err != nil {
			return nil, err
		}

		if err := writeJSON("DEBUG/templates.json", templatesOutput); err != nil {
			return nil, err
		}

		callersOutput, err := json.MarshalIndent(callers, "", " ")
		if err != nil {
			return nil, err
		}

		if err := writeJSON("DEBUG/callers.json", callersOutput); err != nil {
			return nil, err
		}
	}

	return lines, nil
}

func parseTemplates(data []byte) (*Templates, string, error) {
	var funcTemplates []*TemplatedFunc
	var structTemplates []*TemplatedStruct

	str := string(data)

	curlyPattern, err := regexp.Compile(BETWEEN_CURLYS)
	if err != nil {
		return nil, str, err
	}

	funcPattern, err := regexp.Compile(TEMPLATE_BODY)
	if err != nil {
		return nil, str, err
	}

	structPattern, err := regexp.Compile(TEMPLATE_STRUCT)
	if err != nil {
		return nil, str, err
	}

	structFieldPattern, err := regexp.Compile(STRUCT_FIELD)
	if err != nil {
		return nil, str, err
	}

	parseTemplateFunc := func(full string) (*TemplatedFunc, error) {
		tmpl := TemplatedFunc{}
		lines := strings.Split(strings.TrimSpace(full), "\n")
		for idx, line := range lines {
			if strings.HasPrefix(line, "@template") {
				templateLineArgDefinitionString := strings.TrimSpace(strings.Split(strings.Split(line, "(")[1], ")")[0])
				templateLineArgDefinitions := strings.Split(templateLineArgDefinitionString, ",")

				for i, tmplArg := range templateLineArgDefinitions {
					tmplArg = strings.TrimSpace(strings.ReplaceAll(tmplArg, "type", ""))
					tmpl.TemplateArgs = append(tmpl.TemplateArgs, &TemplateArg{
						Position: i,
						Name:     strings.TrimSpace(tmplArg),
					})
				}

				if strings.HasPrefix(lines[idx+1], "func") {
					funcLine := lines[idx+1]

					tmpl.Name = strings.TrimSpace(strings.Split(strings.Split(funcLine, "func")[1], "(")[0])
					funcLineArgDefinitionString := strings.TrimSpace(strings.Split(strings.Split(funcLine, "(")[1], ")")[0])
					funcLineArgDefinitions := strings.Split(funcLineArgDefinitionString, ",")

					for i, funcArg := range funcLineArgDefinitions {
						funcArg := strings.TrimSpace(funcArg)

						tmplType := strings.TrimSpace(strings.Split(funcArg, " ")[1])
						tmplName := strings.TrimSpace(strings.Split(funcArg, " ")[0])
						isGeneric, tmplArg := tmpl.GetTemplateArgByType(tmplType)

						if isGeneric && tmplArg != nil {
							if tmplArg.Name == tmplType {
								tmpl.FuncArgs = append(tmpl.FuncArgs, &FuncArg{
									Name:      &tmplName,
									Type:      tmplType,
									Position:  tmplArg.Position,
									IsBuiltIn: false,
								})
							}
						} else {
							tmpl.FuncArgs = append(tmpl.FuncArgs, &FuncArg{
								Name:      &tmplName,
								Type:      tmplType,
								Position:  i,
								IsBuiltIn: true,
							})
						}
					}

					TryBlock{
						Try: func() {
							funcLineReturnDefintionString := strings.TrimSpace(strings.Split(strings.Split(funcLine, funcLineArgDefinitionString)[1], "{")[0])

							var funcLineReturnArgDefinitions []string

							if len(strings.Split(strings.TrimSpace(strings.ReplaceAll(funcLineReturnDefintionString, ")", "")), ",")) == 1 {
								funcLineReturnArgDefinitions = []string{strings.TrimSpace(strings.ReplaceAll(funcLineReturnDefintionString, ")", ""))}
							} else {
								funcLineReturnDefinitionStringTrim := strings.TrimSpace(strings.Split(strings.Split(funcLineReturnDefintionString, "(")[1], ")")[0])
								funcLineReturnArgDefinitions = strings.Split(funcLineReturnDefinitionStringTrim, ",")
							}

							for i, funcRetArg := range funcLineReturnArgDefinitions {
								funcRetArg = strings.TrimSpace(funcRetArg)
								nameTypeEach := strings.Split(funcRetArg, " ")

								argType := FuncArg{
									IsBuiltIn: false,
								}

								if len(nameTypeEach) == 2 {
									isBuiltIn := true
									for _, funcArg := range tmpl.TemplateArgs {
										if funcArg.Name == nameTypeEach[1] {
											argType.Position = funcArg.Position
											argType.Type = strings.TrimSpace(nameTypeEach[1])
											if len(nameTypeEach[0]) > 0 {
												argType.Name = GetPointerToString(strings.TrimSpace(nameTypeEach[0]))
											}
											isBuiltIn = false
										}
									}

									if isBuiltIn {
										argType.IsBuiltIn = true
										argType.Type = nameTypeEach[1]
										argType.Name = GetPointerToString(strings.TrimSpace(nameTypeEach[0]))
										argType.Position = i
									}
								} else {
									isBuiltIn := true
									for _, funcArg := range tmpl.TemplateArgs {
										if funcArg.Name == nameTypeEach[0] {
											argType.Position = funcArg.Position
											argType.Type = nameTypeEach[0]
											isBuiltIn = false
										}
									}

									if isBuiltIn {
										argType.IsBuiltIn = true
										argType.Type = nameTypeEach[0]
										argType.Position = i
									}
								}

								tmpl.ReturnArgs = append(tmpl.ReturnArgs, &argType)
							}
						},
						Catch: func(e Exception) {
							fmt.Printf("func %s contains no return types. Error: %v\n", tmpl.Name, e)
						},
						Finally: func() {
							var body []string
							bodyIdx := idx + 2

							for {
								if bodyIdx+1 >= len(lines) {
									break
								}

								body = append(body, lines[bodyIdx])
								bodyIdx++
							}

							tmpl.Body = strings.Join(body, "\n")
						},
					}.Run()
				}
			}
		}

		return &tmpl, nil
	}

	parseTemplateStruct := func(full string) (*TemplatedStruct, error) {
		tmpl := TemplatedStruct{}
		lines := strings.Split(strings.TrimSpace(full), "\n")
		for idx, line := range lines {
			if strings.HasPrefix(line, "@template") {
				templateLineArgDefinitionString := strings.TrimSpace(strings.Split(strings.Split(line, "(")[1], ")")[0])
				templateLineArgDefinitions := strings.Split(templateLineArgDefinitionString, ",")

				for i, tmplArg := range templateLineArgDefinitions {
					tmplArg = strings.TrimSpace(strings.ReplaceAll(tmplArg, "type", ""))
					tmpl.TemplateArgs = append(tmpl.TemplateArgs, &TemplateArg{
						Position: i,
						Name:     strings.TrimSpace(tmplArg),
					})
				}

				peek := strings.TrimSpace(strings.ReplaceAll(lines[idx+1], "{", ""))
				if strings.HasPrefix(peek, "type") && strings.HasSuffix(peek, "struct") {
					structLine := lines[idx+1]
					tmpl.Name = strings.TrimSpace(strings.Split(strings.Split(structLine, "type")[1], "struct")[0])

					if matching := curlyPattern.MatchString(full); matching {
						match := curlyPattern.FindString(str)
						match = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(match, "{", ""), "}", ""))

						each := strings.Split(match, "\n")
						for i, def := range each {
							if def != "" {
								fieldArg, err := processStructField(i, def, &tmpl, structFieldPattern)
								if err != nil {
									fmt.Printf("could not process field arg from line %s. Error: %s\n", def, err.Error())
									continue
								}

								tmpl.FieldArgs = append(tmpl.FieldArgs, fieldArg)
							}
						}
					}
				}
			}
		}

		return &tmpl, nil
	}

	if matching := structPattern.MatchString(str); matching {
		matches := structPattern.FindAllString(str, -1)
		if len(matches) > 0 {
			for _, match := range matches {
				template, err := parseTemplateStruct(match)
				if err != nil {
					return nil, str, err
				}

				str = structPattern.ReplaceAllLiteralString(str, "")

				structTemplates = append(structTemplates, template)
			}
		}
	}

	if matching := funcPattern.MatchString(str); matching {
		matches := funcPattern.FindAllString(str, -1)
		if len(matches) > 0 {
			for _, match := range matches {
				template, err := parseTemplateFunc(match)
				if err != nil {
					return nil, str, err
				}

				str = funcPattern.ReplaceAllLiteralString(str, "")

				funcTemplates = append(funcTemplates, template)
			}
		}
	}

	return &Templates{
		Funcs:   funcTemplates,
		Structs: structTemplates,
	}, str, nil
}

func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func writeJSON(outputPath string, data []byte) error {
	return ioutil.WriteFile(outputPath, data, 0644)
}

func writeFile(outputPath string, data []string) error {
	file, err := os.OpenFile(outputPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	for _, l := range data {
		_, err := writer.WriteString(fmt.Sprintf("%s\n", l))
		if err != nil {
			return err
		}
	}

	writer.Flush()
	file.Close()

	return nil
}
