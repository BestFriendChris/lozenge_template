package tokenizer

import (
	"fmt"
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
	"github.com/BestFriendChris/lozenge_template/input"
	"github.com/BestFriendChris/lozenge_template/interfaces"
	"github.com/BestFriendChris/lozenge_template/internal/logic/token"
)

func TestContentTokenizer_ReadTokensUntil(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		inS := `
Try:
◊^{ import "fmt" }
◊{ v := 1 }
v = ◊v ◊
	  1 + 2 = ◊(1 + 2)
DONE
`[1:]
		c := ic.New(t)
		c.PrintSection("Input")
		c.Println(inS)

		tokenizer := NewDefault(interfaces.NewMacros())
		in := input.NewInput("test", inS)
		tokens, _ := tokenizer.ReadTokensUntil(in, "DONE")
		rest := in.Rest()

		logTokens(c, tokens)

		c.PrintSection("Rest")
		c.Println(rest)

		c.Expect(`
			################################################################################
			# Input
			################################################################################
			Try:
			◊^{ import "fmt" }
			◊{ v := 1 }
			v = ◊v ◊
				  1 + 2 = ◊(1 + 2)
			DONE
			
			################################################################################
			# tokens
			################################################################################
			   | TT                 | S                           | E |
			---+--------------------+-----------------------------+---+
			 1 | TT.Content         | test:1 - "Try:"             |   |
			---+--------------------+-----------------------------+---+
			 2 | TT.NL              | test:1 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			 3 | TT.CodeGlobalBlock | test:2 - " import \"fmt\" " |   |
			---+--------------------+-----------------------------+---+
			 4 | TT.NL              | test:2 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			 5 | TT.CodeLocalBlock  | test:3 - " v := 1 "         |   |
			---+--------------------+-----------------------------+---+
			 6 | TT.NL              | test:3 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			 7 | TT.Content         | test:4 - "v"                |   |
			---+--------------------+-----------------------------+---+
			 8 | TT.WS              | test:4 - " "                |   |
			---+--------------------+-----------------------------+---+
			 9 | TT.Content         | test:4 - "="                |   |
			---+--------------------+-----------------------------+---+
			10 | TT.WS              | test:4 - " "                |   |
			---+--------------------+-----------------------------+---+
			11 | TT.CodeLocalExpr   | test:4 - "v"                |   |
			---+--------------------+-----------------------------+---+
			12 | TT.WS              | test:4 - " "                |   |
			---+--------------------+-----------------------------+---+
			13 | TT.Content         | test:4 - "◊"                |   |
			---+--------------------+-----------------------------+---+
			14 | TT.NL              | test:4 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			15 | TT.WS              | test:5 - "\t  "             |   |
			---+--------------------+-----------------------------+---+
			16 | TT.Content         | test:5 - "1"                |   |
			---+--------------------+-----------------------------+---+
			17 | TT.WS              | test:5 - " "                |   |
			---+--------------------+-----------------------------+---+
			18 | TT.Content         | test:5 - "+"                |   |
			---+--------------------+-----------------------------+---+
			19 | TT.WS              | test:5 - " "                |   |
			---+--------------------+-----------------------------+---+
			20 | TT.Content         | test:5 - "2"                |   |
			---+--------------------+-----------------------------+---+
			21 | TT.WS              | test:5 - " "                |   |
			---+--------------------+-----------------------------+---+
			22 | TT.Content         | test:5 - "="                |   |
			---+--------------------+-----------------------------+---+
			23 | TT.WS              | test:5 - " "                |   |
			---+--------------------+-----------------------------+---+
			24 | TT.CodeLocalExpr   | test:5 - "(1 + 2)"          |   |
			---+--------------------+-----------------------------+---+
			25 | TT.NL              | test:5 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			################################################################################
			# Rest
			################################################################################
			DONE
			
			`)
	})
	t.Run("simple macro", func(t *testing.T) {
		s := `Hi: ◊.SimpleMacro(1 + 2) DONE`
		c := ic.New(t)
		c.PrintSection("Input")
		c.Println(s)

		tokenizer := NewDefault(mapOfSimpleMacro())
		in := input.NewInput("test", s)
		tokens, err := tokenizer.ReadTokensUntil(in, "DONE")
		if err != nil {
			panic(err)
		}
		rest := in.Rest()

		logTokens(c, tokens)

		c.PrintSection("Rest")
		c.Println(rest)

		c.Expect(`
			################################################################################
			# Input
			################################################################################
			Hi: ◊.SimpleMacro(1 + 2) DONE
			################################################################################
			# tokens
			################################################################################
			   | TT               | S                      | E |
			---+------------------+------------------------+---+
			 1 | TT.Content       | test:1 - "Hi:"         |   |
			---+------------------+------------------------+---+
			 2 | TT.WS            | test:1 - " "           |   |
			---+------------------+------------------------+---+
			 3 | TT.Macro         | test:1 - "SimpleMacro" |   |
			---+------------------+------------------------+---+
			 4 | TT.Content       | test:1 - "(1 + 2) = "  |   |
			---+------------------+------------------------+---+
			 5 | TT.CodeLocalExpr | test:1 - "(1 + 2)"     |   |
			---+------------------+------------------------+---+
			 6 | TT.WS            | test:1 - " "           |   |
			---+------------------+------------------------+---+
			################################################################################
			# Rest
			################################################################################
			DONE
			`)
	})
	t.Run("stop at end of stream", func(t *testing.T) {
		s := `Hi: ◊.SimpleMacro(1 + 2) DONE`
		c := ic.New(t)
		c.Println(`ReadTokensUntil(input, "") will read to end of stream`)
		c.PrintSection("Input")
		c.Println(s)

		tokenizer := NewDefault(mapOfSimpleMacro())
		in := input.NewInput("test", s)
		tokens, _ := tokenizer.ReadTokensUntil(in, "")
		rest := in.Rest()

		logTokens(c, tokens)

		c.PrintSection("Rest")
		c.Println(rest)

		c.Expect(`
			ReadTokensUntil(input, "") will read to end of stream
			################################################################################
			# Input
			################################################################################
			Hi: ◊.SimpleMacro(1 + 2) DONE
			################################################################################
			# tokens
			################################################################################
			   | TT               | S                      | E |
			---+------------------+------------------------+---+
			 1 | TT.Content       | test:1 - "Hi:"         |   |
			---+------------------+------------------------+---+
			 2 | TT.WS            | test:1 - " "           |   |
			---+------------------+------------------------+---+
			 3 | TT.Macro         | test:1 - "SimpleMacro" |   |
			---+------------------+------------------------+---+
			 4 | TT.Content       | test:1 - "(1 + 2) = "  |   |
			---+------------------+------------------------+---+
			 5 | TT.CodeLocalExpr | test:1 - "(1 + 2)"     |   |
			---+------------------+------------------------+---+
			 6 | TT.WS            | test:1 - " "           |   |
			---+------------------+------------------------+---+
			 7 | TT.Content       | test:1 - "DONE"        |   |
			---+------------------+------------------------+---+
			################################################################################
			# Rest
			################################################################################
			
			`)
	})
	t.Run("undefined macro", func(t *testing.T) {
		s := `Hi: ◊.UndefinedMacro(1) DONE`
		c := ic.New(t)
		c.PrintSection("Input")
		c.Println(s)

		tokenizer := NewDefault(interfaces.NewMacros())
		in := input.NewInput("test", s)
		_, err := tokenizer.ReadTokensUntil(in, "DONE")

		c.PrintSection("error")
		c.Println(err)

		c.Expect(`
			################################################################################
			# Input
			################################################################################
			Hi: ◊.UndefinedMacro(1) DONE
			################################################################################
			# error
			################################################################################
			unknown macro "UndefinedMacro"
			`)
	})
	t.Run("error if stopAt not found", func(t *testing.T) {
		s := `Hello END`
		c := ic.New(t)
		c.PrintSection("Input")
		c.Println(s)

		tokenizer := NewDefault(interfaces.NewMacros())
		in := input.NewInput("test", s)
		_, err := tokenizer.ReadTokensUntil(in, "WILL_NOT_FIND")

		c.PrintSection("error")
		c.Println(err)

		c.Expect(`
			################################################################################
			# Input
			################################################################################
			Hello END
			################################################################################
			# error
			################################################################################
			line 1: Hello END
			                 ▲
			                 └── did not find "WILL_NOT_FIND"
			`)
	})
}

