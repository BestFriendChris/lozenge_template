package main

import (
	"bytes"
	"fmt"
	"strings"
)

func myName() string {
	return "chris"
}

func main() {
	buf := new(bytes.Buffer)
	buf.WriteString("\n")
	foo := strings.ToUpper(myName())
	buf.WriteString("\n")
	buf.WriteString("Hello ")
	buf.WriteString(fmt.Sprintf("%v", foo))
	fmt.Print(buf.String())
}
