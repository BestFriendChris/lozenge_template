package tokenizer

import (
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
	"github.com/BestFriendChris/lozenge_template/internal/logic/token"
)

func TestOptimize_noTrimSpaces(t *testing.T) {
	t.Run("multiple content blocks", func(t *testing.T) {
		toks := []*token.Token{
			contentToken("foo"),
			contentToken("-"),
			contentToken("bar"),
		}
		c := ic.New(t)
		optimized := Optimize(toks, false)
		c.PT(optimized)
		c.Expect(`
			   | TT         | S         | E |
			---+------------+-----------+---+
			 1 | TT.Content | "foo-bar" |   |
			---+------------+-----------+---+
			`)
	})
	t.Run("multiple content blocks splits lines longer than 60", func(t *testing.T) {
		toks := []*token.Token{
			contentToken("1234567890123456789012345678901234567890"),
			contentToken("12345678901234567890"),
			contentToken("-"),
			contentToken("bar"),
		}
		c := ic.New(t)
		optimized := Optimize(toks, false)
		c.PT(optimized)
		c.Expect(`
			   | TT         | S                                                              | E |
			---+------------+----------------------------------------------------------------+---+
			 1 | TT.Content | "123456789012345678901234567890123456789012345678901234567890" |   |
			---+------------+----------------------------------------------------------------+---+
			 2 | TT.Content | "-bar"                                                         |   |
			---+------------+----------------------------------------------------------------+---+
			`)
	})
	t.Run("multiple content with ws and nl blocks", func(t *testing.T) {
		toks := []*token.Token{
			wsToken("  "),
			contentToken("foo"),
			nlToken("\n"),
			wsToken("\t"),
			contentToken("bar"),
		}
		c := ic.New(t)
		optimized := Optimize(toks, false)
		c.PT(optimized)
		c.Expect(`
			   | TT         | S              | E |
			---+------------+----------------+---+
			 1 | TT.Content | "  foo\n\tbar" |   |
			---+------------+----------------+---+
			`)
	})
	t.Run("multiple code blocks", func(t *testing.T) {
		toks := []*token.Token{
			wsToken("  "),
			mkToken(token.TTcodeGlobalBlock, `import "foo"`),
			nlToken("\n"),
			mkToken(token.TTcodeGlobalBlock, `import "bar"`),
			nlToken("\n"),
			wsToken("\t"),
			mkToken(token.TTcodeLocalBlock, `val := "bar"`),
			nlToken("\n"),
			mkToken(token.TTcodeLocalBlock, `if val == "bar" {`),
			contentToken("bar"),
		}
		c := ic.New(t)
		optimized := Optimize(toks, false)
		c.PT(optimized)
		c.Expect(`
			   | TT                 | S                     | E |
			---+--------------------+-----------------------+---+
			 1 | TT.Content         | "  "                  |   |
			---+--------------------+-----------------------+---+
			 2 | TT.CodeGlobalBlock | "import \"foo\""      |   |
			---+--------------------+-----------------------+---+
			 3 | TT.Content         | "\n"                  |   |
			---+--------------------+-----------------------+---+
			 4 | TT.CodeGlobalBlock | "import \"bar\""      |   |
			---+--------------------+-----------------------+---+
			 5 | TT.Content         | "\n\t"                |   |
			---+--------------------+-----------------------+---+
			 6 | TT.CodeLocalBlock  | "val := \"bar\""      |   |
			---+--------------------+-----------------------+---+
			 7 | TT.Content         | "\n"                  |   |
			---+--------------------+-----------------------+---+
			 8 | TT.CodeLocalBlock  | "if val == \"bar\" {" |   |
			---+--------------------+-----------------------+---+
			 9 | TT.Content         | "bar"                 |   |
			---+--------------------+-----------------------+---+
			`)
	})
	t.Run("multiple content code blocks", func(t *testing.T) {
		toks := []*token.Token{
			wsToken("  "),
			contentToken("foo"),
			mkToken(token.TTcodeGlobalBlock, `import "foo"`),
			nlToken("\n"),
			mkToken(token.TTcodeLocalBlock, `val := "bar"`),
			wsToken("\t"),
			contentToken("bar"),
			wsToken(" "),
			contentToken("="),
			wsToken(" "),
			mkToken(token.TTcodeLocalExpr, `val`),
		}
		c := ic.New(t)
		optimized := Optimize(toks, false)
		c.PT(optimized)
		c.Expect(`
			   | TT                 | S                | E |
			---+--------------------+------------------+---+
			 1 | TT.Content         | "  foo"          |   |
			---+--------------------+------------------+---+
			 2 | TT.CodeGlobalBlock | "import \"foo\"" |   |
			---+--------------------+------------------+---+
			 3 | TT.Content         | "\n"             |   |
			---+--------------------+------------------+---+
			 4 | TT.CodeLocalBlock  | "val := \"bar\"" |   |
			---+--------------------+------------------+---+
			 5 | TT.Content         | "\tbar = "       |   |
			---+--------------------+------------------+---+
			 6 | TT.CodeLocalExpr   | "val"            |   |
			---+--------------------+------------------+---+
			`)
	})
	t.Run("multiple content code blocks and macros", func(t *testing.T) {
		toks := []*token.Token{
			wsToken("  "),
			contentToken("foo"),
			mkToken(token.TTcodeGlobalBlock, `import "foo"`),
			nlToken("\n"),
			mkToken(token.TTcodeLocalBlock, `val := "bar"`),
			wsToken("\t"),
			mkToken(token.TTmacro, `If`),
			mkToken(token.TTcodeLocalBlock, `if val == "bar" {`),
			contentToken("bar"),
			wsToken(" "),
			contentToken("="),
			wsToken(" "),
			mkToken(token.TTcodeLocalExpr, `val`),
			nlToken("\n"),
			mkToken(token.TTcodeLocalBlock, `}`),
		}
		c := ic.New(t)
		optimized := Optimize(toks, false)
		c.PT(optimized)
		c.Expect(`
			   | TT                 | S                     | E |
			---+--------------------+-----------------------+---+
			 1 | TT.Content         | "  foo"               |   |
			---+--------------------+-----------------------+---+
			 2 | TT.CodeGlobalBlock | "import \"foo\""      |   |
			---+--------------------+-----------------------+---+
			 3 | TT.Content         | "\n"                  |   |
			---+--------------------+-----------------------+---+
			 4 | TT.CodeLocalBlock  | "val := \"bar\""      |   |
			---+--------------------+-----------------------+---+
			 5 | TT.Content         | "\t"                  |   |
			---+--------------------+-----------------------+---+
			 6 | TT.Macro           | "If"                  |   |
			---+--------------------+-----------------------+---+
			 7 | TT.CodeLocalBlock  | "if val == \"bar\" {" |   |
			---+--------------------+-----------------------+---+
			 8 | TT.Content         | "bar = "              |   |
			---+--------------------+-----------------------+---+
			 9 | TT.CodeLocalExpr   | "val"                 |   |
			---+--------------------+-----------------------+---+
			10 | TT.Content         | "\n"                  |   |
			---+--------------------+-----------------------+---+
			11 | TT.CodeLocalBlock  | "}"                   |   |
			---+--------------------+-----------------------+---+
			`)
	})
}

