package interfaces

type TemplateHandler interface {
	DefaultMacros() MapMacros
	WriteTextContent(string)
	WriteCodeGlobalBlock(string)
	WriteCodeLocalExpression(string)
	WriteCodeLocalBlock(string)
	Done() (string, error)
}
