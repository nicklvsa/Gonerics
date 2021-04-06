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

	for _, arg := range caller.TemplateReplacementArgs {
		lineArg := strings.TrimSpace(lineArgs[arg.Position])
		argType := strings.TrimSpace(arg.Type)

		if lineArg == argType {
			return true
		}
	}

	return false
}

func ReplaceGenericCalls(line string) string {
	replacer := fmt.Sprintf("<%s>", strings.Split(strings.Split(line, strings.TrimSpace("<"))[1], ">")[0])
	return strings.ReplaceAll(line, replacer, "")
} 