package parser

type Handler interface {
	DefaultMacros() map[string]Macro
	WriteContent(string)
	WriteCodeExpression(string)
	WriteCodeBlock(string)
	WriteCodeGlobalBlock(string)
	Done() (string, error)
}
