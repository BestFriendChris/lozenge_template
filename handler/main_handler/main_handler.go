package main_handler

import (
	"fmt"
	"strings"

	"github.com/BestFriendChris/lozenge_template/input"
	"github.com/BestFriendChris/lozenge_template/interfaces"
)

type MainHandler struct {
	Content      []string
	GlobalCode   []string
	InlineOutput []string
}

func (th *MainHandler) DefaultMacros() *interfaces.Macros {
	return nil
}

func (th *MainHandler) WriteTextContent(slc input.Slice) {
	th.Content = append(th.Content, slc.String())
	s := fmt.Sprintf(
		"%sbuf.WriteString(%q)",
		locationComment(slc),
		slc.S,
	)
	th.InlineOutput = append(th.InlineOutput, s)
}

func (th *MainHandler) WriteCodeLocalExpression(slc input.Slice) {
	s := fmt.Sprintf(
		"%sbuf.WriteString(fmt.Sprintf(%q, %s))",
		locationComment(slc),
		"%v",
		slc.S,
	)
	th.InlineOutput = append(th.InlineOutput, s)
}

func (th *MainHandler) WriteCodeLocalBlock(slc input.Slice) {
	s := fmt.Sprintf("%s%s", locationComment(slc), slc.S)
	th.InlineOutput = append(th.InlineOutput, s)
}

func (th *MainHandler) WriteCodeGlobalBlock(slc input.Slice) {
	s := fmt.Sprintf("%s%s", locationComment(slc), slc.S)
	th.GlobalCode = append(th.GlobalCode, s)
}

var format = `
package main
import "bytes"
import "fmt"
%s
func main() {
	buf := new(bytes.Buffer)
%s
	fmt.Print(buf.String())
}
`[1:]

func (th *MainHandler) Done() (string, error) {
	return fmt.Sprintf(
		format,
		strings.Join(th.GlobalCode, "\n"),
		strings.Join(th.InlineOutput, "\n"),
	), nil
}

func locationComment(slc input.Slice) string {
	return fmt.Sprintf("//line %s:%d\n", slc.Name, slc.Start.Row)
}
