package macro_for

import (
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
	"github.com/BestFriendChris/lozenge/interfaces"
	"github.com/BestFriendChris/lozenge/internal/logic/token"
	"github.com/BestFriendChris/lozenge/internal/logic/tokenizer"
)

func TestMacroFor_NextTokens(t *testing.T) {
	t.Run("basic for", func(t *testing.T) {
		ct := tokenizer.NewDefault(interfaces.NewMacros())

		rest := `
for _, v := range vals {◊
	<span>◊v</span>
◊}bar`[1:]

		macroFor := New()

		var tokens []*token.Token
		tokens, rest, _ = macroFor.NextTokens(ct, rest)

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
	t.Run("no open brace with if", func(t *testing.T) {
		ct := tokenizer.NewDefault(interfaces.NewMacros())

		input := `
for _, v := range vals 
	<span>◊v</span>
◊}bar`[1:]

		macroFor := New()

		_, _, err := macroFor.NextTokens(ct, input)

		c := ic.New(t)
		c.PrintSection("error")
		c.Println(err)

		c.Expect(`
			################################################################################
			# error
			################################################################################
			macro(for): no open brace found
			`)
	})
	t.Run("no close brace with if", func(t *testing.T) {
		ct := tokenizer.NewDefault(interfaces.NewMacros())

		input := `
for _, v := range vals {◊
	<span>◊v</span>
`[1:]

		macroFor := New()

		_, _, err := macroFor.NextTokens(ct, input)

		c := ic.New(t)
		c.PrintSection("error")
		c.Println(err)

		c.Expect(`
			################################################################################
			# error
			################################################################################
			macro(for): did not find "◊}"
			`)
	})
}