func TestContentTokenizer_NextTokens(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		rest := `
Try:
◊^{ import "fmt" }
◊{ v := 1 }
v = ◊v ◊
	  1 + 2 = ◊(1 + 2)
DONE
`[1:]
		c := ic.New(t)
		c.PrintSection("Input")
		c.Println(rest)

		tokens := make([]*token.Token, 0)
		var toks []*token.Token
		tokenizer := NewDefault(interfaces.NewMacros())
		for {
			in := input.NewInput("test", rest)
			toks, _ = tokenizer.NextTokens(in)
			rest = in.Rest()

			for _, tok := range toks {
				tokens = append(tokens, tok)
			}
			if rest == "" {
				break
			}
		}

		logTokens(c, tokens)

		c.Expect(`
			################################################################################
			# Input
			################################################################################
			Try:
			◊^{ import "fmt" }
			◊{ v := 1 }
			v = ◊v ◊
				  1 + 2 = ◊(1 + 2)
			DONE
			
			################################################################################
			# tokens
			################################################################################
			   | TT                 | S                           | E |
			---+--------------------+-----------------------------+---+
			 1 | TT.Content         | test:1 - "Try:"             |   |
			---+--------------------+-----------------------------+---+
			 2 | TT.NL              | test:1 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			 3 | TT.CodeGlobalBlock | test:1 - " import \"fmt\" " |   |
			---+--------------------+-----------------------------+---+
			 4 | TT.NL              | test:1 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			 5 | TT.CodeLocalBlock  | test:1 - " v := 1 "         |   |
			---+--------------------+-----------------------------+---+
			 6 | TT.NL              | test:1 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			 7 | TT.Content         | test:1 - "v"                |   |
			---+--------------------+-----------------------------+---+
			 8 | TT.WS              | test:1 - " "                |   |
			---+--------------------+-----------------------------+---+
			 9 | TT.Content         | test:1 - "="                |   |
			---+--------------------+-----------------------------+---+
			10 | TT.WS              | test:1 - " "                |   |
			---+--------------------+-----------------------------+---+
			11 | TT.CodeLocalExpr   | test:1 - "v"                |   |
			---+--------------------+-----------------------------+---+
			12 | TT.WS              | test:1 - " "                |   |
			---+--------------------+-----------------------------+---+
			13 | TT.Content         | test:1 - "◊"                |   |
			---+--------------------+-----------------------------+---+
			14 | TT.NL              | test:1 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			15 | TT.WS              | test:1 - "\t  "             |   |
			---+--------------------+-----------------------------+---+
			16 | TT.Content         | test:1 - "1"                |   |
			---+--------------------+-----------------------------+---+
			17 | TT.WS              | test:1 - " "                |   |
			---+--------------------+-----------------------------+---+
			18 | TT.Content         | test:1 - "+"                |   |
			---+--------------------+-----------------------------+---+
			19 | TT.WS              | test:1 - " "                |   |
			---+--------------------+-----------------------------+---+
			20 | TT.Content         | test:1 - "2"                |   |
			---+--------------------+-----------------------------+---+
			21 | TT.WS              | test:1 - " "                |   |
			---+--------------------+-----------------------------+---+
			22 | TT.Content         | test:1 - "="                |   |
			---+--------------------+-----------------------------+---+
			23 | TT.WS              | test:1 - " "                |   |
			---+--------------------+-----------------------------+---+
			24 | TT.CodeLocalExpr   | test:1 - "(1 + 2)"          |   |
			---+--------------------+-----------------------------+---+
			25 | TT.NL              | test:1 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			26 | TT.Content         | test:1 - "DONE"             |   |
			---+--------------------+-----------------------------+---+
			27 | TT.NL              | test:1 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			`)
	})

	t.Run("happy path (change loz)", func(t *testing.T) {
		rest := `
Try:
^^{ import "fmt" }
^{ v := 1 }
v = ^v ^
	  1 + 2 = ^(1 + 2)
DONE
`[1:]
		c := ic.New(t)
		c.PrintSection("Input")
		c.Println(rest)

		tokens := make([]*token.Token, 0)
		var toks []*token.Token
		tokenizer := New('^', interfaces.NewMacros())
		for {
			in := input.NewInput("test", rest)
			toks, _ = tokenizer.NextTokens(in)
			rest = in.Rest()
			for _, tok := range toks {
				tokens = append(tokens, tok)
			}
			if rest == "" {
				break
			}
		}

		logTokens(c, tokens)

		c.Expect(`
			################################################################################
			# Input
			################################################################################
			Try:
			^^{ import "fmt" }
			^{ v := 1 }
			v = ^v ^
				  1 + 2 = ^(1 + 2)
			DONE
			
			################################################################################
			# tokens
			################################################################################
			   | TT                 | S                           | E |
			---+--------------------+-----------------------------+---+
			 1 | TT.Content         | test:1 - "Try:"             |   |
			---+--------------------+-----------------------------+---+
			 2 | TT.NL              | test:1 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			 3 | TT.CodeGlobalBlock | test:1 - " import \"fmt\" " |   |
			---+--------------------+-----------------------------+---+
			 4 | TT.NL              | test:1 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			 5 | TT.CodeLocalBlock  | test:1 - " v := 1 "         |   |
			---+--------------------+-----------------------------+---+
			 6 | TT.NL              | test:1 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			 7 | TT.Content         | test:1 - "v"                |   |
			---+--------------------+-----------------------------+---+
			 8 | TT.WS              | test:1 - " "                |   |
			---+--------------------+-----------------------------+---+
			 9 | TT.Content         | test:1 - "="                |   |
			---+--------------------+-----------------------------+---+
			10 | TT.WS              | test:1 - " "                |   |
			---+--------------------+-----------------------------+---+
			11 | TT.CodeLocalExpr   | test:1 - "v"                |   |
			---+--------------------+-----------------------------+---+
			12 | TT.WS              | test:1 - " "                |   |
			---+--------------------+-----------------------------+---+
			13 | TT.Content         | test:1 - "^"                |   |
			---+--------------------+-----------------------------+---+
			14 | TT.NL              | test:1 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			15 | TT.WS              | test:1 - "\t  "             |   |
			---+--------------------+-----------------------------+---+
			16 | TT.Content         | test:1 - "1"                |   |
			---+--------------------+-----------------------------+---+
			17 | TT.WS              | test:1 - " "                |   |
			---+--------------------+-----------------------------+---+
			18 | TT.Content         | test:1 - "+"                |   |
			---+--------------------+-----------------------------+---+
			19 | TT.WS              | test:1 - " "                |   |
			---+--------------------+-----------------------------+---+
			20 | TT.Content         | test:1 - "2"                |   |
			---+--------------------+-----------------------------+---+
			21 | TT.WS              | test:1 - " "                |   |
			---+--------------------+-----------------------------+---+
			22 | TT.Content         | test:1 - "="                |   |
			---+--------------------+-----------------------------+---+
			23 | TT.WS              | test:1 - " "                |   |
			---+--------------------+-----------------------------+---+
			24 | TT.CodeLocalExpr   | test:1 - "(1 + 2)"          |   |
			---+--------------------+-----------------------------+---+
			25 | TT.NL              | test:1 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			26 | TT.Content         | test:1 - "DONE"             |   |
			---+--------------------+-----------------------------+---+
			27 | TT.NL              | test:1 - "\n"               |   |
			---+--------------------+-----------------------------+---+
			`)
	})
	t.Run("whitespace", func(t *testing.T) {
		tok, rest, _ := readNextToken(t, "\t   hi")

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.WS
			Token.Slc: test:1 - "\t   "
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"hi"
			`)
	})
	t.Run("newline", func(t *testing.T) {
		tok, rest, _ := readNextToken(t, "\n\nfoo")

		c := ic.New(t)
		c.Println("Only read one newline")
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			Only read one newline
			################################################################################
			# token
			################################################################################
			Token.TT: TT.NL
			Token.Slc: test:1 - "\n"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"\nfoo"
			`)
	})
	t.Run("content", func(t *testing.T) {
		tok, rest, _ := readNextToken(t, "foo\nbar")

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.Content
			Token.Slc: test:1 - "foo"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"\nbar"
			`)
	})
	t.Run("◊◊foo", func(t *testing.T) {
		tok, rest, _ := readNextToken(t, `◊◊foo`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.Content
			Token.Slc: test:1 - "◊"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"foo"
			`)
	})
	t.Run("◊ bar", func(t *testing.T) {
		tok, rest, _ := readNextToken(t, `◊ bar`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.Content
			Token.Slc: test:1 - "◊"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			" bar"
			`)
	})
	t.Run(`◊\nbar`, func(t *testing.T) {
		tok, rest, _ := readNextToken(t, "◊\nbar")

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.Content
			Token.Slc: test:1 - "◊"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"\nbar"
			`)
	})
	t.Run(`◊`, func(t *testing.T) {
		tok, rest, _ := readNextToken(t, "◊")

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.Content
			Token.Slc: test:1 - "◊"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			""
			`)
	})
	t.Run("◊foo bar", func(t *testing.T) {
		tok, rest, _ := readNextToken(t, `◊foo bar`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.CodeLocalExpr
			Token.Slc: test:1 - "foo"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			" bar"
			`)
	})
	t.Run("◊foo", func(t *testing.T) {
		tok, rest, _ := readNextToken(t, `◊foo`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.CodeLocalExpr
			Token.Slc: test:1 - "foo"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			""
			`)
	})
	t.Run("◊(1 + 2)foo", func(t *testing.T) {
		tok, rest, _ := readNextToken(t, `◊(1 + 2)foo`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.CodeLocalExpr
			Token.Slc: test:1 - "(1 + 2)"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"foo"
			`)
	})
	t.Run("◊(1 + 2", func(t *testing.T) {
		_, _, err := readNextToken(t, `◊(1 + 2`)

		c := ic.New(t)
		logErr(c, err)
		c.Expect(`
			################################################################################
			# error
			################################################################################
			line 1: ◊(1 + 2
			         ▲
			         └── did not find matched ')'
			`)
	})
	t.Run("◊{ GOCODE }", func(t *testing.T) {
		tok, rest, _ := readNextToken(t, `◊{ var foo, bar, baz := struct{a string}{"}\""}, '}', '\'' }foo`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.CodeLocalBlock
			Token.Slc: test:1 - " var foo, bar, baz := struct{a string}{\"}\\\"\"}, '}', '\\'' "
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"foo"
			`)
	})
	t.Run(`◊{ GOCODE \n GOCODE }`, func(t *testing.T) {
		s := `
◊{
  var foo := struct{a string}{"}\""}
  var bar := '}'
  var baz := '\''
}foo`[1:]
		toks, rest, _ := readNextNTokens(t, 3, s)

		c := ic.New(t)
		logTokens(c, toks)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# tokens
			################################################################################
			   | TT                | S                                                   | E |
			---+-------------------+-----------------------------------------------------+---+
			 1 | TT.CodeLocalBlock | test:2 - "  var foo := struct{a string}{\"}\\\"\"}" |   |
			---+-------------------+-----------------------------------------------------+---+
			 2 | TT.CodeLocalBlock | test:3 - "  var bar := '}'"                         |   |
			---+-------------------+-----------------------------------------------------+---+
			 3 | TT.CodeLocalBlock | test:4 - "  var baz := '\\''"                       |   |
			---+-------------------+-----------------------------------------------------+---+
			################################################################################
			# rest
			################################################################################
			"foo"
			`)
	})
	t.Run("◊{ GOCODE ", func(t *testing.T) {
		_, _, err := readNextToken(t, `◊{ var foo struct{a string}{"}"} foo`)

		c := ic.New(t)
		logErr(c, err)
		c.Expect(`
			################################################################################
			# error
			################################################################################
			line 1: ◊{ var foo struct{a string}{"}"} foo
			         ▲
			         └── did not find matched '}'
			`)
	})
	t.Run("◊^{ GOCODE }", func(t *testing.T) {
		tok, rest, _ := readNextToken(t, `◊^{ import "bar" }foo`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.CodeGlobalBlock
			Token.Slc: test:1 - " import \"bar\" "
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"foo"
			`)
	})
	t.Run(`◊^{ GOCODE \n GOCODE }`, func(t *testing.T) {
		toks, rest, _ := readNextNTokens(t, 3, `
◊^{
	type foo struct{
		a string
	}
}foo`[1:])

		c := ic.New(t)
		logTokens(c, toks)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# tokens
			################################################################################
			   | TT                 | S                             | E |
			---+--------------------+-------------------------------+---+
			 1 | TT.CodeGlobalBlock | test:2 - "\ttype foo struct{" |   |
			---+--------------------+-------------------------------+---+
			 2 | TT.CodeGlobalBlock | test:3 - "\t\ta string"       |   |
			---+--------------------+-------------------------------+---+
			 3 | TT.CodeGlobalBlock | test:4 - "\t}"                |   |
			---+--------------------+-------------------------------+---+
			################################################################################
			# rest
			################################################################################
			"foo"
			`)
	})
	t.Run("◊^{ GOCODE ", func(t *testing.T) {
		_, _, err := readNextToken(t, `
◊^{
	type foo struct{
		a string
	}
foo`[1:])

		c := ic.New(t)
		logErr(c, err)
		c.Expect(`
			################################################################################
			# error
			################################################################################
			line 1: ◊^{
			          ▲
			          └── did not find matched '}'
			`)
	})
	t.Run("◊^foo", func(t *testing.T) {
		tok, rest, _ := readNextToken(t, `◊^foo`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.Content
			Token.Slc: test:1 - "◊"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"^foo"
			`)
	})
	t.Run("◊.macro", func(t *testing.T) {
		tokenizer := NewDefault(mapOfSimpleMacro())
		in := input.NewInput("test", `◊.SimpleMacro(1 + 2)`)
		tokens, _ := tokenizer.NextTokens(in)
		rest := in.Rest()

		c := ic.New(t)
		logTokens(c, tokens)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# tokens
			################################################################################
			   | TT               | S                      | E |
			---+------------------+------------------------+---+
			 1 | TT.Macro         | test:1 - "SimpleMacro" |   |
			---+------------------+------------------------+---+
			 2 | TT.Content       | test:1 - "(1 + 2) = "  |   |
			---+------------------+------------------------+---+
			 3 | TT.CodeLocalExpr | test:1 - "(1 + 2)"     |   |
			---+------------------+------------------------+---+
			################################################################################
			# rest
			################################################################################
			""
			`)
	})
}

