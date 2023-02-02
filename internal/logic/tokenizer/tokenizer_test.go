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
			   | TT                 | S                  | E |
			---+--------------------+--------------------+---+
			 1 | TT.Content         | "Try:"             |   |
			---+--------------------+--------------------+---+
			 2 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
			 3 | TT.CodeGlobalBlock | " import \"fmt\" " |   |
			---+--------------------+--------------------+---+
			 4 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
			 5 | TT.CodeLocalBlock  | " v := 1 "         |   |
			---+--------------------+--------------------+---+
			 6 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
			 7 | TT.Content         | "v"                |   |
			---+--------------------+--------------------+---+
			 8 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			 9 | TT.Content         | "="                |   |
			---+--------------------+--------------------+---+
			10 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			11 | TT.CodeLocalExpr   | "v"                |   |
			---+--------------------+--------------------+---+
			12 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			13 | TT.Content         | "◊"                |   |
			---+--------------------+--------------------+---+
			14 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
			15 | TT.WS              | "\t  "             |   |
			---+--------------------+--------------------+---+
			16 | TT.Content         | "1"                |   |
			---+--------------------+--------------------+---+
			17 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			18 | TT.Content         | "+"                |   |
			---+--------------------+--------------------+---+
			19 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			20 | TT.Content         | "2"                |   |
			---+--------------------+--------------------+---+
			21 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			22 | TT.Content         | "="                |   |
			---+--------------------+--------------------+---+
			23 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			24 | TT.CodeLocalExpr   | "(1 + 2)"          |   |
			---+--------------------+--------------------+---+
			25 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
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
			   | TT               | S             | E |
			---+------------------+---------------+---+
			 1 | TT.Content       | "Hi:"         |   |
			---+------------------+---------------+---+
			 2 | TT.WS            | " "           |   |
			---+------------------+---------------+---+
			 3 | TT.Macro         | "SimpleMacro" |   |
			---+------------------+---------------+---+
			 4 | TT.Content       | "(1 + 2) = "  |   |
			---+------------------+---------------+---+
			 5 | TT.CodeLocalExpr | "(1 + 2)"     |   |
			---+------------------+---------------+---+
			 6 | TT.WS            | " "           |   |
			---+------------------+---------------+---+
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
			   | TT               | S             | E |
			---+------------------+---------------+---+
			 1 | TT.Content       | "Hi:"         |   |
			---+------------------+---------------+---+
			 2 | TT.WS            | " "           |   |
			---+------------------+---------------+---+
			 3 | TT.Macro         | "SimpleMacro" |   |
			---+------------------+---------------+---+
			 4 | TT.Content       | "(1 + 2) = "  |   |
			---+------------------+---------------+---+
			 5 | TT.CodeLocalExpr | "(1 + 2)"     |   |
			---+------------------+---------------+---+
			 6 | TT.WS            | " "           |   |
			---+------------------+---------------+---+
			 7 | TT.Content       | "DONE"        |   |
			---+------------------+---------------+---+
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

