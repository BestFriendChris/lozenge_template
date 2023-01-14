package main

import (
	"bytes"
	"fmt"
)

func main() {
	buf := new(bytes.Buffer)
	buf.WriteString("hi\n")
	buf.WriteString("there")
	fmt.Print(buf.String())
}
