package parser

const (
	TEMPLATE_BODY = `@template([A-z0-9\(*\)\s].+)\nfunc\s*([A-z0-9\(*\)\s]+)?\s*\((?:[^)(]+|\((?:[^)(]+|\([^)(]*\))*\))*\)\s*[\(\)a-zA-Z]?.*\s?\{(?:[^}{]+|\{(?:[^}{]+|\{[^}{]*\})*\})*\}`
	CALLER_BODY   = `([a-zA-Z0-9]*.<.*[a-zA-Z0-9]>\(.*[a-zA-Z0-9]*\))`
)
