package interfaces

import "github.com/BestFriendChris/lozenge_template/input"

type TemplateHandler interface {
	DefaultMacros() *Macros
	WriteTextContent(input.Slice)
	WriteCodeGlobalBlock(input.Slice)
	WriteCodeLocalExpression(input.Slice)
	WriteCodeLocalBlock(input.Slice)
	Done() (string, error)
}
