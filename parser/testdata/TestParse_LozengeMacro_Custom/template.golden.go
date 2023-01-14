package main

import (
	"bytes"
	"fmt"
)

func main() {
	buf := new(bytes.Buffer)
	buf.WriteString("Try:\n")
	buf.WriteString("CUSTOM bar")
	buf.WriteString("\n")
	buf.WriteString("DONE")
	fmt.Print(buf.String())
}
