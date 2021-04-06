package parser

type Exception interface{}

type TryBlock struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type TemplateArg struct {
	Position int    `json:"position"`
	Name     string `json:"name"`
}

type FuncArg struct {
	Position int     `json:"position"`
	Name     *string `json:"name"`
	Type     string  `json:"type"`
}

type TemplatedFunc struct {
	TemplateArgs []*TemplateArg `json:"template_args"`
	ReturnArgs   []*FuncArg     `json:"return_args"`
	FuncArgs     []*FuncArg     `json:"func_args"`
	Name         string         `json:"name"`
	Body         string         `json:"body"`
}

type Caller struct {
	TemplateReplacementArgs []*FuncArg     `json:"func_args"`
	LinkedTemplate          *TemplatedFunc `json:"linked_template"`
	Generator               string         `json:"generator"`
}

func (t *TemplatedFunc) GetTemplateArgAtPos(pos int) *TemplateArg {
	for _, arg := range t.TemplateArgs {
		if arg.Position == pos {
			return arg
		}
	}

	return nil
}

func (t *TemplatedFunc) GetTemplateArgByType(argType string) *TemplateArg {
	for _, arg := range t.TemplateArgs {
		if arg.Name == argType {
			return arg
		}
	}

	return nil
}

func (t *TemplatedFunc) DoesTemplateArgExist(input string) bool {
	for _, arg := range t.TemplateArgs {
		if arg.Name == input {
			return true
		}
	}

	return false
}

func (t *TemplatedFunc) GetFuncArgAtPos(pos int) *FuncArg {
	for _, arg := range t.FuncArgs {
		if arg.Position == pos {
			return arg
		}
	}

	return nil
}

func (t *TemplatedFunc) GetReturnArgAtPos(pos int) *FuncArg {
	for _, arg := range t.ReturnArgs {
		if arg.Position == pos {
			return arg
		}
	}

	return nil
}

func (tb TryBlock) Run() {
	if tb.Finally != nil {
		defer tb.Finally()
	}

	if tb.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tb.Catch(r)
			}
		}()
	}

	tb.Try()
}