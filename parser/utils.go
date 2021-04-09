package parser

import (
	"fmt"
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

	for _, arg := range caller.LinkedTemplate.ReturnArgs {
		callerArgs = append(callerArgs, caller.TemplateReplacementArgs[arg.Position].Type)
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
		if sl0[i] != sl1[i] {
			return false
		}
	}

	return true
}
