package macro_for

import (
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
	"github.com/BestFriendChris/lozenge_template/input"
	"github.com/BestFriendChris/lozenge_template/interfaces"
	"github.com/BestFriendChris/lozenge_template/internal/logic/token"
	"github.com/BestFriendChris/lozenge_template/internal/logic/tokenizer"
)

func TestMacroFor_NextTokens(t *testing.T) {
	t.Run("basic for", func(t *testing.T) {
		ct := tokenizer.NewDefault(interfaces.NewMacros())

		s := `
for _, v := range vals {◊
	<span>◊v</span>
◊}bar`[1:]

		macroFor := New()

		var tokens []*token.Token
		in := input.NewInput(s)
		tokens, _ = macroFor.NextTokens(ct, in)
		rest := in.Rest()

		c := ic.New(t)
		c.PrintSection("tokens")
		c.PT(tokens)

		c.PrintSection("rest")
		c.Println(rest)
		c.Expect(`
			################################################################################
			# tokens
			################################################################################
			   | TT                | S                          | E |
			---+-------------------+----------------------------+---+
			 1 | TT.CodeLocalBlock | "for _, v := range vals {" |   |
			---+-------------------+----------------------------+---+
			 2 | TT.NL             | "\n"                       |   |
			---+-------------------+----------------------------+---+
			 3 | TT.WS             | "\t"                       |   |
			---+-------------------+----------------------------+---+
			 4 | TT.Content        | "<span>"                   |   |
			---+-------------------+----------------------------+---+
			 5 | TT.CodeLocalExpr  | "v"                        |   |
			---+-------------------+----------------------------+---+
			 6 | TT.Content        | "</span>"                  |   |
			---+-------------------+----------------------------+---+
			 7 | TT.NL             | "\n"                       |   |
			---+-------------------+----------------------------+---+
			 8 | TT.CodeLocalBlock | "}"                        |   |
			---+-------------------+----------------------------+---+
			################################################################################
			# rest
			################################################################################
			bar
			`)
	})
}

func TestMacroFor_NextTokens_errorCases(t *testing.T) {
	t.Run("no open brace with for", func(t *testing.T) {
		ct := tokenizer.NewDefault(interfaces.NewMacros())

		s := `
for _, v := range vals 
	<span>◊v</span>
◊}bar`[1:]

		macroFor := New()

		in := input.NewInput(s)
		_, err := macroFor.NextTokens(ct, in)

		c := ic.New(t)
		c.PrintSection("error")
		c.Println(err)

		c.Expect(`
			################################################################################
			# error
			################################################################################
			line 1: for _, v := range vals 
			        ▲
			        └── no open brace found
			`)
	})
	t.Run("no close brace with if", func(t *testing.T) {
		ct := tokenizer.NewDefault(interfaces.NewMacros())

		s := `
for _, v := range vals {◊
	<span>◊v</span>
`[1:]

		macroFor := New()

		in := input.NewInput(s)
		_, err := macroFor.NextTokens(ct, in)

		c := ic.New(t)
		c.PrintSection("error")
		c.Println(err)

		c.Expect(`
			################################################################################
			# error
			################################################################################
			line 3: 
			        ▲
			        └── did not find "◊}"
			`)
	})
}
