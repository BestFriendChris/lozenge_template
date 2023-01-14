package main

import (
	"bytes"
	"fmt"
)

func main() {
	buf := new(bytes.Buffer)
	foo := "Chris"
	buf.WriteString("\n")
	buf.WriteString("Hello ")
	buf.WriteString(fmt.Sprintf("%v", foo))
	fmt.Print(buf.String())
}
