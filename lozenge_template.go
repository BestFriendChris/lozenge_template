package lozenge_template

import (
	"github.com/BestFriendChris/lozenge_template/input"
	"github.com/BestFriendChris/lozenge_template/interfaces"
	"github.com/BestFriendChris/lozenge_template/internal/infra/go_format"
	"github.com/BestFriendChris/lozenge_template/internal/logic/macro/macro_for"
	"github.com/BestFriendChris/lozenge_template/internal/logic/macro/macro_if"
	"github.com/BestFriendChris/lozenge_template/internal/logic/parser"
	"github.com/BestFriendChris/lozenge_template/internal/logic/token"
	"github.com/BestFriendChris/lozenge_template/internal/logic/tokenizer"
)

func New(overrideMacros *interfaces.Macros, config ParserConfig) *LozengeTemplate {
	return &LozengeTemplate{
		config:        config,
		defaultMacros: defaultMacros(overrideMacros),
	}
}

type LozengeTemplate struct {
	config        ParserConfig
	defaultMacros *interfaces.Macros
}

func (lt *LozengeTemplate) Generate(h interfaces.TemplateHandler, in *input.Input) (goCode string, err error) {
	macros := lt.defaultMacros.Merge(h.DefaultMacros())

	ct := tokenizer.New(lt.config.Loz, macros)

	var toks []*token.Token

	toks, err = ct.ReadAll(in)
	if err != nil {
		return "", err
	}
	toks = tokenizer.Optimize(toks, lt.config.TrimSpaces)

	prs := parser.New(macros)
	_, err = prs.Parse(h, toks)

	if err != nil {
		return "", err
	}

	goCode, err = h.Done()
	if err != nil {
		return "", err
	}

	return go_format.Format(goCode)
}

func defaultMacros(overrideMacros *interfaces.Macros) *interfaces.Macros {
	macros := interfaces.NewMacros()
	macros.Add(macro_if.New())
	macros.Add(macro_for.New())

	macros = macros.Merge(overrideMacros)

	return macros
}
