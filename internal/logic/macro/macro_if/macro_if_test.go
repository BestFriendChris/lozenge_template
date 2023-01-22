package macro_if

import (
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
	"github.com/BestFriendChris/lozenge/interfaces"
	"github.com/BestFriendChris/lozenge/internal/logic/token"
	"github.com/BestFriendChris/lozenge/internal/logic/tokenizer"
)

func TestMacroIf_NextTokens(t *testing.T) {
	t.Run("basic if", func(t *testing.T) {
		ct := tokenizer.NewDefault(interfaces.NewMapMacros())

		rest := `
if reflect.DeepEqual(val, []string{"foo"}) {◊
  hi
◊}bar`[1:]

		macroIf := New(ct)

		var tokens []*token.Token
		tokens, rest = macroIf.NextTokens(rest)

		c := ic.New(t)
		c.PrintSection("tokens")
		c.PT(tokens)

		c.PrintSection("rest")
		c.Println(rest)
		c.Expect(`
			################################################################################
			# tokens
			################################################################################
			   | TT                | S                                                | E |
			---+-------------------+--------------------------------------------------+---+
			 1 | TT.CodeLocalBlock | "if reflect.DeepEqual(val, []string{\"foo\"}) {" |   |
			---+-------------------+--------------------------------------------------+---+
			 2 | TT.NL             | "\n"                                             |   |
			---+-------------------+--------------------------------------------------+---+
			 3 | TT.WS             | "  "                                             |   |
			---+-------------------+--------------------------------------------------+---+
			 4 | TT.Content        | "hi"                                             |   |
			---+-------------------+--------------------------------------------------+---+
			 5 | TT.NL             | "\n"                                             |   |
			---+-------------------+--------------------------------------------------+---+
			 6 | TT.CodeLocalBlock | "}"                                              |   |
			---+-------------------+--------------------------------------------------+---+
			################################################################################
			# rest
			################################################################################
			bar
			`)
	})
	t.Run("basic if else", func(t *testing.T) {
		ct := tokenizer.NewDefault(interfaces.NewMapMacros())

		rest := `
if true {◊
  foo
◊}  else  {◊
  bar
◊}baz
`[1:]

		macroIf := New(ct)

		var tokens []*token.Token
		tokens, rest = macroIf.NextTokens(rest)

		c := ic.New(t)
		c.PrintSection("tokens")
		c.PT(tokens)

		c.PrintSection("rest")
		c.Println(rest)

		c.Expect(`
			################################################################################
			# tokens
			################################################################################
			   | TT                | S            | E |
			---+-------------------+--------------+---+
			 1 | TT.CodeLocalBlock | "if true {"  |   |
			---+-------------------+--------------+---+
			 2 | TT.NL             | "\n"         |   |
			---+-------------------+--------------+---+
			 3 | TT.WS             | "  "         |   |
			---+-------------------+--------------+---+
			 4 | TT.Content        | "foo"        |   |
			---+-------------------+--------------+---+
			 5 | TT.NL             | "\n"         |   |
			---+-------------------+--------------+---+
			 6 | TT.CodeLocalBlock | "}  else  {" |   |
			---+-------------------+--------------+---+
			 7 | TT.NL             | "\n"         |   |
			---+-------------------+--------------+---+
			 8 | TT.WS             | "  "         |   |
			---+-------------------+--------------+---+
			 9 | TT.Content        | "bar"        |   |
			---+-------------------+--------------+---+
			10 | TT.NL             | "\n"         |   |
			---+-------------------+--------------+---+
			11 | TT.CodeLocalBlock | "}"          |   |
			---+-------------------+--------------+---+
			################################################################################
			# rest
			################################################################################
			baz
			
			`)
	})
	t.Run("basic if else if", func(t *testing.T) {
		ct := tokenizer.NewDefault(interfaces.NewMapMacros())

		rest := `
if v == 1 {◊
  one
◊}  else  if v == 2 {◊
  two
◊}  else  if  v == 3 {◊
  three
◊}  else {◊
  four
◊}baz
`[1:]

		macroIf := New(ct)

		var tokens []*token.Token
		tokens, rest = macroIf.NextTokens(rest)

		c := ic.New(t)
		c.PrintSection("tokens")
		c.PT(tokens)

		c.PrintSection("rest")
		c.Println(rest)

		c.Expect(`
			################################################################################
			# tokens
			################################################################################
			   | TT                | S                       | E |
			---+-------------------+-------------------------+---+
			 1 | TT.CodeLocalBlock | "if v == 1 {"           |   |
			---+-------------------+-------------------------+---+
			 2 | TT.NL             | "\n"                    |   |
			---+-------------------+-------------------------+---+
			 3 | TT.WS             | "  "                    |   |
			---+-------------------+-------------------------+---+
			 4 | TT.Content        | "one"                   |   |
			---+-------------------+-------------------------+---+
			 5 | TT.NL             | "\n"                    |   |
			---+-------------------+-------------------------+---+
			 6 | TT.CodeLocalBlock | "}  else  if v == 2 {"  |   |
			---+-------------------+-------------------------+---+
			 7 | TT.NL             | "\n"                    |   |
			---+-------------------+-------------------------+---+
			 8 | TT.WS             | "  "                    |   |
			---+-------------------+-------------------------+---+
			 9 | TT.Content        | "two"                   |   |
			---+-------------------+-------------------------+---+
			10 | TT.NL             | "\n"                    |   |
			---+-------------------+-------------------------+---+
			11 | TT.CodeLocalBlock | "}  else  if  v == 3 {" |   |
			---+-------------------+-------------------------+---+
			12 | TT.NL             | "\n"                    |   |
			---+-------------------+-------------------------+---+
			13 | TT.WS             | "  "                    |   |
			---+-------------------+-------------------------+---+
			14 | TT.Content        | "three"                 |   |
			---+-------------------+-------------------------+---+
			15 | TT.NL             | "\n"                    |   |
			---+-------------------+-------------------------+---+
			16 | TT.CodeLocalBlock | "}  else {"             |   |
			---+-------------------+-------------------------+---+
			17 | TT.NL             | "\n"                    |   |
			---+-------------------+-------------------------+---+
			18 | TT.WS             | "  "                    |   |
			---+-------------------+-------------------------+---+
			19 | TT.Content        | "four"                  |   |
			---+-------------------+-------------------------+---+
			20 | TT.NL             | "\n"                    |   |
			---+-------------------+-------------------------+---+
			21 | TT.CodeLocalBlock | "}"                     |   |
			---+-------------------+-------------------------+---+
			################################################################################
			# rest
			################################################################################
			baz
			
			`)
	})
}

