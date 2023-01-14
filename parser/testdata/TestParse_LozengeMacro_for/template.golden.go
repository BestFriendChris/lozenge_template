package main

import (
	"bytes"
	"fmt"
)

func main() {
	buf := new(bytes.Buffer)
	buf.WriteString("Try:\n")
	vals := []string{"a", "b"}
	buf.WriteString("\n")
	for _, v := range vals {
		buf.WriteString("\n")
		buf.WriteString("\t<span>")
		buf.WriteString(fmt.Sprintf("%v", v))
		buf.WriteString("</span>\n")
	}
	buf.WriteString("\n")
	buf.WriteString("DONE")
	fmt.Print(buf.String())
}
