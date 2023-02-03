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
		in := input.NewInput("test", `
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
			/* line test:1 */ input "fmt"
			
			################################################################################
			# LOCAL
			################################################################################
			/* line test:2 */ val := 1
			/* line test:3 */ fmt.Print("foo = ")
			/* line test:4 */ fmt.Printf("%v", val)
			/* line test:4 */ fmt.Print("\n")
			`)
	})
	t.Run("with macros", func(t *testing.T) {
		in := input.NewInput("test", `
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
			/* line test:1 */ input "fmt"
			
			################################################################################
			# LOCAL
			################################################################################
			/* line test:2 */ val := 1
			/* line test:3 */ if val == 1 {
			/* line test:4 */ fmt.Print("foo = ")
			/* line test:5 */ fmt.Printf("%v", val)
			/* line test:6 */ fmt.Print("\n")
			/* line test:7 */ }
			`)
	})
}

type testHandler struct {
	GlobalOutput, LocalOutput strings.Builder
}

func (h *testHandler) DefaultMacros() *interfaces.Macros {
	return nil
}

func (h *testHandler) WriteTextContent(slc input.Slice) {
	_, _ = fmt.Fprintf(&h.LocalOutput, "%s fmt.Print(%q)\n", line(slc), slc.S)
}

func (h *testHandler) WriteCodeLocalExpression(slc input.Slice) {
	_, _ = fmt.Fprintf(&h.LocalOutput, "%s fmt.Printf(%q, %s)\n", line(slc), `%v`, slc.S)
}

func (h *testHandler) WriteCodeLocalBlock(slc input.Slice) {
	s := fmt.Sprintf("%s %s", line(slc), slc.S)
	_, _ = fmt.Fprintln(&h.LocalOutput, s)
}

func (h *testHandler) WriteCodeGlobalBlock(slc input.Slice) {
	s := fmt.Sprintf("%s %s", line(slc), slc.S)
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

func line(slc input.Slice) string {
	return fmt.Sprintf("/* line %s:%d */", slc.Name, slc.Start.Row)
}
