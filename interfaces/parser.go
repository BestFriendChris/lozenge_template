package interfaces

import (
	"github.com/BestFriendChris/lozenge_template/internal/logic/token"
)

type Parser interface {
	Parse(h TemplateHandler, toks []*token.Token) (rest []*token.Token, err error)
}

type ParserSlc interface {
	ParseSlc(h TemplateHandler, toks []*token.TokenSlice) (rest []*token.TokenSlice, err error)
}
