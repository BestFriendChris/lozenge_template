package lozenge_template

import (
	"github.com/BestFriendChris/lozenge/interfaces"
	"github.com/BestFriendChris/lozenge/internal/infra/go_format"
	"github.com/BestFriendChris/lozenge/internal/logic/macro/macro_for"
	"github.com/BestFriendChris/lozenge/internal/logic/macro/macro_if"
	"github.com/BestFriendChris/lozenge/internal/logic/parser"
	"github.com/BestFriendChris/lozenge/internal/logic/tokenizer"
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

func (lt *LozengeTemplate) Generate(h interfaces.TemplateHandler, input string) (goCode string, err error) {
	macros := lt.defaultMacros.Merge(h.DefaultMacros())

	ct := tokenizer.NewDefault(macros)
	toks := ct.ReadAll(input)
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

	goCode, err = go_format.Format(goCode)
	return
}

func defaultMacros(overrideMacros *interfaces.Macros) *interfaces.Macros {
	macros := interfaces.NewMacros()
	macros.Add(macro_if.New())
	macros.Add(macro_for.New())

	macros = macros.Merge(overrideMacros)

	return macros
}
