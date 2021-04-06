package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func Parse(inputFile, outputFile string) error {
	test, err := os.Stat(inputFile)
	if err != nil {
		return err
	}

	if test.IsDir() {
		files, err := ioutil.ReadDir(inputFile)
		if err != nil {
			return err
		}

		for _, f := range files {
			Parse(fmt.Sprintf("%s/%s", inputFile, f.Name()), fmt.Sprintf("%s/%s/%s%s", inputFile, outputFile, f.Name(), ".go"))
		}

		return nil
	}

	input, err := readFile(inputFile)
	if err != nil {
		return err
	}

	templates, cleaned, err := parseTemplates(input)
	if err != nil {
		return err
	}

	out, err := buildFuncData(cleaned, templates)
	if err != nil {
		return err
	}

	if err := writeFile(outputFile, out); err != nil {
		return err
	}

	return nil
}

func buildFuncData(data string, tmpls []*TemplatedFunc) ([]string, error) {
	pattern, err := regexp.Compile(CALLER_BODY)
	if err != nil {
		return nil, err
	}

	parseCaller := func(full string) (*Caller, error) {
		full = strings.TrimSpace(full)

		caller := Caller{
			Generator: NewGenerator(),
		}

		for _, tmpl := range tmpls {
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
							caller.TemplateReplacementArgs = append(caller.TemplateReplacementArgs, &FuncArg{
								Position: i,
								Name:     &call,
								Type:     funcCallTypeDefinitions[i],
							})
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

					if line == match /*|| IsMatchingGenericCall(line, parsed)*/ {
						lines[idx] = strings.ReplaceAll(lines[idx], parsed.LinkedTemplate.Name, newName)
						lines[idx] = ReplaceGenericCalls(lines[idx])
					}
				}

				var params []string
				for _, param := range parsed.LinkedTemplate.FuncArgs {
					name := strings.TrimSpace(*param.Name)
					namedType := strings.TrimSpace(parsed.TemplateReplacementArgs[param.Position].Type)
					params = append(params, fmt.Sprintf("%s %s", name, namedType))
				}

				var args []string
				for _, arg := range parsed.LinkedTemplate.ReturnArgs {
					replace := "%s %s"

					argType := strings.TrimSpace(parsed.TemplateReplacementArgs[arg.Position].Type)

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

func parseTemplates(data []byte) ([]*TemplatedFunc, string, error) {
	var templates []*TemplatedFunc

	str := string(data)

	pattern, err := regexp.Compile(TEMPLATE_BODY)
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

					for _, funcArg := range funcLineArgDefinitions {
						funcArg := strings.TrimSpace(funcArg)

						tmplType := strings.TrimSpace(strings.Split(funcArg, " ")[1])
						tmplName := strings.TrimSpace(strings.Split(funcArg, " ")[0])
						tmplArg := tmpl.GetTemplateArgByType(tmplType)

						if tmplArg.Name == tmplType {
							tmpl.FuncArgs = append(tmpl.FuncArgs, &FuncArg{
								Name:     &tmplName,
								Type:     tmplType,
								Position: tmplArg.Position,
							})
						}
					}

					TryBlock{
						Try: func() {
							funcLineReturnDefintionString := strings.TrimSpace(strings.Split(strings.Split(funcLine, funcLineArgDefinitionString)[1], "{")[0])

							var funcLineReturnArgDefinitions []string

							if len(strings.TrimSpace(strings.ReplaceAll(funcLineReturnDefintionString, ")", ""))) == 1 {
								funcLineReturnArgDefinitions = []string{strings.TrimSpace(strings.ReplaceAll(funcLineReturnDefintionString, ")", ""))}
							} else {
								funcLineReturnDefinitionStringTrim := strings.TrimSpace(strings.Split(strings.Split(funcLineReturnDefintionString, "(")[1], ")")[0])
								funcLineReturnArgDefinitions = strings.Split(funcLineReturnDefinitionStringTrim, ",")
							}

							for _, funcRetArg := range funcLineReturnArgDefinitions {
								funcRetArg = strings.TrimSpace(funcRetArg)
								nameTypeEach := strings.Split(funcRetArg, " ")

								argType := FuncArg{}

								if len(nameTypeEach) == 2 {
									argType.Type = strings.TrimSpace(nameTypeEach[1])
									if len(nameTypeEach[0]) > 0 {
										argType.Name = GetPointerToString(strings.TrimSpace(nameTypeEach[0]))
									}

									for _, funcArg := range tmpl.TemplateArgs {
										if funcArg.Name == nameTypeEach[1] {
											argType.Position = funcArg.Position
										}
									}
								} else {
									for _, funcArg := range tmpl.TemplateArgs {
										if funcArg.Name == nameTypeEach[0] {
											argType.Position = funcArg.Position
											argType.Type = nameTypeEach[0]
										}
									}
								}

								if tmpl.DoesTemplateArgExist(argType.Type) {
									tmpl.ReturnArgs = append(tmpl.ReturnArgs, &argType)
								}
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

	if matching := pattern.MatchString(str); matching {
		matches := pattern.FindAllString(str, -1)
		if len(matches) > 0 {
			for _, match := range matches {
				template, err := parseTemplateFunc(match)
				if err != nil {
					return nil, str, err
				}

				str = pattern.ReplaceAllLiteralString(str, "")

				templates = append(templates, template)
			}
		}
	}

	return templates, str, nil
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
