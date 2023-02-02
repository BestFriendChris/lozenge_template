package tokenizer

import (
	"fmt"
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
	"github.com/BestFriendChris/lozenge_template/input"
	"github.com/BestFriendChris/lozenge_template/internal/logic/token"
)

func TestOptimize_noTrimSpaces(t *testing.T) {
	t.Run("multiple content blocks", func(t *testing.T) {
		i := input.NewInput("test", "foo-bar")
		toks := []*token.Token{
			contentToken(i, "foo"),
			contentToken(i, "-"),
			contentToken(i, "bar"),
		}
		c := ic.New(t)
		optimized := Optimize(toks, false)
		c.PT(optimized)
		c.Expect(`
			   | TT         | Slc                | E |
			---+------------+--------------------+---+
			 1 | TT.Content | test:1 - "foo-bar" |   |
			---+------------+--------------------+---+
			`)
	})
	t.Run("multiple content blocks splits lines longer than 60", func(t *testing.T) {
		i := input.NewInput("test", "123456789012345678901234567890123456789012345678901234567890-bar")
		toks := []*token.Token{
			contentToken(i, "1234567890123456789012345678901234567890"),
			contentToken(i, "12345678901234567890"),
			contentToken(i, "-"),
			contentToken(i, "bar"),
		}
		c := ic.New(t)
		optimized := Optimize(toks, false)
		c.PT(optimized)
		c.Expect(`
			   | TT         | Slc                                                                     | E |
			---+------------+-------------------------------------------------------------------------+---+
			 1 | TT.Content | test:1 - "123456789012345678901234567890123456789012345678901234567890" |   |
			---+------------+-------------------------------------------------------------------------+---+
			 2 | TT.Content | test:1 - "-bar"                                                         |   |
			---+------------+-------------------------------------------------------------------------+---+
			`)
	})
	t.Run("multiple content with ws and nl blocks", func(t *testing.T) {
		i := input.NewInput("test", "  foo\n\tbar")
		toks := []*token.Token{
			wsToken(i, "  "),
			contentToken(i, "foo"),
			nlToken(i, "\n"),
			wsToken(i, "\t"),
			contentToken(i, "bar"),
		}
		c := ic.New(t)
		optimized := Optimize(toks, false)
		c.PT(optimized)
		c.Expect(`
			   | TT         | Slc                | E |
			---+------------+--------------------+---+
			 1 | TT.Content | test:1 - "  foo\n" |   |
			---+------------+--------------------+---+
			 2 | TT.Content | test:2 - "\tbar"   |   |
			---+------------+--------------------+---+
			`)
	})
	t.Run("multiple code blocks", func(t *testing.T) {
		i := input.NewInput("test", `
  import "foo"
import "bar"
	what
val := "bar"
if val == "bar" {bar
`[1:])
		toks := []*token.Token{
			wsToken(i, "  "),
			mkToken(token.TTcodeGlobalBlock, i, `import "foo"`),
			nlToken(i, "\n"),
			mkToken(token.TTcodeGlobalBlock, i, `import "bar"`),
			nlToken(i, "\n"),
			wsToken(i, "\t"),
			contentToken(i, "what"),
			nlToken(i, "\n"),
			mkToken(token.TTcodeLocalBlock, i, `val := "bar"`),
			nlToken(i, "\n"),
			mkToken(token.TTcodeLocalBlock, i, `if val == "bar" {`),
			contentToken(i, "bar"),
		}
		c := ic.New(t)
		optimized := Optimize(toks, false)
		c.PT(optimized)
		c.Expect(`
			   | TT                 | Slc                            | E |
			---+--------------------+--------------------------------+---+
			 1 | TT.Content         | test:1 - "  "                  |   |
			---+--------------------+--------------------------------+---+
			 2 | TT.CodeGlobalBlock | test:1 - "import \"foo\""      |   |
			---+--------------------+--------------------------------+---+
			 3 | TT.Content         | test:1 - "\n"                  |   |
			---+--------------------+--------------------------------+---+
			 4 | TT.CodeGlobalBlock | test:2 - "import \"bar\""      |   |
			---+--------------------+--------------------------------+---+
			 5 | TT.Content         | test:2 - "\n"                  |   |
			---+--------------------+--------------------------------+---+
			 6 | TT.Content         | test:3 - "\twhat\n"            |   |
			---+--------------------+--------------------------------+---+
			 7 | TT.CodeLocalBlock  | test:4 - "val := \"bar\""      |   |
			---+--------------------+--------------------------------+---+
			 8 | TT.Content         | test:4 - "\n"                  |   |
			---+--------------------+--------------------------------+---+
			 9 | TT.CodeLocalBlock  | test:5 - "if val == \"bar\" {" |   |
			---+--------------------+--------------------------------+---+
			10 | TT.Content         | test:5 - "bar"                 |   |
			---+--------------------+--------------------------------+---+
			`)
	})
	t.Run("multiple content code blocks", func(t *testing.T) {
		i := input.NewInput("test", `
  foo
import "foo"
val := bar
	bar = val
`[1:])
		toks := []*token.Token{
			wsToken(i, "  "),
			contentToken(i, "foo"),
			nlToken(i, "\n"),
			mkToken(token.TTcodeGlobalBlock, i, `import "foo"`),
			nlToken(i, "\n"),
			mkToken(token.TTcodeLocalBlock, i, `val := bar`),
			nlToken(i, "\n"),
			wsToken(i, "\t"),
			contentToken(i, "bar"),
			wsToken(i, " "),
			contentToken(i, "="),
			wsToken(i, " "),
			mkToken(token.TTcodeLocalExpr, i, `val`),
		}
		c := ic.New(t)
		optimized := Optimize(toks, false)
		c.PT(optimized)
		c.Expect(`
			   | TT                 | Slc                       | E |
			---+--------------------+---------------------------+---+
			 1 | TT.Content         | test:1 - "  foo\n"        |   |
			---+--------------------+---------------------------+---+
			 2 | TT.CodeGlobalBlock | test:2 - "import \"foo\"" |   |
			---+--------------------+---------------------------+---+
			 3 | TT.Content         | test:2 - "\n"             |   |
			---+--------------------+---------------------------+---+
			 4 | TT.CodeLocalBlock  | test:3 - "val := bar"     |   |
			---+--------------------+---------------------------+---+
			 5 | TT.Content         | test:3 - "\n"             |   |
			---+--------------------+---------------------------+---+
			 6 | TT.Content         | test:4 - "\tbar = "       |   |
			---+--------------------+---------------------------+---+
			 7 | TT.CodeLocalExpr   | test:4 - "val"            |   |
			---+--------------------+---------------------------+---+
			`)
	})
	t.Run("multiple content code blocks and macros", func(t *testing.T) {
		i := input.NewInput("test", `
  foo
import "foo"
val := "bar"
	if val == "bar" {
bar = val
}
`[1:])
		toks := []*token.Token{
			wsToken(i, "  "),
			contentToken(i, "foo"),
			nlToken(i, "\n"),
			mkToken(token.TTcodeGlobalBlock, i, `import "foo"`),
			nlToken(i, "\n"),
			mkToken(token.TTcodeLocalBlock, i, `val := "bar"`),
			nlToken(i, "\n"),
			wsToken(i, "\t"),
			mkTokenSlice(token.TTmacro, i.SliceOffset(2)), // if
			mkToken(token.TTcodeLocalBlock, i, `if val == "bar" {`),
			nlToken(i, "\n"),
			contentToken(i, "bar"),
			wsToken(i, " "),
			contentToken(i, "="),
			wsToken(i, " "),
			mkToken(token.TTcodeLocalExpr, i, `val`),
			nlToken(i, "\n"),
			mkToken(token.TTcodeLocalBlock, i, `}`),
		}
		c := ic.New(t)
		optimized := Optimize(toks, false)
		c.PT(optimized)
		c.Expect(`
			   | TT                 | Slc                            | E |
			---+--------------------+--------------------------------+---+
			 1 | TT.Content         | test:1 - "  foo\n"             |   |
			---+--------------------+--------------------------------+---+
			 2 | TT.CodeGlobalBlock | test:2 - "import \"foo\""      |   |
			---+--------------------+--------------------------------+---+
			 3 | TT.Content         | test:2 - "\n"                  |   |
			---+--------------------+--------------------------------+---+
			 4 | TT.CodeLocalBlock  | test:3 - "val := \"bar\""      |   |
			---+--------------------+--------------------------------+---+
			 5 | TT.Content         | test:3 - "\n"                  |   |
			---+--------------------+--------------------------------+---+
			 6 | TT.Content         | test:4 - "\t"                  |   |
			---+--------------------+--------------------------------+---+
			 7 | TT.Macro           | test:4 - "if"                  |   |
			---+--------------------+--------------------------------+---+
			 8 | TT.CodeLocalBlock  | test:4 - "if val == \"bar\" {" |   |
			---+--------------------+--------------------------------+---+
			 9 | TT.Content         | test:4 - "\n"                  |   |
			---+--------------------+--------------------------------+---+
			10 | TT.Content         | test:5 - "bar = "              |   |
			---+--------------------+--------------------------------+---+
			11 | TT.CodeLocalExpr   | test:5 - "val"                 |   |
			---+--------------------+--------------------------------+---+
			12 | TT.Content         | test:5 - "\n"                  |   |
			---+--------------------+--------------------------------+---+
			13 | TT.CodeLocalBlock  | test:6 - "}"                   |   |
			---+--------------------+--------------------------------+---+
			`)
	})
}

