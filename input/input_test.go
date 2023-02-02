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
	in := NewInput("test", "foo bar")
	ident := in.ReadWhile(func(r rune) bool {
		isLetter := unicode.IsLetter(r) || r == '_'
		return isLetter || unicode.IsNumber(r)
	})

	c := ic.New(t)
	c.PrintSection("Rest")
	c.Printf("%q\n", in.Rest())

	c.PrintSection("ident")
	c.Println(ident)
	c.Expect(`
		################################################################################
		# Rest
		################################################################################
		" bar"
		################################################################################
		# ident
		################################################################################
		test:1 - "foo"
		`)
}

func TestInput_TryReadWhile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		in := NewInput("test", "(1 + (2 * 3))next")
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
		c.Println(str)

		c.PrintSection("err")
		c.Printf("%s\n", err)

		c.PrintSection("in.rest")
		c.Printf("%q\n", in.Rest())

		c.Expect(`
			################################################################################
			# str
			################################################################################
			test:1 - "(1 + (2 * 3))"
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
		in := NewInput("test", "(1 + (2 * 3)next")
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
		c.Println(str)

		c.PrintSection("err")
		c.Printf("%s\n", err)

		c.PrintSection("in.rest")
		c.Printf("%q\n", in.Rest())

		c.Expect(`
			################################################################################
			# str
			################################################################################
			<empty slice>
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
	in := NewInput("test", `}   else {◊ foo`)
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
		test:1 - "}   else {"
		################################################################################
		# rest
		################################################################################
		"◊ foo"
		`)
}

func TestInput_Pos(t *testing.T) {
	t.Run("single line", func(t *testing.T) {
		s := `only one line`
		in := NewInput("test", s)
		c := ic.New(t)

		c.PrintSection("start")
		c.Println(in.Pos())

		c.PrintSection(`seek to "one"`)
		in.Seek(strings.Index(s, "one"))
		c.Println(in.Pos())

		c.PrintSection("seek to end")
		in.Seek(len(s))
		c.Println(in.Pos())

		c.Expect(`
			################################################################################
			# start
			################################################################################
			Pos[line=1;col=1]
			################################################################################
			# seek to "one"
			################################################################################
			Pos[line=1;col=6]
			################################################################################
			# seek to end
			################################################################################
			Pos[line=1;col=14]
			`)
	})
	t.Run("multi line", func(t *testing.T) {
		s := `
this is line 1
this is line 2
this is line 3`[1:]
		in := NewInput("test", s)
		c := ic.New(t)

		c.PrintSection("start")
		c.PVWN("Position", in.Pos())
		c.PVWN("Rest", in.Rest())

		c.PrintSection(`seek to first newline`)
		in.Seek(strings.Index(s, "\nthis is line 2"))
		c.PVWN("Position", in.Pos())
		c.PVWN("Rest", in.Rest())

		c.PrintSection(`seek to "one"`)
		in.Seek(strings.Index(s, "line 2"))
		c.PVWN("Position", in.Pos())
		c.PVWN("Rest", in.Rest())

		c.PrintSection("seek to end")
		in.Seek(len(s))
		c.PVWN("Position", in.Pos())
		c.PVWN("Rest", in.Rest())

		c.Expect(`
			################################################################################
			# start
			################################################################################
			Position: Pos[line=1;col=1]
			Rest: "this is line 1\nthis is line 2\nthis is line 3"
			################################################################################
			# seek to first newline
			################################################################################
			Position: Pos[line=1;col=15]
			Rest: "\nthis is line 2\nthis is line 3"
			################################################################################
			# seek to "one"
			################################################################################
			Position: Pos[line=2;col=9]
			Rest: "line 2\nthis is line 3"
			################################################################################
			# seek to end
			################################################################################
			Position: Pos[line=3;col=15]
			Rest: ""
			`)
	})
}

func TestInput_SliceAt(t *testing.T) {
	t.Run("at newline break", func(t *testing.T) {
		i := NewInput("test", "\nfoo")

		c := ic.New(t)
		{
			c.PrintSection("the newline")
			fmt.Printf("i.str[0:1]: %q\n", i.str[0:1])
			slc := i.SliceAt(0, 1)
			c.PVWN("S", slc.S)
			c.PVWN("Start", slc.Start)
			c.PVWN("End", slc.End)
		}
		{
			c.PrintSection("second line")
			fmt.Printf("i.str[1:2]: %q\n", i.str[1:2])
			slc := i.SliceAt(1, 2)
			c.PVWN("S", slc.S)
			c.PVWN("Start", slc.Start)
			c.PVWN("End", slc.End)
		}

		c.Expect(`
			################################################################################
			# the newline
			################################################################################
			S: "\n"
			Start: Pos[line=1;col=1]
			End: Pos[line=2;col=1]
			################################################################################
			# second line
			################################################################################
			S: "f"
			Start: Pos[line=2;col=1]
			End: Pos[line=2;col=2]
			`)
	})
}