func TestContentTokenizer_NextTokenCodeUntilOpenBrace(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		tokenizer := NewDefault(interfaces.NewMacros())
		in := input.NewInput("test", `if strings.DeepEqual(v, []string{"\"", "{"}) {◊foo`)
		tok, _ := tokenizer.NextTokenCodeUntilOpenBraceLoz(in)
		rest := in.Rest()
		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.CodeLocalBlock
			Token.Slc: test:1 - "if strings.DeepEqual(v, []string{\"\\\"\", \"{\"}) {"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"foo"
			`)
	})
	t.Run("no open brace", func(t *testing.T) {
		tokenizer := NewDefault(interfaces.NewMacros())
		in := input.NewInput("test", `if "\"" == "{" ◊} else { foo ◊}`)
		_, err := tokenizer.NextTokenCodeUntilOpenBraceLoz(in)
		c := ic.New(t)
		logErr(c, err)
		c.Expect(`
			################################################################################
			# error
			################################################################################
			line 1: if "\"" == "{" ◊} else { foo ◊}
			        ▲
			        └── no open brace found
			`)
	})
}

func readNextToken(t *testing.T, s string) (*token.Token, string, error) {
	tokens, rest, err := readNextNTokens(t, 1, s)
	if tokens == nil {
		return nil, rest, err
	}
	return tokens[0], rest, err
}

func readNextNTokens(t *testing.T, n int, s string) ([]*token.Token, string, error) {
	tokenizer := NewDefault(mapOfSimpleMacro())
	in := input.NewInput("test", s)

	tokens, err := tokenizer.NextTokens(in)
	if err != nil {
		return nil, "", err
	}

	rest := in.Rest()
	if len(tokens) != n {
		t.Fatalf("got len(%d) want len(%d):\n%v", len(tokens), n, tokens)
	}
	return tokens, rest, err
}

func logErr(c *ic.IC, err error) {
	c.PrintSection("error")
	c.Println(err)
}

func logToken(c *ic.IC, t *token.Token) {
	c.PrintSection("token")
	if t == nil {
		c.Println("<null>")
	} else {
		c.PV(t)
	}
}

func logTokens(c *ic.IC, tokens []*token.Token) {
	c.PrintSection("tokens")
	type tokensTable struct {
		TT token.TokenType
		S  input.Slice
		E  *any
	}
	tt := make([]tokensTable, len(tokens))
	for i, toks := range tokens {
		tt[i] = tokensTable{toks.TT, toks.Slc, toks.E}
	}
	c.PT(tt)
}

func logRest(c *ic.IC, rest string) {
	c.PrintSection("rest")
	c.Printf("%q\n", rest)
}

func mapOfSimpleMacro() *interfaces.Macros {
	macros := interfaces.NewMacros()
	macros.Add(&simpleMacro{})
	return macros
}

// ◊.SimpleMacro(1 + 2)
type simpleMacro struct {
}

func (m simpleMacro) Name() string {
	return "SimpleMacro"
}

func (m simpleMacro) NextTokens(ct interfaces.ContentTokenizer, in *input.Input) (toks []*token.Token, err error) {
	_, _ = in.ConsumeString(m.Name())
	var valToks []*token.Token
	valToks, err = ct.ParseGoCodeFromTo(in, token.TTcodeLocalExpr, '(', ')', true)
	if err != nil {
		return nil, err
	}
	if len(valToks) != 1 {
		return nil, fmt.Errorf("got len(%d) want len(1) of valtoks\n%v", len(valToks), valToks)
	}
	valTok := valToks[0]
	contentSlc := input.NewSlice("test", fmt.Sprintf("%s = ", valTok.Slc.S), valTok.Slc.Start, valTok.Slc.End)
	contentToken := token.NewToken(token.TTcontent, contentSlc)
	return []*token.Token{contentToken, valTok}, nil
}

func (m simpleMacro) Parse(_ interfaces.TemplateHandler, toks []*token.Token) (rest []*token.Token, err error) {
	return toks, nil
}