func TestMacroIf_NextTokens_errorCases(t *testing.T) {
	t.Run("no closing brace with if", func(t *testing.T) {
		input := `if true `
		ct := tokenizer.NewDefault(interfaces.NewMapMacros())
		macroIf := New(ct)

		tokens, rest := macroIf.NextTokens(input)

		c := ic.New(t)
		c.PrintSection("tokens")
		c.PT(tokens)

		c.PrintSection("rest")
		c.Println(rest)

		c.Expect(`
			################################################################################
			# tokens
			################################################################################
			   | TT | S | E |
			---+----+---+---+
			################################################################################
			# rest
			################################################################################
			if true 
			`)
	})
	t.Run("no closing brace with else", func(t *testing.T) {
		input := `
if v == 1 {◊
  one
◊}  else 
  four
◊}baz
`[1:]
		ct := tokenizer.NewDefault(interfaces.NewMapMacros())
		macroIf := New(ct)

		tokens, rest := macroIf.NextTokens(input)

		c := ic.New(t)
		c.PrintSection("tokens")
		c.PT(tokens)

		c.PrintSection("rest")
		c.Println(rest)

		c.Expect(`
			################################################################################
			# tokens
			################################################################################
			   | TT | S | E |
			---+----+---+---+
			################################################################################
			# rest
			################################################################################
			if v == 1 {◊
			  one
			◊}  else 
			  four
			◊}baz
			
			`)
	})
	t.Run("no closing brace with else if", func(t *testing.T) {
		input := `
if v == 1 {◊
  one
◊}  else  if v == 2 {◊
  two
◊}  else  if  v == 3 
  three
◊}  else {◊
  four
◊}baz
`[1:]
		ct := tokenizer.NewDefault(interfaces.NewMapMacros())
		macroIf := New(ct)

		tokens, rest := macroIf.NextTokens(input)

		c := ic.New(t)
		c.PrintSection("tokens")
		c.PT(tokens)

		c.PrintSection("rest")
		c.Println(rest)

		c.Expect(`
			################################################################################
			# tokens
			################################################################################
			   | TT | S | E |
			---+----+---+---+
			################################################################################
			# rest
			################################################################################
			if v == 1 {◊
			  one
			◊}  else  if v == 2 {◊
			  two
			◊}  else  if  v == 3 
			  three
			◊}  else {◊
			  four
			◊}baz
			
			`)
	})
}
