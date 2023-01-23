package main_handler

import (
	"fmt"
	"strings"

	"github.com/BestFriendChris/lozenge/interfaces"
)

type MainHandler struct {
	Content      []string
	GlobalCode   []string
	InlineOutput []string
}

func (th *MainHandler) DefaultMacros() *interfaces.Macros {
	return nil
}

func (th *MainHandler) WriteTextContent(s string) {
	th.Content = append(th.Content, s)
	th.InlineOutput = append(th.InlineOutput, fmt.Sprintf("buf.WriteString(%q)", s))
}

func (th *MainHandler) WriteCodeLocalExpression(s string) {
	th.InlineOutput = append(th.InlineOutput, fmt.Sprintf("buf.WriteString(fmt.Sprintf(%q, %s))", "%v", s))
}

func (th *MainHandler) WriteCodeLocalBlock(s string) {
	th.InlineOutput = append(th.InlineOutput, s)
}

func (th *MainHandler) WriteCodeGlobalBlock(s string) {
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
