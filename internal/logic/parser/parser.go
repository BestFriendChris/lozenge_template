package parser

import (
	"fmt"

	"github.com/BestFriendChris/lozenge_template/interfaces"
	"github.com/BestFriendChris/lozenge_template/internal/logic/token"
)

func New(macros *interfaces.Macros) *DefaultParser {
	return &DefaultParser{
		macros: macros,
	}
}

type DefaultParser struct {
	macros *interfaces.Macros
}

func (p *DefaultParser) Parse(h interfaces.TemplateHandler, toks []*token.Token) (rest []*token.Token, err error) {
	var idx int
	for i, tok := range toks {
		idx = i
		switch tok.TT {
		case token.TTcontent, token.TTnl, token.TTws:
			h.WriteTextContent(tok.S)
		case token.TTcodeGlobalBlock:
			h.WriteCodeGlobalBlock(tok.S)
		case token.TTcodeLocalBlock:
			h.WriteCodeLocalBlock(tok.S)
		case token.TTcodeLocalExpr:
			h.WriteCodeLocalExpression(tok.S)
		case token.TTmacro:
			if p.macros == nil {
				return toks, fmt.Errorf("parser: unknown macro %q", tok.S)
			}
			m, found := p.macros.Get(tok.S)
			if found {
				rest, err := m.Parse(h, toks)
				if err != nil {
					return rest, err
				}
			} else {
				return toks, fmt.Errorf("parser: unknown macro %q", tok.S)
			}
		default:
			return toks[i:], fmt.Errorf("parser: unrecognized token type %q: %s", tok.TT, tok)
		}
	}
	return toks[idx:], nil
}

func (p *DefaultParser) ParseSlc(h interfaces.TemplateHandler, toks []*token.TokenSlice) (rest []*token.TokenSlice, err error) {
	var idx int
	for i, tok := range toks {
		idx = i
		switch tok.TT {
		case token.TTcontent, token.TTnl, token.TTws:
			h.WriteTextContent(tok.Slc.S)
		case token.TTcodeGlobalBlock:
			h.WriteCodeGlobalBlock(tok.Slc.S)
		case token.TTcodeLocalBlock:
			h.WriteCodeLocalBlock(tok.Slc.S)
		case token.TTcodeLocalExpr:
			h.WriteCodeLocalExpression(tok.Slc.S)
		case token.TTmacro:
			if p.macros == nil {
				return toks, fmt.Errorf("parser: unknown macro %q", tok.Slc.S)
			}
			m, found := p.macros.Get(tok.Slc.S)
			if found {
				rest, err := m.ParseSlc(h, toks)
				if err != nil {
					return rest, err
				}
			} else {
				return toks, fmt.Errorf("parser: unknown macro %q", tok.Slc.S)
			}
		default:
			return toks[i:], fmt.Errorf("parser: unrecognized token type %q: %s", tok.TT, tok)
		}
	}
	return toks[idx:], nil
}