func TestOptimize_trimSpaces(t *testing.T) {
	t.Run("multiple code blocks", func(t *testing.T) {
		i := input.NewInput("test", `
  import "foo"
import "bar"
	val := "bar"
if val == "bar" {

bar
`[1:])
		toks := []*token.Token{
			wsToken(i, "  "),
			mkToken(token.TTcodeGlobalBlock, i, `import "foo"`),
			nlToken(i, "\n"),
			mkToken(token.TTcodeGlobalBlock, i, `import "bar"`),
			nlToken(i, "\n"),
			wsToken(i, "\t"),
			mkToken(token.TTcodeLocalBlock, i, `val := "bar"`),
			nlToken(i, "\n"),
			mkToken(token.TTcodeLocalBlock, i, `if val == "bar" {`),
			nlToken(i, "\n"),
			nlToken(i, "\n"),
			contentToken(i, "bar"),
		}
		c := ic.New(t)
		optimized := Optimize(toks, true)
		c.PT(optimized)
		c.Expect(`
			   | TT                 | Slc                            | E |
			---+--------------------+--------------------------------+---+
			 1 | TT.CodeGlobalBlock | test:1 - "import \"foo\""      |   |
			---+--------------------+--------------------------------+---+
			 2 | TT.CodeGlobalBlock | test:2 - "import \"bar\""      |   |
			---+--------------------+--------------------------------+---+
			 3 | TT.CodeLocalBlock  | test:3 - "val := \"bar\""      |   |
			---+--------------------+--------------------------------+---+
			 4 | TT.CodeLocalBlock  | test:4 - "if val == \"bar\" {" |   |
			---+--------------------+--------------------------------+---+
			 5 | TT.Content         | test:5 - "\n"                  |   |
			---+--------------------+--------------------------------+---+
			 6 | TT.Content         | test:6 - "bar"                 |   |
			---+--------------------+--------------------------------+---+
			`)
	})

	t.Run("multiple content code blocks", func(t *testing.T) {
		i := input.NewInput("test", `
  foo
  import "foo"
bar
    	val := "bar"
	bar = val
`[1:])
		toks := []*token.Token{
			wsToken(i, "  "),
			contentToken(i, "foo"),
			nlToken(i, "\n"),
			wsToken(i, "  "), // will trim
			mkToken(token.TTcodeGlobalBlock, i, `import "foo"`),
			nlToken(i, "\n"), // will trim
			contentToken(i, "bar"),
			nlToken(i, "\n"),
			wsToken(i, "    "), // will trim
			wsToken(i, "\t"),   // will trim
			mkToken(token.TTcodeLocalBlock, i, `val := "bar"`),
			nlToken(i, "\n"), // will trim
			wsToken(i, "\t"),
			contentToken(i, "bar"),
			wsToken(i, " "),
			contentToken(i, "="),
			wsToken(i, " "), // will NOT trim
			mkToken(token.TTcodeLocalExpr, i, `val`),
			nlToken(i, "\n"), // will NOT trim
		}
		c := ic.New(t)
		optimized := Optimize(toks, true)
		c.PT(optimized)
		c.Expect(`
			   | TT                 | Slc                       | E |
			---+--------------------+---------------------------+---+
			 1 | TT.Content         | test:1 - "  foo\n"        |   |
			---+--------------------+---------------------------+---+
			 2 | TT.CodeGlobalBlock | test:2 - "import \"foo\"" |   |
			---+--------------------+---------------------------+---+
			 3 | TT.Content         | test:3 - "bar\n"          |   |
			---+--------------------+---------------------------+---+
			 4 | TT.CodeLocalBlock  | test:4 - "val := \"bar\"" |   |
			---+--------------------+---------------------------+---+
			 5 | TT.Content         | test:5 - "\tbar = "       |   |
			---+--------------------+---------------------------+---+
			 6 | TT.CodeLocalExpr   | test:5 - "val"            |   |
			---+--------------------+---------------------------+---+
			 7 | TT.Content         | test:5 - "\n"             |   |
			---+--------------------+---------------------------+---+
			`)
	})
	t.Run("multiple content code blocks and macros", func(t *testing.T) {
		i := input.NewInput("test", `
  foo
import "foo"
val := "bar"
	if val == "bar" {
bar = val
}
`[1:])
		toks := []*token.Token{
			wsToken(i, "  "),
			contentToken(i, "foo"),
			nlToken(i, "\n"),
			mkToken(token.TTcodeGlobalBlock, i, `import "foo"`),
			nlToken(i, "\n"),
			mkToken(token.TTcodeLocalBlock, i, `val := "bar"`),
			nlToken(i, "\n"),
			wsToken(i, "\t"),
			mkTokenSlice(token.TTmacro, i.SliceOffset(2)),
			mkToken(token.TTcodeLocalBlock, i, `if val == "bar" {`),
			nlToken(i, "\n"),
			contentToken(i, "bar"),
			wsToken(i, " "),
			contentToken(i, "="),
			wsToken(i, " "),
			mkToken(token.TTcodeLocalExpr, i, `val`),
			nlToken(i, "\n"),
			mkToken(token.TTcodeLocalBlock, i, `}`),
		}
		c := ic.New(t)
		optimized := Optimize(toks, true)
		c.PT(optimized)
		c.Expect(`
			   | TT                 | Slc                            | E |
			---+--------------------+--------------------------------+---+
			 1 | TT.Content         | test:1 - "  foo\n"             |   |
			---+--------------------+--------------------------------+---+
			 2 | TT.CodeGlobalBlock | test:2 - "import \"foo\""      |   |
			---+--------------------+--------------------------------+---+
			 3 | TT.CodeLocalBlock  | test:3 - "val := \"bar\""      |   |
			---+--------------------+--------------------------------+---+
			 4 | TT.Macro           | test:4 - "if"                  |   |
			---+--------------------+--------------------------------+---+
			 5 | TT.CodeLocalBlock  | test:4 - "if val == \"bar\" {" |   |
			---+--------------------+--------------------------------+---+
			 6 | TT.Content         | test:5 - "bar = "              |   |
			---+--------------------+--------------------------------+---+
			 7 | TT.CodeLocalExpr   | test:5 - "val"                 |   |
			---+--------------------+--------------------------------+---+
			 8 | TT.Content         | test:5 - "\n"                  |   |
			---+--------------------+--------------------------------+---+
			 9 | TT.CodeLocalBlock  | test:6 - "}"                   |   |
			---+--------------------+--------------------------------+---+
			`)
	})
}

func contentToken(i *input.Input, prefix string) *token.Token {
	return mkToken(token.TTcontent, i, prefix)
}

func nlToken(i *input.Input, prefix string) *token.Token {
	return mkToken(token.TTnl, i, prefix)
}

func wsToken(i *input.Input, prefix string) *token.Token {
	return mkToken(token.TTws, i, prefix)
}

func mkToken(tt token.TokenType, i *input.Input, prefix string) *token.Token {
	slc, found := i.ConsumeString(prefix)
	if !found {
		panic(fmt.Sprintf("prefix not found: %q\ninput: %q", prefix, i.Rest()))
	}
	return mkTokenSlice(tt, slc)
}

func mkTokenSlice(tt token.TokenType, slc input.Slice) *token.Token {
	return &token.Token{TT: tt, Slc: slc}
}
