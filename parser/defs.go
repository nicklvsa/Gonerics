package parser

const (
	BETWEEN_CURLYS = `\{(?:[^}{]+|\{(?:[^}{]+|\{[^}{]*\})*\})*\}`
)

// func pattern matches
const (
	TEMPLATE_BODY = `@template([A-z0-9\(*\)\s].+)\nfunc\s*([A-z0-9\(*\)\s]+)?\s*\((?:[^)(]+|\((?:[^)(]+|\([^)(]*\))*\))*\)\s*[\(\)a-zA-Z]?.*\s?\{(?:[^}{]+|\{(?:[^}{]+|\{[^}{]*\})*\})*\}`
	CALLER_BODY   = `([a-zA-Z0-9]*.<.*[a-zA-Z0-9]>\(.*[a-zA-Z0-9]*\))`
)

// struct pattern matches
const (
	TEMPLATE_STRUCT = `@template([A-z0-9\(*\)\s].+)\ntype\s*([A-z0-9\s]+)?\s*[\(\)a-zA-Z]?.*[a-zA-Z0-9]\sstruct\s?\{(?:[^}{]+|\{(?:[^}{]+|\{[^}{]*\})*\})*\}`
	STRUCT_FIELD    = "\\`([A-z0-9\\(*\\)\\s].+)\\`"
)
