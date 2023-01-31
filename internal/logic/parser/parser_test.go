package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
	"github.com/BestFriendChris/lozenge_template/input"
	"github.com/BestFriendChris/lozenge_template/interfaces"
	"github.com/BestFriendChris/lozenge_template/internal/logic/macro/macro_if"
	"github.com/BestFriendChris/lozenge_template/internal/logic/token"
)

func TestParser_Parse(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		in := input.NewInput(`
input "fmt"
val := 1
foo = 
val

`[1:])
		toks := []*token.Token{
			token.NewToken(token.TTcodeGlobalBlock, in.SliceAt(0, 11)),
			token.NewToken(token.TTcodeLocalBlock, in.SliceAt(12, 20)),
			token.NewToken(token.TTcontent, in.SliceAt(21, 27)),
			token.NewToken(token.TTcodeLocalExpr, in.SliceAt(28, 31)),
			token.NewToken(token.TTcontent, in.SliceAt(31, 32)),
		}
		h := &testHandler{}
		_, _ = New(nil).Parse(h, toks)
		s, _ := h.Done()

		c := ic.New(t)
		c.Print(s)
		c.Expect(`
			################################################################################
			# GLOBAL
			################################################################################
			input "fmt"
			
			################################################################################
			# LOCAL
			################################################################################
			val := 1
			fmt.Print("foo = ")
			fmt.Printf("%v", val)
			fmt.Print("\n")
			`)
	})
	t.Run("with macros", func(t *testing.T) {
		in := input.NewInput(`
input "fmt"
val := 1
if val == 1 {
foo = 
val

}`[1:])

		toks := []*token.Token{
			token.NewToken(token.TTcodeGlobalBlock, in.SliceAt(0, 11)),
			token.NewToken(token.TTcodeLocalBlock, in.SliceAt(12, 20)),
			token.NewToken(token.TTmacro, in.SliceAt(21, 23)),
			token.NewToken(token.TTcodeLocalBlock, in.SliceAt(21, 34)),
			token.NewToken(token.TTcontent, in.SliceAt(35, 41)),
			token.NewToken(token.TTcodeLocalExpr, in.SliceAt(42, 45)),
			token.NewToken(token.TTcontent, in.SliceAt(46, 47)),
			token.NewToken(token.TTcodeLocalBlock, in.SliceAt(47, 48)),
		}
		h := &testHandler{}
		macros := interfaces.NewMacros()
		macros.Add(macro_if.New())
		_, err := New(macros).Parse(h, toks)
		if err != nil {
			t.Fatal(err)
		}
		s, _ := h.Done()

		c := ic.New(t)
		c.Print(s)
		c.Expect(`
			################################################################################
			# GLOBAL
			################################################################################
			input "fmt"
			
			################################################################################
			# LOCAL
			################################################################################
			val := 1
			if val == 1 {
			fmt.Print("foo = ")
			fmt.Printf("%v", val)
			fmt.Print("\n")
			}
			`)
	})
}

type testHandler struct {
	GlobalOutput, LocalOutput strings.Builder
}

func (h *testHandler) DefaultMacros() *interfaces.Macros {
	return nil
}

func (h *testHandler) WriteTextContent(s string) {
	_, _ = fmt.Fprintf(&h.LocalOutput, "fmt.Print(%q)\n", s)
}

func (h *testHandler) WriteCodeLocalExpression(s string) {
	_, _ = fmt.Fprintf(&h.LocalOutput, "fmt.Printf(%q, %s)\n", `%v`, s)
}

func (h *testHandler) WriteCodeLocalBlock(s string) {
	_, _ = fmt.Fprintln(&h.LocalOutput, s)
}

func (h *testHandler) WriteCodeGlobalBlock(s string) {
	_, _ = fmt.Fprintln(&h.GlobalOutput, s)
}

func (h *testHandler) Done() (string, error) {
	var sb strings.Builder
	_, _ = fmt.Fprintln(&sb, "################################################################################")
	_, _ = fmt.Fprintln(&sb, "# GLOBAL")
	_, _ = fmt.Fprintln(&sb, "################################################################################")
	_, _ = fmt.Fprint(&sb, h.GlobalOutput.String())

	_, _ = fmt.Fprintln(&sb)
	_, _ = fmt.Fprintln(&sb, "################################################################################")
	_, _ = fmt.Fprintln(&sb, "# LOCAL")
	_, _ = fmt.Fprintln(&sb, "################################################################################")
	_, _ = fmt.Fprint(&sb, h.LocalOutput.String())
	return sb.String(), nil
}
