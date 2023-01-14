package main

import (
	"bytes"
	"fmt"
)

func main() {
	buf := new(bytes.Buffer)
	buf.WriteString("Try:\n")
	vals := []string{"a", "b", "c", "d"}
	buf.WriteString("\n")
	for _, v := range vals {
		buf.WriteString("\n\t")
		if v != "c" && v != "d" {
			buf.WriteString("\n\t\t")
			if v == "a" {
				buf.WriteString("\nFOUND A:")
				buf.WriteString(fmt.Sprintf("%v", v))
			} else {
				buf.WriteString("\nFOUND B:")
				buf.WriteString(fmt.Sprintf("%v", v))
			}
		} else if v == "c" {
			buf.WriteString("\nFOUND C:")
			buf.WriteString(fmt.Sprintf("%v", v))
		} else {
			buf.WriteString("\nFOUND D:")
			buf.WriteString(fmt.Sprintf("%v", v))
		}
	}
	buf.WriteString("\n\n")
	buf.WriteString("CUSTOM bar")
	buf.WriteString("\nDONE\n")
	fmt.Print(buf.String())
}
