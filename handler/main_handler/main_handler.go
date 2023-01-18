package main_handler

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/BestFriendChris/lozenge/parser"
)

type MainHandler struct {
	Content      []string
	GlobalCode   []string
	InlineOutput []string
}

func (th *MainHandler) DefaultMacros() map[string]parser.Macro {
	return nil
}

func (th *MainHandler) WriteContent(s string) {
	th.Content = append(th.Content, s)
	th.InlineOutput = append(th.InlineOutput, fmt.Sprintf("buf.WriteString(%q)", s))
}

func (th *MainHandler) WriteCodeExpression(s string) {
	th.InlineOutput = append(th.InlineOutput, fmt.Sprintf("buf.WriteString(fmt.Sprintf(%q, %s))", "%v", s))
}

func (th *MainHandler) WriteCodeBlock(s string) {
	th.InlineOutput = append(th.InlineOutput, s)
}

func (th *MainHandler) WriteCodeGlobalBlock(s string) {
	th.GlobalCode = append(th.GlobalCode, s)
}

var STATIC = []string{
	`package main
import "bytes"
import "fmt"
`,

	`
func main() {
	buf := new(bytes.Buffer)
`,

	`
fmt.Print(buf.String())
}`,
}

func (th *MainHandler) Done() (string, error) {
	var fullOutput bytes.Buffer

	fullOutput.WriteString(STATIC[0])
	fullOutput.WriteString(strings.Join(th.GlobalCode, "\n"))
	fullOutput.WriteString(STATIC[1])
	fullOutput.WriteString(strings.Join(th.InlineOutput, "\n"))
	fullOutput.WriteString(STATIC[2])

	return fullOutput.String(), nil
}
