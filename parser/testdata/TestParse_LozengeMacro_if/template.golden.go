package main

import (
	"bytes"
	"fmt"
)

func main() {
	buf := new(bytes.Buffer)
	buf.WriteString("Try:\n")
	val := "hi"
	buf.WriteString("\n")
	if val != "" {
		buf.WriteString("\n")
		buf.WriteString("\t<span>")
		buf.WriteString(fmt.Sprintf("%v", val))
		buf.WriteString("</span>\n")
	} else if 1 == 0 {
		buf.WriteString("\n")
		buf.WriteString("\t<span>impossible</span>\n")
	} else {
		buf.WriteString("\n")
		buf.WriteString("\t<span>default</span>\n")
	}
	buf.WriteString("\n")
	buf.WriteString("DONE")
	fmt.Print(buf.String())
}