func TestOptimize_trimSpaces(t *testing.T) {
	t.Run("multiple code blocks", func(t *testing.T) {
		toks := []*token.Token{
			wsToken("  "),
			mkToken(token.TTcodeGlobalBlock, `import "foo"`),
			nlToken("\n"),
			mkToken(token.TTcodeGlobalBlock, `import "bar"`),
			nlToken("\n"),
			wsToken("\t"),
			mkToken(token.TTcodeLocalBlock, `val := "bar"`),
			mkToken(token.TTcodeLocalBlock, `if val == "bar" {`),
			nlToken("\n"),
			nlToken("\n"),
			contentToken("bar"),
		}
		c := ic.New(t)
		optimized := Optimize(toks, true)
		c.PT(optimized)
		c.Expect(`
			   | TT                 | S                                     | E |
			---+--------------------+---------------------------------------+---+
			 1 | TT.CodeGlobalBlock | "import \"foo\"\nimport \"bar\""      |   |
			---+--------------------+---------------------------------------+---+
			 2 | TT.CodeLocalBlock  | "val := \"bar\"\nif val == \"bar\" {" |   |
			---+--------------------+---------------------------------------+---+
			 3 | TT.Content         | "\nbar"                               |   |
			---+--------------------+---------------------------------------+---+
			`)
	})

	t.Run("multiple content code blocks", func(t *testing.T) {
		toks := []*token.Token{
			wsToken("  "),
			contentToken("foo"),
			wsToken("  "), // will trim
			mkToken(token.TTcodeGlobalBlock, `import "foo"`),
			nlToken("\n"), // will trim
			contentToken("bar"),
			wsToken("    "), // will trim
			wsToken("\t"),   // will trim
			mkToken(token.TTcodeLocalBlock, `val := "bar"`),
			nlToken("\n"), // will trim
			wsToken("\t"),
			contentToken("bar"),
			wsToken(" "),
			contentToken("="),
			wsToken(" "), // will NOT trim
			mkToken(token.TTcodeLocalExpr, `val`),
			nlToken("\n"), // will NOT trim
		}
		c := ic.New(t)
		optimized := Optimize(toks, true)
		c.PT(optimized)
		c.Expect(`
			   | TT                 | S                | E |
			---+--------------------+------------------+---+
			 1 | TT.Content         | "  foo"          |   |
			---+--------------------+------------------+---+
			 2 | TT.CodeGlobalBlock | "import \"foo\"" |   |
			---+--------------------+------------------+---+
			 3 | TT.Content         | "bar"            |   |
			---+--------------------+------------------+---+
			 4 | TT.CodeLocalBlock  | "val := \"bar\"" |   |
			---+--------------------+------------------+---+
			 5 | TT.Content         | "\tbar = "       |   |
			---+--------------------+------------------+---+
			 6 | TT.CodeLocalExpr   | "val"            |   |
			---+--------------------+------------------+---+
			 7 | TT.Content         | "\n"             |   |
			---+--------------------+------------------+---+
			`)
	})
	t.Run("multiple content code blocks and macros", func(t *testing.T) {
		toks := []*token.Token{
			wsToken("  "),
			contentToken("foo"),
			mkToken(token.TTcodeGlobalBlock, `import "foo"`),
			nlToken("\n"),
			mkToken(token.TTcodeLocalBlock, `val := "bar"`),
			wsToken("\t"),
			mkToken(token.TTmacro, `If`),
			mkToken(token.TTcodeLocalBlock, `if val == "bar" {`),
			contentToken("bar"),
			wsToken(" "),
			contentToken("="),
			wsToken(" "),
			mkToken(token.TTcodeLocalExpr, `val`),
			nlToken("\n"),
			mkToken(token.TTcodeLocalBlock, `}`),
		}
		c := ic.New(t)
		optimized := Optimize(toks, true)
		c.PT(optimized)
		c.Expect(`
			   | TT                 | S                     | E |
			---+--------------------+-----------------------+---+
			 1 | TT.Content         | "  foo"               |   |
			---+--------------------+-----------------------+---+
			 2 | TT.CodeGlobalBlock | "import \"foo\""      |   |
			---+--------------------+-----------------------+---+
			 3 | TT.CodeLocalBlock  | "val := \"bar\""      |   |
			---+--------------------+-----------------------+---+
			 4 | TT.Macro           | "If"                  |   |
			---+--------------------+-----------------------+---+
			 5 | TT.CodeLocalBlock  | "if val == \"bar\" {" |   |
			---+--------------------+-----------------------+---+
			 6 | TT.Content         | "bar = "              |   |
			---+--------------------+-----------------------+---+
			 7 | TT.CodeLocalExpr   | "val"                 |   |
			---+--------------------+-----------------------+---+
			 8 | TT.Content         | "\n"                  |   |
			---+--------------------+-----------------------+---+
			 9 | TT.CodeLocalBlock  | "}"                   |   |
			---+--------------------+-----------------------+---+
			`)
	})
}

func contentToken(val string) *token.Token {
	return &token.Token{TT: token.TTcontent, S: val}
}

func nlToken(val string) *token.Token {
	return &token.Token{TT: token.TTnl, S: val}
}

func wsToken(val string) *token.Token {
	return &token.Token{TT: token.TTws, S: val}
}

func mkToken(tt token.TokenType, val string) *token.Token {
	return &token.Token{TT: tt, S: val}
}