func TestContentTokenizer_NextToken(t *testing.T) {
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
			   | TT                 | S                  | E |
			---+--------------------+--------------------+---+
			 1 | TT.Content         | "Try:"             |   |
			---+--------------------+--------------------+---+
			 2 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
			 3 | TT.CodeGlobalBlock | " import \"fmt\" " |   |
			---+--------------------+--------------------+---+
			 4 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
			 5 | TT.CodeLocalBlock  | " v := 1 "         |   |
			---+--------------------+--------------------+---+
			 6 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
			 7 | TT.Content         | "v"                |   |
			---+--------------------+--------------------+---+
			 8 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			 9 | TT.Content         | "="                |   |
			---+--------------------+--------------------+---+
			10 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			11 | TT.CodeLocalExpr   | "v"                |   |
			---+--------------------+--------------------+---+
			12 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			13 | TT.Content         | "◊"                |   |
			---+--------------------+--------------------+---+
			14 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
			15 | TT.WS              | "\t  "             |   |
			---+--------------------+--------------------+---+
			16 | TT.Content         | "1"                |   |
			---+--------------------+--------------------+---+
			17 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			18 | TT.Content         | "+"                |   |
			---+--------------------+--------------------+---+
			19 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			20 | TT.Content         | "2"                |   |
			---+--------------------+--------------------+---+
			21 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			22 | TT.Content         | "="                |   |
			---+--------------------+--------------------+---+
			23 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			24 | TT.CodeLocalExpr   | "(1 + 2)"          |   |
			---+--------------------+--------------------+---+
			25 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
			26 | TT.Content         | "DONE"             |   |
			---+--------------------+--------------------+---+
			27 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
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
			   | TT                 | S                  | E |
			---+--------------------+--------------------+---+
			 1 | TT.Content         | "Try:"             |   |
			---+--------------------+--------------------+---+
			 2 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
			 3 | TT.CodeGlobalBlock | " import \"fmt\" " |   |
			---+--------------------+--------------------+---+
			 4 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
			 5 | TT.CodeLocalBlock  | " v := 1 "         |   |
			---+--------------------+--------------------+---+
			 6 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
			 7 | TT.Content         | "v"                |   |
			---+--------------------+--------------------+---+
			 8 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			 9 | TT.Content         | "="                |   |
			---+--------------------+--------------------+---+
			10 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			11 | TT.CodeLocalExpr   | "v"                |   |
			---+--------------------+--------------------+---+
			12 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			13 | TT.Content         | "^"                |   |
			---+--------------------+--------------------+---+
			14 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
			15 | TT.WS              | "\t  "             |   |
			---+--------------------+--------------------+---+
			16 | TT.Content         | "1"                |   |
			---+--------------------+--------------------+---+
			17 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			18 | TT.Content         | "+"                |   |
			---+--------------------+--------------------+---+
			19 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			20 | TT.Content         | "2"                |   |
			---+--------------------+--------------------+---+
			21 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			22 | TT.Content         | "="                |   |
			---+--------------------+--------------------+---+
			23 | TT.WS              | " "                |   |
			---+--------------------+--------------------+---+
			24 | TT.CodeLocalExpr   | "(1 + 2)"          |   |
			---+--------------------+--------------------+---+
			25 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
			26 | TT.Content         | "DONE"             |   |
			---+--------------------+--------------------+---+
			27 | TT.NL              | "\n"               |   |
			---+--------------------+--------------------+---+
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
		tok, rest, _ := readNextToken(t, `
◊^{
	type foo struct{
		a string
	}
}foo`[1:])

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.CodeGlobalBlock
			Token.Slc: test:1 - "\n\ttype foo struct{\n\t\ta string\n\t}\n"
			Token.E: 
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
			   | TT               | S             | E |
			---+------------------+---------------+---+
			 1 | TT.Macro         | "SimpleMacro" |   |
			---+------------------+---------------+---+
			 2 | TT.Content       | "(1 + 2) = "  |   |
			---+------------------+---------------+---+
			 3 | TT.CodeLocalExpr | "(1 + 2)"     |   |
			---+------------------+---------------+---+
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
	tokenizer := NewDefault(mapOfSimpleMacro())
	in := input.NewInput("test", s)

	tokens, err := tokenizer.NextTokens(in)
	if err != nil {
		return nil, "", err
	}

	rest := in.Rest()
	if len(tokens) != 1 {
		t.Fatalf("got len(%d) want len(1): %#v", len(tokens), tokens)
	}
	return tokens[0], rest, err
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
		S  string
		E  *any
	}
	tt := make([]tokensTable, len(tokens))
	for i, toks := range tokens {
		tt[i] = tokensTable{toks.TT, toks.Slc.S, toks.E}
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
	var valTok *token.Token
	valTok, err = ct.ParseGoCodeFromTo(in, token.TTcodeLocalExpr, '(', ')', true)
	if err != nil {
		return nil, err
	}
	contentSlc := input.NewSlice("test", fmt.Sprintf("%s = ", valTok.Slc.S), valTok.Slc.Start, valTok.Slc.End)
	contentToken := token.NewToken(token.TTcontent, contentSlc)
	return []*token.Token{contentToken, valTok}, nil
}

func (m simpleMacro) Parse(_ interfaces.TemplateHandler, toks []*token.Token) (rest []*token.Token, err error) {
	return toks, nil
}
