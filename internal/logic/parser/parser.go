package parser

import (
	"fmt"

	"github.com/BestFriendChris/lozenge/interfaces"
	"github.com/BestFriendChris/lozenge/internal/logic/token"
)

func New() *DefaultParser {
	return &DefaultParser{}
}

type DefaultParser struct {
	macros interfaces.MapMacros
}

func (p *DefaultParser) Parse(h interfaces.TemplateHandler, toks []*token.Token) ([]*token.Token, error) {
	var idx int
	for i, tok := range toks {
		i = idx
		switch tok.TT {
		case token.TTcontent, token.TTnl, token.TTws:
			h.WriteTextContent(tok.S)
		case token.TTcodeGlobalBlock:
			h.WriteCodeGlobalBlock(tok.S)
		case token.TTcodeLocalBlock:
			h.WriteCodeLocalBlock(tok.S)
		case token.TTcodeLocalExpr:
			h.WriteCodeLocalExpression(tok.S)
		default:
			return toks[i:], fmt.Errorf("unrecognized token type %q: %s", tok.TT, tok)
		}
	}
	return toks[idx:], nil
}
