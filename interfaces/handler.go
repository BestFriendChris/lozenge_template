package interfaces

type TemplateHandler interface {
	DefaultMacros() *Macros
	WriteTextContent(string)
	WriteCodeGlobalBlock(string)
	WriteCodeLocalExpression(string)
	WriteCodeLocalBlock(string)
	Done() (string, error)
}
