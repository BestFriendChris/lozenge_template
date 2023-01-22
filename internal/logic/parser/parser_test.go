package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
	"github.com/BestFriendChris/lozenge/interfaces"
	"github.com/BestFriendChris/lozenge/internal/logic/token"
)

func TestParser_Parse(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		toks := []*token.Token{
			token.NewToken(token.TTcodeGlobalBlock, `input "fmt"`),
			token.NewToken(token.TTcodeLocalBlock, "val := 1"),
			token.NewToken(token.TTcontent, "foo = "),
			token.NewToken(token.TTcodeLocalExpr, "val"),
			token.NewToken(token.TTcontent, "\n"),
		}
		h := &testHandler{}
		_, _ = New().Parse(h, toks)
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
}

type testHandler struct {
	GlobalOutput, LocalOutput strings.Builder
}

func (h *testHandler) DefaultMacros() interfaces.MapMacros {
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
