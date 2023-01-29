package errors

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
)

func TestTokenizerError_ShowFailureContext(t *testing.T) {
	t.Run("on a single line", func(t *testing.T) {
		input := "foo ◊(1 + 2 bar"

		e := fmt.Errorf("did not find matched ')'")
		idx := strings.Index(input, "(1")
		e = NewTokenizerError(input, idx, e)

		c := ic.New(t)
		c.PrintSection("error with context")
		c.Println(e)
		c.Expect(`
			################################################################################
			# error with context
			################################################################################
			line 1: foo ◊(1 + 2 bar
			             ▲
			             └── did not find matched ')'
			`)
	})
	t.Run("on multiple line", func(t *testing.T) {
		input := `
START:
◊{if 1 == 2 }◊
nope
◊} else 
will fail
◊}
`[1:]
		e := fmt.Errorf("no open brace found")
		idx := strings.Index(input, "} else")
		e = NewTokenizerError(input, idx, e)

		c := ic.New(t)
		c.PrintSection("error with context")
		c.Println(e)
		c.Expect(`
			################################################################################
			# error with context
			################################################################################
			line 4: ◊} else 
			         ▲
			         └── no open brace found
			`)
	})
}

func Test_findLine(t *testing.T) {
	t.Run("one line", func(t *testing.T) {
		lineNo, line, idx := findLine(`only one line`, 5)
		c := ic.New(t)
		c.PrintSection("line")
		c.Printf("line %d: %s\n", lineNo, line)

		c.PrintSection("idx")
		c.Println(idx)

		c.Expect(`
			################################################################################
			# line
			################################################################################
			line 1: only one line
			################################################################################
			# idx
			################################################################################
			5
			`)
	})
	t.Run("two lines", func(t *testing.T) {
		input := `
this line is one
this line is two
`[1:]
		globalIdx := strings.Index(input, "two")
		lineNo, line, idx := findLine(input, globalIdx)
		c := ic.New(t)
		c.PrintSection("line")
		c.Printf("line %d: %s\n", lineNo, line)

		c.PrintSection("idx")
		c.Println(idx)

		c.Expect(`
			################################################################################
			# line
			################################################################################
			line 2: this line is two
			################################################################################
			# idx
			################################################################################
			13
			`)
	})
}
