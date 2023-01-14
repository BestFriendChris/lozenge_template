package main

import (
	"bytes"
	"fmt"
)

func main() {
	buf := new(bytes.Buffer)
	foo := 1
	baz_123 := 2
	buf.WriteString("hi ")
	buf.WriteString(fmt.Sprintf("%v", foo))
	buf.WriteString(" bar\n")
	buf.WriteString("<span>")
	buf.WriteString(fmt.Sprintf("%v", baz_123))
	buf.WriteString("</span>there\n")
	buf.WriteString("Loz-space is ignored \"")
	buf.WriteString("◊ ")
	buf.WriteString("\"\n")
	buf.WriteString("Loz-newline is also ignored ")
	buf.WriteString("◊\n")
	buf.WriteString("Loz-Loz is also ignored \"")
	buf.WriteString("◊")
	buf.WriteString("\"")
	fmt.Print(buf.String())
}
