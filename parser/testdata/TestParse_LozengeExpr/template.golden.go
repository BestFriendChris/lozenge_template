package main

import (
	"bytes"
	"fmt"
)

func main() {
	buf := new(bytes.Buffer)
	buf.WriteString("foo ")
	buf.WriteString(fmt.Sprintf("%v", (1 + 2)))
	buf.WriteString(" bar")
	fmt.Print(buf.String())
}
