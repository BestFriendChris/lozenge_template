package tokenizer

import (
	"strings"
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
	"github.com/BestFriendChris/lozenge/interfaces"
	"github.com/BestFriendChris/lozenge/internal/logic/token"
)

func TestContentTokenizer_ReadTokensUntil(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		input := `
Try:
◊^{ import "fmt" }
◊{ v := 1 }
v = ◊v ◊
	  1 + 2 = ◊(1 + 2)
DONE
`[1:]
		c := ic.New(t)
		c.PrintSection("Input")
		c.Println(input)

		tokenizer := NewDefault(interfaces.NewMacros())
		tokens, rest := tokenizer.ReadTokensUntil(input, "DONE")

		c.PrintSection("Tokens")
		c.PT(tokens)

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
			# Tokens
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
		input := `Hi: ◊.SimpleMacro(1 + 2) DONE`
		c := ic.New(t)
		c.PrintSection("Input")
		c.Println(input)

		tokenizer := NewDefault(mapOfSimpleMacro())
		tokens, rest := tokenizer.ReadTokensUntil(input, "DONE")

		c.PrintSection("Tokens")
		c.PT(tokens)

		c.PrintSection("Rest")
		c.Println(rest)

		c.Expect(`
			################################################################################
			# Input
			################################################################################
			Hi: ◊.SimpleMacro(1 + 2) DONE
			################################################################################
			# Tokens
			################################################################################
			   | TT               | S             | E |
			---+------------------+---------------+---+
			 1 | TT.Content       | "Hi:"         |   |
			---+------------------+---------------+---+
			 2 | TT.WS            | " "           |   |
			---+------------------+---------------+---+
			 3 | TT.Macro         | "SimpleMacro" |   |
			---+------------------+---------------+---+
			 4 | TT.CodeLocalExpr | "(1 + 2)"     |   |
			---+------------------+---------------+---+
			 5 | TT.WS            | " "           |   |
			---+------------------+---------------+---+
			################################################################################
			# Rest
			################################################################################
			DONE
			`)
	})
	t.Run("stop at end of stream", func(t *testing.T) {
		input := `Hi: ◊.SimpleMacro(1 + 2) DONE`
		c := ic.New(t)
		c.Println(`ReadTokensUntil(input, "") will read to end of stream`)
		c.PrintSection("Input")
		c.Println(input)

		tokenizer := NewDefault(mapOfSimpleMacro())
		tokens, rest := tokenizer.ReadTokensUntil(input, "")

		c.PrintSection("Tokens")
		c.PT(tokens)

		c.PrintSection("Rest")
		c.Println(rest)

		c.Expect(`
			ReadTokensUntil(input, "") will read to end of stream
			################################################################################
			# Input
			################################################################################
			Hi: ◊.SimpleMacro(1 + 2) DONE
			################################################################################
			# Tokens
			################################################################################
			   | TT               | S             | E |
			---+------------------+---------------+---+
			 1 | TT.Content       | "Hi:"         |   |
			---+------------------+---------------+---+
			 2 | TT.WS            | " "           |   |
			---+------------------+---------------+---+
			 3 | TT.Macro         | "SimpleMacro" |   |
			---+------------------+---------------+---+
			 4 | TT.CodeLocalExpr | "(1 + 2)"     |   |
			---+------------------+---------------+---+
			 5 | TT.WS            | " "           |   |
			---+------------------+---------------+---+
			 6 | TT.Content       | "DONE"        |   |
			---+------------------+---------------+---+
			################################################################################
			# Rest
			################################################################################
			
			`)
	})
	t.Run("undefined macro", func(t *testing.T) {
		input := `Hi: ◊.UndefinedMacro(1) DONE`
		c := ic.New(t)
		c.PrintSection("Input")
		c.Println(input)

		tokenizer := NewDefault(interfaces.NewMacros())
		tokens, rest := tokenizer.ReadTokensUntil(input, "DONE")

		c.PrintSection("Tokens")
		c.PT(tokens)

		c.PrintSection("Rest")
		c.Println(rest)

		c.Expect(`
			################################################################################
			# Input
			################################################################################
			Hi: ◊.UndefinedMacro(1) DONE
			################################################################################
			# Tokens
			################################################################################
			   | TT         | S                    | E |
			---+------------+----------------------+---+
			 1 | TT.Content | "Hi:"                |   |
			---+------------+----------------------+---+
			 2 | TT.WS      | " "                  |   |
			---+------------+----------------------+---+
			 3 | TT.Content | "◊"                  |   |
			---+------------+----------------------+---+
			 4 | TT.Content | ".UndefinedMacro(1)" |   |
			---+------------+----------------------+---+
			 5 | TT.WS      | " "                  |   |
			---+------------+----------------------+---+
			################################################################################
			# Rest
			################################################################################
			DONE
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
			toks, rest = tokenizer.NextTokens(rest)
			for _, tok := range toks {
				tokens = append(tokens, tok)
			}
			if rest == "" {
				break
			}
		}

		c.PrintSection("Tokens")
		c.PT(tokens)
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
			# Tokens
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
			toks, rest = tokenizer.NextTokens(rest)
			for _, tok := range toks {
				tokens = append(tokens, tok)
			}
			if rest == "" {
				break
			}
		}

		c.PrintSection("Tokens")
		c.PT(tokens)
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
			# Tokens
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
		tok, rest := readNextToken(t, "\t   hi")

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.WS
			Token.S: "\t   "
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"hi"
			`)
	})
	t.Run("newline", func(t *testing.T) {
		tok, rest := readNextToken(t, "\n\nfoo")

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
			Token.S: "\n"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"\nfoo"
			`)
	})
	t.Run("content", func(t *testing.T) {
		tok, rest := readNextToken(t, "foo\nbar")

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.Content
			Token.S: "foo"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"\nbar"
			`)
	})
	t.Run("◊◊foo", func(t *testing.T) {
		tok, rest := readNextToken(t, `◊◊foo`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.Content
			Token.S: "◊"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"foo"
			`)
	})
	t.Run("◊ bar", func(t *testing.T) {
		tok, rest := readNextToken(t, `◊ bar`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.Content
			Token.S: "◊"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			" bar"
			`)
	})
	t.Run(`◊\nbar`, func(t *testing.T) {
		tok, rest := readNextToken(t, "◊\nbar")

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.Content
			Token.S: "◊"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"\nbar"
			`)
	})
	t.Run(`◊`, func(t *testing.T) {
		tok, rest := readNextToken(t, "◊")

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.Content
			Token.S: "◊"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			""
			`)
	})
	t.Run("◊foo bar", func(t *testing.T) {
		tok, rest := readNextToken(t, `◊foo bar`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.CodeLocalExpr
			Token.S: "foo"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			" bar"
			`)
	})
	t.Run("◊foo", func(t *testing.T) {
		tok, rest := readNextToken(t, `◊foo`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.CodeLocalExpr
			Token.S: "foo"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			""
			`)
	})
	t.Run("◊(1 + 2)foo", func(t *testing.T) {
		tok, rest := readNextToken(t, `◊(1 + 2)foo`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.CodeLocalExpr
			Token.S: "(1 + 2)"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"foo"
			`)
	})
	t.Run("◊(1 + 2", func(t *testing.T) {
		tok, rest := readNextToken(t, `◊(1 + 2`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.Content
			Token.S: "◊"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"(1 + 2"
			`)
	})
	t.Run("◊{ GOCODE }", func(t *testing.T) {
		tok, rest := readNextToken(t, `◊{ var foo, bar, baz := struct{a string}{"}\""}, '}', '\'' }foo`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.CodeLocalBlock
			Token.S: " var foo, bar, baz := struct{a string}{\"}\\\"\"}, '}', '\\'' "
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"foo"
			`)
	})
	t.Run("◊{ GOCODE ", func(t *testing.T) {
		tok, rest := readNextToken(t, `◊{ var foo struct{a string}{"}"} foo`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.Content
			Token.S: "◊"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"{ var foo struct{a string}{\"}\"} foo"
			`)
	})
	t.Run("◊^{ GOCODE }", func(t *testing.T) {
		tok, rest := readNextToken(t, `
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
			Token.S: "\n\ttype foo struct{\n\t\ta string\n\t}\n"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"foo"
			`)
	})
	t.Run("◊^{ GOCODE ", func(t *testing.T) {
		tok, rest := readNextToken(t, `
◊^{
	type foo struct{
		a string
	}
foo`[1:])

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.Content
			Token.S: "◊"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"{\n\ttype foo struct{\n\t\ta string\n\t}\nfoo"
			`)
	})
	t.Run("◊^foo", func(t *testing.T) {
		tok, rest := readNextToken(t, `◊^foo`)

		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.Content
			Token.S: "◊"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"^foo"
			`)
	})
	t.Run("◊.macro", func(t *testing.T) {
		tokenizer := NewDefault(mapOfSimpleMacro())
		tokens, rest := tokenizer.NextTokens(`◊.SimpleMacro(1 + 2)`)

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
			 2 | TT.CodeLocalExpr | "(1 + 2)"     |   |
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
		tok, rest := tokenizer.NextTokenCodeUntilOpenBraceLoz(`if strings.DeepEqual(v, []string{"\"", "{"}) {◊foo`)
		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			Token.TT: TT.CodeLocalBlock
			Token.S: "if strings.DeepEqual(v, []string{\"\\\"\", \"{\"}) {"
			Token.E: 
			################################################################################
			# rest
			################################################################################
			"foo"
			`)
	})
	t.Run("no open brace", func(t *testing.T) {
		tokenizer := NewDefault(interfaces.NewMacros())
		tok, rest := tokenizer.NextTokenCodeUntilOpenBraceLoz(`if "\"" == "{" ◊} else { foo ◊}`)
		c := ic.New(t)
		logToken(c, tok)
		logRest(c, rest)
		c.Expect(`
			################################################################################
			# token
			################################################################################
			<null>
			################################################################################
			# rest
			################################################################################
			"if \"\\\"\" == \"{\" ◊} else { foo ◊}"
			`)
	})
}

func readNextToken(t *testing.T, s string) (*token.Token, string) {
	tokenizer := NewDefault(mapOfSimpleMacro())
	tokens, rest := tokenizer.NextTokens(s)
	if len(tokens) != 1 {
		t.Fatalf("got len(%d) want len(1): %#v", len(tokens), tokens)
	}
	return tokens[0], rest
}

func logToken(c *ic.IC, t *token.Token) {

	c.PrintSection("token")
	if t == nil {
		c.Println("<null>")
	} else {
		c.PV(t)
	}
}
func logTokens(c *ic.IC, toks []*token.Token) {
	c.PrintSection("tokens")
	c.PT(toks)
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

func (m simpleMacro) NextTokens(ct interfaces.ContentTokenizer, rest string) ([]*token.Token, string) {
	rest = strings.TrimPrefix(rest, m.Name())

	runes := []rune(rest)
	var tok *token.Token
	tok, rest = ct.ParseGoCodeFromTo(runes, token.TTcodeLocalExpr, '(', ')', true)
	return []*token.Token{tok}, rest
}

func (m simpleMacro) Parse(_ interfaces.TemplateHandler, toks []*token.Token) (rest []*token.Token, err error) {
	return toks, nil
}
