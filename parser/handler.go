package parser

import (
	"mvdan.cc/gofumpt/format"
)

type Handler interface {
	DefaultMacros() map[string]Macro
	WriteContent(string)
	WriteCodeInline(string)
	WriteCodeBlock(string)
	WriteCodeGlobalBlock(string)
	Done(format.Options) (string, error)
}
