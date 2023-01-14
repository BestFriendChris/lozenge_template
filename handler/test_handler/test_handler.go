package test_handler

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/BestFriendChris/lozenge/parser"
	"mvdan.cc/gofumpt/format"
)

type TestHandler struct {
	Content      []string
	GlobalCode   []string
	InlineOutput []string
}

func (th *TestHandler) DefaultMacros() map[string]parser.Macro {
	return nil
}

func (th *TestHandler) WriteContent(s string) {
	th.Content = append(th.Content, s)
	th.InlineOutput = append(th.InlineOutput, fmt.Sprintf("buf.WriteString(%q)", s))
}

func (th *TestHandler) WriteCodeInline(s string) {
	th.InlineOutput = append(th.InlineOutput, fmt.Sprintf("buf.WriteString(fmt.Sprintf(%q, %s))", "%v", s))
}

func (th *TestHandler) WriteCodeBlock(s string) {
	th.InlineOutput = append(th.InlineOutput, s)
}

func (th *TestHandler) WriteCodeGlobalBlock(s string) {
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

func (th *TestHandler) Done(opts format.Options) (string, error) {
	var fullOutput bytes.Buffer

	fullOutput.WriteString(STATIC[0])
	fullOutput.WriteString(strings.Join(th.GlobalCode, "\n"))
	fullOutput.WriteString(STATIC[1])
	fullOutput.WriteString(strings.Join(th.InlineOutput, "\n"))
	fullOutput.WriteString(STATIC[2])

	formatted, err := format.Source([]byte(fullOutput.String()), opts)
	if err != nil {
		fmt.Printf("unable to format:\n%s\n", fullOutput.String())
		return "", err
	}
	return string(formatted), nil
}
