package input

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"unicode"

	"github.com/BestFriendChris/go-ic/ic"
)

func TestInput_ReadWhile(t *testing.T) {
	in := NewInput("foo bar")
	ident := in.ReadWhile(func(r rune) bool {
		isLetter := unicode.IsLetter(r) || r == '_'
		return isLetter || unicode.IsNumber(r)
	})

	c := ic.New(t)
	c.PrintSection("Rest")
	c.Printf("%q\n", in.Rest())

	c.PrintSection("ident")
	c.Printf("%q\n", ident)
	c.Expect(`
		################################################################################
		# Rest
		################################################################################
		" bar"
		################################################################################
		# ident
		################################################################################
		"foo"
		`)
}

func TestInput_TryReadWhile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		in := NewInput("(1 + (2 * 3))next")
		var parenCount int
		var foundAll bool
		str, err := in.TryReadWhile(func(r rune, last bool) (bool, error) {
			if foundAll {
				return false, nil
			}
			if r == '(' {
				parenCount++
			} else if r == ')' {
				parenCount--
			}
			if parenCount == 0 {
				foundAll = true
			} else if last && parenCount > 0 {
				return false, fmt.Errorf("unbalanced parens")
			}
			return true, nil
		})
		c := ic.New(t)
		c.PrintSection("str")
		c.Printf("%q\n", str)

		c.PrintSection("err")
		c.Printf("%s\n", err)

		c.PrintSection("in.rest")
		c.Printf("%q\n", in.Rest())

		c.Expect(`
			################################################################################
			# str
			################################################################################
			"(1 + (2 * 3))"
			################################################################################
			# err
			################################################################################
			%!s(<nil>)
			################################################################################
			# in.rest
			################################################################################
			"next"
			`)
	})
	t.Run("error path", func(t *testing.T) {
		in := NewInput("(1 + (2 * 3)next")
		var parenCount int
		var foundAll bool
		str, err := in.TryReadWhile(func(r rune, last bool) (bool, error) {
			if foundAll {
				return false, nil
			}
			if r == '(' {
				parenCount++
			} else if r == ')' {
				parenCount--
			}
			if parenCount == 0 {
				foundAll = true
			} else if last && parenCount > 0 {
				return false, fmt.Errorf("unbalanced parens")
			}
			return true, nil
		})
		c := ic.New(t)
		c.PrintSection("str")
		c.Printf("%q\n", str)

		c.PrintSection("err")
		c.Printf("%s\n", err)

		c.PrintSection("in.rest")
		c.Printf("%q\n", in.Rest())

		c.Expect(`
			################################################################################
			# str
			################################################################################
			""
			################################################################################
			# err
			################################################################################
			line 1: (1 + (2 * 3)next
			        ▲
			        └── unbalanced parens
			################################################################################
			# in.rest
			################################################################################
			"(1 + (2 * 3)next"
			`)
	})
}

func TestInput_ConsumeWhenMatchesRegexp(t *testing.T) {
	in := NewInput(`}   else {◊ foo`)
	elseRegex := regexp.MustCompile(`}\s*else\s*{`)
	match, found := in.ConsumeRegexp(elseRegex)
	if !found {
		t.Fatal("should have matched regex")
	}
	c := ic.New(t)
	c.PrintSection("match")
	c.Println(match)
	c.PrintSection("rest")
	c.Printf("%q\n", in.Rest())

	c.Expect(`
		################################################################################
		# match
		################################################################################
		}   else {
		################################################################################
		# rest
		################################################################################
		"◊ foo"
		`)
}

func TestInput_Line(t *testing.T) {
	t.Run("single line", func(t *testing.T) {
		s := `only one line`
		in := NewInput(s)
		c := ic.New(t)

		c.PrintSection("start")
		c.Println(in.Line())

		c.PrintSection(`seek to "one"`)
		in.Seek(strings.Index(s, "one"))
		c.Println(in.Line())

		c.PrintSection("seek to end")
		in.Seek(len(s))
		c.Println(in.Line())

		c.Expect(`
			################################################################################
			# start
			################################################################################
			1
			################################################################################
			# seek to "one"
			################################################################################
			1
			################################################################################
			# seek to end
			################################################################################
			1
			`)
	})
	t.Run("multi line", func(t *testing.T) {
		s := `
this is line 1
this is line 2
this is line 3`[1:]
		in := NewInput(s)
		c := ic.New(t)

		c.PrintSection("start")
		c.PVWN("Line", in.Line())
		c.PVWN("Rest", in.Rest())

		c.PrintSection(`seek to first newline`)
		in.Seek(strings.Index(s, "\nthis is line 2"))
		c.PVWN("Line", in.Line())
		c.PVWN("Rest", in.Rest())

		c.PrintSection(`seek to "one"`)
		in.Seek(strings.Index(s, "line 2"))
		c.PVWN("Line", in.Line())
		c.PVWN("Rest", in.Rest())

		c.PrintSection("seek to end")
		in.Seek(len(s))
		c.PVWN("Line", in.Line())
		c.PVWN("Rest", in.Rest())

		c.Expect(`
			################################################################################
			# start
			################################################################################
			Line: 1
			Rest: "this is line 1\nthis is line 2\nthis is line 3"
			################################################################################
			# seek to first newline
			################################################################################
			Line: 1
			Rest: "\nthis is line 2\nthis is line 3"
			################################################################################
			# seek to "one"
			################################################################################
			Line: 2
			Rest: "line 2\nthis is line 3"
			################################################################################
			# seek to end
			################################################################################
			Line: 3
			Rest: ""
			`)
	})
}
