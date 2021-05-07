package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/xid"
)

func GetPointerToString(s string) *string {
	return &s
}

func NewGenerator() string {
	return xid.New().String()
}

func IsMatchingGenericCall(compare string, caller *Caller) bool {
	var lineArgs []string

	TryBlock{
		Try: func() {
			lineArgs = strings.Split(strings.Split(strings.Split(strings.TrimSpace(compare), "<")[1], ">")[0], ",")
		},
		Catch: func(e Exception) {
			lineArgs = []string{}
		},
		Finally: nil,
	}.Run()

	if len(lineArgs) <= 0 || caller == nil {
		return false
	}

	var callerArgs []string

	for _, arg := range caller.LinkedTemplate.TemplateArgs {
		callerArgs = append(callerArgs, strings.TrimSpace(caller.TemplateReplacementArgs[arg.Position].Type))
	}

	return areSlicesEqual(lineArgs, callerArgs)
}

func ReplaceGenericCalls(line string) string {
	replacer := fmt.Sprintf("<%s>", strings.Split(strings.Split(line, strings.TrimSpace("<"))[1], ">")[0])
	return strings.ReplaceAll(line, replacer, "")
}

func areSlicesEqual(sl0, sl1 []string) bool {
	if (sl0 == nil) != (sl1 == nil) {
		return false
	}

	if len(sl0) != len(sl1) {
		return false
	}

	for i := range sl0 {
		if !strings.EqualFold(sl0[i], sl1[i]) {
			return false
		}
	}

	return true
}

func processStructField(idx int, field string, tmpl *TemplatedStruct, pattern *regexp.Regexp) (*FieldArg, error) {
	if pattern == nil {
		return nil, fmt.Errorf("provided field pattern was nil")
	}

	fieldArg := FieldArg{
		IsBuiltIn: true,
	}

	if matching := pattern.MatchString(field); matching {
		offset := 0
		tag := pattern.FindString(field)

		if tag != "" {
			field = field[:strings.Index(field, tag)]
			offset += 1
		}

		fieldLineArgs := strings.Split(strings.TrimSpace(field), " ")

		if len(fieldLineArgs)+offset == 3 {
			fieldArg.Name = strings.TrimSpace(fieldLineArgs[0])
			fieldArg.Type = strings.TrimSpace(fieldLineArgs[1])
			fieldArg.Tags = GetPointerToString(strings.TrimSpace(tag))
			fieldArg.Position = idx
		} else {
			fieldArg.Name = strings.TrimSpace(fieldLineArgs[0])
			fieldArg.Type = strings.TrimSpace(fieldLineArgs[1])
			fieldArg.Position = idx
		}
	} else {
		fieldLineArgs := strings.Split(strings.TrimSpace(field), " ")

		fieldArg.Name = strings.TrimSpace(fieldLineArgs[0])
		fieldArg.Type = strings.TrimSpace(fieldLineArgs[1])
		fieldArg.Position = idx
	}

	isGeneric, tmplArg := tmpl.GetTemplateArgByType(fieldArg.Type)
	if isGeneric && tmplArg != nil {
		fieldArg.IsBuiltIn = false
	}

	return &fieldArg, nil
}
