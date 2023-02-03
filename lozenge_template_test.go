package lozenge_template

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
	"github.com/BestFriendChris/lozenge_template/handler/main_handler"
	"github.com/BestFriendChris/lozenge_template/input"
	"github.com/BestFriendChris/lozenge_template/interfaces"
	"github.com/BestFriendChris/lozenge_template/internal/logic/token"
)

func TestLozengeTemplate_Generate(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		s := `
hi
there`[1:]
		output := GenerateWithTestHandler(t, s)

		t.Run("generate go", func(t *testing.T) {
			c := ic.New(t)
			c.Print(output)
			c.Expect(`
				package main
				
				import (
					"bytes"
					"fmt"
				)
				
				func main() {
					buf := new(bytes.Buffer)
				//line test.txt.◊:1
					buf.WriteString("hi\n")
				//line test.txt.◊:2
					buf.WriteString("there")
					fmt.Print(buf.String())
				}
				`)
		})
		t.Run("compile and run", func(t *testing.T) {
			if testing.Short() {
				t.Skip()
			}
			stdout := execAndReturnStdOut(t, "simple", output)
			c := ic.New(t)
			c.Print(stdout)
			c.Expect(`
				hi
				there`)
		})

	})
	t.Run("lozenge simple", func(t *testing.T) {
		s := `
◊{
foo := 1
baz_123 := 2}hi ◊foo bar
<span>◊baz_123</span>there
Loz-space is ignored "◊ "
Loz-newline is also ignored ◊
Loz-Loz is also ignored "◊◊"
Loz-EOL is also ignored ◊`[1:]
		output := GenerateWithTestHandler(t, s)

		t.Run("generate go", func(t *testing.T) {
			c := ic.New(t)
			c.Print(output)
			c.Expect(`
				package main
				
				import (
					"bytes"
					"fmt"
				)
				
				func main() {
					buf := new(bytes.Buffer)
				//line test.txt.◊:2
					foo := 1
				//line test.txt.◊:3
					baz_123 := 2
				//line test.txt.◊:3
					buf.WriteString("hi ")
				//line test.txt.◊:3
					buf.WriteString(fmt.Sprintf("%v", foo))
				//line test.txt.◊:3
					buf.WriteString(" bar\n")
				//line test.txt.◊:4
					buf.WriteString("<span>")
				//line test.txt.◊:4
					buf.WriteString(fmt.Sprintf("%v", baz_123))
				//line test.txt.◊:4
					buf.WriteString("</span>there\n")
				//line test.txt.◊:5
					buf.WriteString("Loz-space is ignored \"◊ \"\n")
				//line test.txt.◊:6
					buf.WriteString("Loz-newline is also ignored ◊\n")
				//line test.txt.◊:7
					buf.WriteString("Loz-Loz is also ignored \"◊")
				//line test.txt.◊:7
					buf.WriteString("\"\n")
				//line test.txt.◊:8
					buf.WriteString("Loz-EOL is also ignored ◊")
					fmt.Print(buf.String())
				}
				`)
		})
		t.Run("compile and run", func(t *testing.T) {
			if testing.Short() {
				t.Skip()
			}
			stdout := execAndReturnStdOut(t, "simple", output)
			c := ic.New(t)
			c.Print(stdout)
			c.Expect(`
				hi 1 bar
				<span>2</span>there
				Loz-space is ignored "◊ "
				Loz-newline is also ignored ◊
				Loz-Loz is also ignored "◊"
				Loz-EOL is also ignored ◊`)
		})

	})
	t.Run("lozenge simple - different marker", func(t *testing.T) {
		s := `
∆{
foo := 1
baz_123 := 2}hi ∆foo bar
<span>∆baz_123</span>there
Loz-space is ignored "∆ "
Loz-newline is also ignored ∆
Loz-Loz is also ignored "∆∆"
Loz-EOL is also ignored ∆`[1:]

		config := NewParserConfig().WithMarker('∆')
		output := GenerateWithTestHandlerWithMacrosWithConfig(t, s, nil, config)

		t.Run("generate go", func(t *testing.T) {
			c := ic.New(t)
			c.Print(output)
			c.Expect(`
				package main
				
				import (
					"bytes"
					"fmt"
				)
				
				func main() {
					buf := new(bytes.Buffer)
				//line test.txt.◊:2
					foo := 1
				//line test.txt.◊:3
					baz_123 := 2
				//line test.txt.◊:3
					buf.WriteString("hi ")
				//line test.txt.◊:3
					buf.WriteString(fmt.Sprintf("%v", foo))
				//line test.txt.◊:3
					buf.WriteString(" bar\n")
				//line test.txt.◊:4
					buf.WriteString("<span>")
				//line test.txt.◊:4
					buf.WriteString(fmt.Sprintf("%v", baz_123))
				//line test.txt.◊:4
					buf.WriteString("</span>there\n")
				//line test.txt.◊:5
					buf.WriteString("Loz-space is ignored \"∆ \"\n")
				//line test.txt.◊:6
					buf.WriteString("Loz-newline is also ignored ∆\n")
				//line test.txt.◊:7
					buf.WriteString("Loz-Loz is also ignored \"∆")
				//line test.txt.◊:7
					buf.WriteString("\"\n")
				//line test.txt.◊:8
					buf.WriteString("Loz-EOL is also ignored ∆")
					fmt.Print(buf.String())
				}
				`)
		})
		t.Run("compile and run", func(t *testing.T) {
			if testing.Short() {
				t.Skip()
			}
			stdout := execAndReturnStdOut(t, "simple", output)
			c := ic.New(t)
			c.Print(stdout)
			c.Expect(`
				hi 1 bar
				<span>2</span>there
				Loz-space is ignored "∆ "
				Loz-newline is also ignored ∆
				Loz-Loz is also ignored "∆"
				Loz-EOL is also ignored ∆`)
		})

	})
	t.Run("lozenge expression", func(t *testing.T) {
		s := `foo ◊(1 + 2) bar`

		output := GenerateWithTestHandler(t, s)

		t.Run("generate go", func(t *testing.T) {
			c := ic.New(t)
			c.Print(output)
			c.Expect(`
				package main
				
				import (
					"bytes"
					"fmt"
				)
				
				func main() {
					buf := new(bytes.Buffer)
				//line test.txt.◊:1
					buf.WriteString("foo ")
				//line test.txt.◊:1
					buf.WriteString(fmt.Sprintf("%v", (1 + 2)))
				//line test.txt.◊:1
					buf.WriteString(" bar")
					fmt.Print(buf.String())
				}
				`)
		})
		t.Run("compile and run", func(t *testing.T) {
			if testing.Short() {
				t.Skip()
			}
			stdout := execAndReturnStdOut(t, "simple", output)
			c := ic.New(t)
			c.Print(stdout)
			c.Expect(`foo 3 bar`)
		})

	})
	t.Run("lozenge local block", func(t *testing.T) {
		s := `
◊{ foo := "Chris" }
Hello ◊foo`[1:]

		output := GenerateWithTestHandler(t, s)

		t.Run("generate go", func(t *testing.T) {
			c := ic.New(t)
			c.Print(output)
			c.Expect(`
				package main
				
				import (
					"bytes"
					"fmt"
				)
				
				func main() {
					buf := new(bytes.Buffer)
				//line test.txt.◊:1
					foo := "Chris"
				//line test.txt.◊:1
					buf.WriteString("\n")
				//line test.txt.◊:2
					buf.WriteString("Hello ")
				//line test.txt.◊:2
					buf.WriteString(fmt.Sprintf("%v", foo))
					fmt.Print(buf.String())
				}
				`)
		})
		t.Run("compile and run", func(t *testing.T) {
			if testing.Short() {
				t.Skip()
			}
			stdout := execAndReturnStdOut(t, "simple", output)
			c := ic.New(t)
			c.Print(stdout)
			c.Expect(`Hello Chris`)
		})

	})
	t.Run("lozenge global block", func(t *testing.T) {
		s := `
◊^{import "strings"
	func myName() string {
		return "chris"
	}
}
◊{ foo := strings.ToUpper(myName()) }
Hello ◊foo`[1:]

		output := GenerateWithTestHandler(t, s)
		config := NewParserConfig().WithTrimSpaces()
		outputTrimSpaces := GenerateWithTestHandlerWithMacrosWithConfig(t, s, nil, config)

		t.Run("generate go", func(t *testing.T) {
			c := ic.New(t)
			c.Print(output)
			c.Expect(`
				package main
				
				import (
					"bytes"
					"fmt"
				//line test.txt.◊:1
					"strings"
				)
				
				//line test.txt.◊:2
				func myName() string {
				//line test.txt.◊:3
					return "chris"
				//line test.txt.◊:4
				}
				
				func main() {
					buf := new(bytes.Buffer)
				//line test.txt.◊:5
					buf.WriteString("\n")
				//line test.txt.◊:6
					foo := strings.ToUpper(myName())
				//line test.txt.◊:6
					buf.WriteString("\n")
				//line test.txt.◊:7
					buf.WriteString("Hello ")
				//line test.txt.◊:7
					buf.WriteString(fmt.Sprintf("%v", foo))
					fmt.Print(buf.String())
				}
				`)
		})
		t.Run("generate go - trim spaces", func(t *testing.T) {
			c := ic.New(t)
			c.Print(outputTrimSpaces)
			c.Expect(`
				package main
				
				import (
					"bytes"
					"fmt"
				//line test.txt.◊:1
					"strings"
				)
				
				//line test.txt.◊:2
				func myName() string {
				//line test.txt.◊:3
					return "chris"
				//line test.txt.◊:4
				}
				
				func main() {
					buf := new(bytes.Buffer)
				//line test.txt.◊:6
					foo := strings.ToUpper(myName())
				//line test.txt.◊:7
					buf.WriteString("Hello ")
				//line test.txt.◊:7
					buf.WriteString(fmt.Sprintf("%v", foo))
					fmt.Print(buf.String())
				}
				`)
		})
		t.Run("compile and run", func(t *testing.T) {
			if testing.Short() {
				t.Skip()
			}
			stdout := execAndReturnStdOut(t, "simple", output)
			c := ic.New(t)
			c.Print(stdout)
			c.Expect(`
				
				Hello CHRIS`)
		})
		t.Run("compile and run - trim spaces", func(t *testing.T) {
			if testing.Short() {
				t.Skip()
			}
			stdout := execAndReturnStdOut(t, "simple", outputTrimSpaces)
			c := ic.New(t)
			c.Print(stdout)
			c.Expect(`Hello CHRIS`)
		})

	})
	t.Run("lozenge macro - if", func(t *testing.T) {
		s := `
Try:
◊{val := "hi"}
◊.if val != "" {◊
	<span>◊val</span>
◊} else if 1 == 0 {◊
	<span>impossible</span>
◊} else {◊
	<span>default</span>
◊}
DONE`[1:]

		output := GenerateWithTestHandler(t, s)
		config := NewParserConfig().WithTrimSpaces()
		outputTrimSpaces := GenerateWithTestHandlerWithMacrosWithConfig(t, s, nil, config)

		t.Run("generate go", func(t *testing.T) {
			c := ic.New(t)
			c.Print(output)
			c.Expect(`
				package main
				
				import (
					"bytes"
					"fmt"
				)
				
				func main() {
					buf := new(bytes.Buffer)
				//line test.txt.◊:1
					buf.WriteString("Try:\n")
				//line test.txt.◊:2
					val := "hi"
				//line test.txt.◊:2
					buf.WriteString("\n")
				//line test.txt.◊:3
					if val != "" {
				//line test.txt.◊:3
						buf.WriteString("\n")
				//line test.txt.◊:4
						buf.WriteString("\t<span>")
				//line test.txt.◊:4
						buf.WriteString(fmt.Sprintf("%v", val))
				//line test.txt.◊:4
						buf.WriteString("</span>\n")
				//line test.txt.◊:5
					} else if 1 == 0 {
				//line test.txt.◊:5
						buf.WriteString("\n")
				//line test.txt.◊:6
						buf.WriteString("\t<span>impossible</span>\n")
				//line test.txt.◊:7
					} else {
				//line test.txt.◊:7
						buf.WriteString("\n")
				//line test.txt.◊:8
						buf.WriteString("\t<span>default</span>\n")
				//line test.txt.◊:9
					}
				//line test.txt.◊:9
					buf.WriteString("\n")
				//line test.txt.◊:10
					buf.WriteString("DONE")
					fmt.Print(buf.String())
				}
				`)
		})
		t.Run("generate go - trim spaces", func(t *testing.T) {
			c := ic.New(t)
			c.Print(outputTrimSpaces)
			c.Expect(`
				package main
				
				import (
					"bytes"
					"fmt"
				)
				
				func main() {
					buf := new(bytes.Buffer)
				//line test.txt.◊:1
					buf.WriteString("Try:\n")
				//line test.txt.◊:2
					val := "hi"
				//line test.txt.◊:3
					if val != "" {
				//line test.txt.◊:4
						buf.WriteString("\t<span>")
				//line test.txt.◊:4
						buf.WriteString(fmt.Sprintf("%v", val))
				//line test.txt.◊:4
						buf.WriteString("</span>\n")
				//line test.txt.◊:5
					} else if 1 == 0 {
				//line test.txt.◊:6
						buf.WriteString("\t<span>impossible</span>\n")
				//line test.txt.◊:7
					} else {
				//line test.txt.◊:8
						buf.WriteString("\t<span>default</span>\n")
				//line test.txt.◊:9
					}
				//line test.txt.◊:10
					buf.WriteString("DONE")
					fmt.Print(buf.String())
				}
				`)
		})
		t.Run("compile and run", func(t *testing.T) {
			if testing.Short() {
				t.Skip()
			}
			stdout := execAndReturnStdOut(t, "simple", output)
			c := ic.New(t)
			c.Print(stdout)
			c.Expect(`
				Try:
				
				
					<span>hi</span>
				
				DONE`)
		})
		t.Run("compile and run - trim spaces", func(t *testing.T) {
			if testing.Short() {
				t.Skip()
			}
			stdout := execAndReturnStdOut(t, "simple", outputTrimSpaces)
			c := ic.New(t)
			c.Print(stdout)
			c.Expect(`
				Try:
					<span>hi</span>
				DONE`)
		})

	})
	t.Run("lozenge macro - for", func(t *testing.T) {
		s := `
Try:
◊{vals := []string{"a", "b"}}
◊.for _, v := range vals {◊
	<span>◊v</span>
◊}
DONE`[1:]

		output := GenerateWithTestHandler(t, s)
		config := NewParserConfig().WithTrimSpaces()
		outputTrimSpaces := GenerateWithTestHandlerWithMacrosWithConfig(t, s, nil, config)

		t.Run("generate go", func(t *testing.T) {
			c := ic.New(t)
			c.Print(output)
			c.Expect(`
				package main
				
				import (
					"bytes"
					"fmt"
				)
				
				func main() {
					buf := new(bytes.Buffer)
				//line test.txt.◊:1
					buf.WriteString("Try:\n")
				//line test.txt.◊:2
					vals := []string{"a", "b"}
				//line test.txt.◊:2
					buf.WriteString("\n")
				//line test.txt.◊:3
					for _, v := range vals {
				//line test.txt.◊:3
						buf.WriteString("\n")
				//line test.txt.◊:4
						buf.WriteString("\t<span>")
				//line test.txt.◊:4
						buf.WriteString(fmt.Sprintf("%v", v))
				//line test.txt.◊:4
						buf.WriteString("</span>\n")
				//line test.txt.◊:5
					}
				//line test.txt.◊:5
					buf.WriteString("\n")
				//line test.txt.◊:6
					buf.WriteString("DONE")
					fmt.Print(buf.String())
				}
				`)
		})
		t.Run("generate go - trim spaces", func(t *testing.T) {
			c := ic.New(t)
			c.Print(outputTrimSpaces)
			c.Expect(`
				package main
				
				import (
					"bytes"
					"fmt"
				)
				
				func main() {
					buf := new(bytes.Buffer)
				//line test.txt.◊:1
					buf.WriteString("Try:\n")
				//line test.txt.◊:2
					vals := []string{"a", "b"}
				//line test.txt.◊:3
					for _, v := range vals {
				//line test.txt.◊:4
						buf.WriteString("\t<span>")
				//line test.txt.◊:4
						buf.WriteString(fmt.Sprintf("%v", v))
				//line test.txt.◊:4
						buf.WriteString("</span>\n")
				//line test.txt.◊:5
					}
				//line test.txt.◊:6
					buf.WriteString("DONE")
					fmt.Print(buf.String())
				}
				`)
		})
		t.Run("compile and run", func(t *testing.T) {
			if testing.Short() {
				t.Skip()
			}
			stdout := execAndReturnStdOut(t, "simple", output)
			c := ic.New(t)
			c.Print(stdout)
			c.Expect(`
				Try:
				
				
					<span>a</span>
				
					<span>b</span>
				
				DONE`)
		})
		t.Run("compile and run - trim spaces", func(t *testing.T) {
			if testing.Short() {
				t.Skip()
			}
			stdout := execAndReturnStdOut(t, "simple", outputTrimSpaces)
			c := ic.New(t)
			c.Print(stdout)
			c.Expect(`
				Try:
					<span>a</span>
					<span>b</span>
				DONE`)
		})

	})
	t.Run("complex example", func(t *testing.T) {
		s := `
Try:
◊{ vals := []string{"a", "b", "c", "d"} }
◊.for _, v := range vals {◊
	◊.if v != "c" && v != "d" {◊
		◊.if v == "a" {◊
FOUND A: ◊v
		◊} else {◊
FOUND B: ◊v
		◊}
	◊} else if v == "c" {◊
FOUND C: ◊v
	◊} else {◊
FOUND D: ◊v
	◊}
◊}

◊.LogValue(1 + 2) bar
DONE
`[1:]

		macros := interfaces.NewMacros()
		macros.Add(LogValue{})
		config := NewParserConfig()

		output := GenerateWithTestHandlerWithMacrosWithConfig(t, s, macros, config)
		outputTrimSpaces := GenerateWithTestHandlerWithMacrosWithConfig(t, s, macros, config.WithTrimSpaces())

		t.Run("generate go", func(t *testing.T) {
			c := ic.New(t)
			c.Print(output)
			c.Expect(`
				package main
				
				import (
					"bytes"
					"fmt"
				)
				
				func main() {
					buf := new(bytes.Buffer)
				//line test.txt.◊:1
					buf.WriteString("Try:\n")
				//line test.txt.◊:2
					vals := []string{"a", "b", "c", "d"}
				//line test.txt.◊:2
					buf.WriteString("\n")
				//line test.txt.◊:3
					for _, v := range vals {
				//line test.txt.◊:3
						buf.WriteString("\n")
				//line test.txt.◊:4
						buf.WriteString("\t")
				//line test.txt.◊:4
						if v != "c" && v != "d" {
				//line test.txt.◊:4
							buf.WriteString("\n")
				//line test.txt.◊:5
							buf.WriteString("\t\t")
				//line test.txt.◊:5
							if v == "a" {
				//line test.txt.◊:5
								buf.WriteString("\n")
				//line test.txt.◊:6
								buf.WriteString("FOUND A: ")
				//line test.txt.◊:6
								buf.WriteString(fmt.Sprintf("%v", v))
				//line test.txt.◊:6
								buf.WriteString("\n")
				//line test.txt.◊:7
								buf.WriteString("\t\t")
				//line test.txt.◊:7
							} else {
				//line test.txt.◊:7
								buf.WriteString("\n")
				//line test.txt.◊:8
								buf.WriteString("FOUND B: ")
				//line test.txt.◊:8
								buf.WriteString(fmt.Sprintf("%v", v))
				//line test.txt.◊:8
								buf.WriteString("\n")
				//line test.txt.◊:9
								buf.WriteString("\t\t")
				//line test.txt.◊:9
							}
				//line test.txt.◊:9
							buf.WriteString("\n")
				//line test.txt.◊:10
							buf.WriteString("\t")
				//line test.txt.◊:10
						} else if v == "c" {
				//line test.txt.◊:10
							buf.WriteString("\n")
				//line test.txt.◊:11
							buf.WriteString("FOUND C: ")
				//line test.txt.◊:11
							buf.WriteString(fmt.Sprintf("%v", v))
				//line test.txt.◊:11
							buf.WriteString("\n")
				//line test.txt.◊:12
							buf.WriteString("\t")
				//line test.txt.◊:12
						} else {
				//line test.txt.◊:12
							buf.WriteString("\n")
				//line test.txt.◊:13
							buf.WriteString("FOUND D: ")
				//line test.txt.◊:13
							buf.WriteString(fmt.Sprintf("%v", v))
				//line test.txt.◊:13
							buf.WriteString("\n")
				//line test.txt.◊:14
							buf.WriteString("\t")
				//line test.txt.◊:14
						}
				//line test.txt.◊:14
						buf.WriteString("\n")
				//line test.txt.◊:15
					}
				//line test.txt.◊:15
					buf.WriteString("\n")
				//line test.txt.◊:16
					buf.WriteString("\n")
				//line test:17
					buf.WriteString("(1 + 2) = ")
				//line test.txt.◊:17
					buf.WriteString(fmt.Sprintf("%v", (1 + 2)))
				//line test.txt.◊:17
					buf.WriteString(" bar\n")
				//line test.txt.◊:18
					buf.WriteString("DONE\n")
					fmt.Print(buf.String())
				}
				`)
		})
		t.Run("generate go - trim spaces", func(t *testing.T) {
			c := ic.New(t)
			c.Print(outputTrimSpaces)
			c.Expect(`
				package main
				
				import (
					"bytes"
					"fmt"
				)
				
				func main() {
					buf := new(bytes.Buffer)
				//line test.txt.◊:1
					buf.WriteString("Try:\n")
				//line test.txt.◊:2
					vals := []string{"a", "b", "c", "d"}
				//line test.txt.◊:3
					for _, v := range vals {
				//line test.txt.◊:4
						if v != "c" && v != "d" {
				//line test.txt.◊:5
							if v == "a" {
				//line test.txt.◊:6
								buf.WriteString("FOUND A: ")
				//line test.txt.◊:6
								buf.WriteString(fmt.Sprintf("%v", v))
				//line test.txt.◊:6
								buf.WriteString("\n")
				//line test.txt.◊:7
							} else {
				//line test.txt.◊:8
								buf.WriteString("FOUND B: ")
				//line test.txt.◊:8
								buf.WriteString(fmt.Sprintf("%v", v))
				//line test.txt.◊:8
								buf.WriteString("\n")
				//line test.txt.◊:9
							}
				//line test.txt.◊:10
						} else if v == "c" {
				//line test.txt.◊:11
							buf.WriteString("FOUND C: ")
				//line test.txt.◊:11
							buf.WriteString(fmt.Sprintf("%v", v))
				//line test.txt.◊:11
							buf.WriteString("\n")
				//line test.txt.◊:12
						} else {
				//line test.txt.◊:13
							buf.WriteString("FOUND D: ")
				//line test.txt.◊:13
							buf.WriteString(fmt.Sprintf("%v", v))
				//line test.txt.◊:13
							buf.WriteString("\n")
				//line test.txt.◊:14
						}
				//line test.txt.◊:15
					}
				//line test.txt.◊:16
					buf.WriteString("\n")
				//line test:17
					buf.WriteString("(1 + 2) = ")
				//line test.txt.◊:17
					buf.WriteString(fmt.Sprintf("%v", (1 + 2)))
				//line test.txt.◊:17
					buf.WriteString(" bar\n")
				//line test.txt.◊:18
					buf.WriteString("DONE\n")
					fmt.Print(buf.String())
				}
				`)
		})
		t.Run("compile and run", func(t *testing.T) {
			if testing.Short() {
				t.Skip()
			}
			stdout := execAndReturnStdOut(t, "simple", output)
			c := ic.New(t)
			c.Print(stdout)
			c.Expect(`
				Try:
				
				
					
						
				FOUND A: a
						
					
				
					
						
				FOUND B: b
						
					
				
					
				FOUND C: c
					
				
					
				FOUND D: d
					
				
				
				(1 + 2) = 3 bar
				DONE
				`)
		})
		t.Run("compile and run - trim spaces", func(t *testing.T) {
			if testing.Short() {
				t.Skip()
			}
			stdout := execAndReturnStdOut(t, "simple", outputTrimSpaces)
			c := ic.New(t)
			c.Print(stdout)
			c.Expect(`
				Try:
				FOUND A: a
				FOUND B: b
				FOUND C: c
				FOUND D: d
				
				(1 + 2) = 3 bar
				DONE
				`)
		})

	})
}

func TestLozengeTemplate_Generate_errorCases(t *testing.T) {
	t.Run("show context around error", func(t *testing.T) {
		t.Run("single line", func(t *testing.T) {
			s := `foo ◊(1 + 2 bar`

			testHandler := &main_handler.MainHandler{}
			p := New(nil, NewParserConfig())

			in := input.NewInput("test.txt.◊", s)
			_, err := p.Generate(testHandler, in)

			c := ic.New(t)
			c.PrintSection("error")
			c.Println(err)

			c.Expect(`
				################################################################################
				# error
				################################################################################
				line 1: foo ◊(1 + 2 bar
				             ▲
				             └── did not find matched ')'
				`)
		})
	})
}

func GenerateWithTestHandler(t testing.TB, s string) string {
	t.Helper()
	return GenerateWithTestHandlerWithMacros(t, s, nil)
}

func GenerateWithTestHandlerWithMacros(t testing.TB, s string, overrideMacros *interfaces.Macros) string {
	t.Helper()
	return GenerateWithTestHandlerWithMacrosWithConfig(t, s, overrideMacros, NewParserConfig())
}

func GenerateWithTestHandlerWithMacrosWithConfig(t testing.TB, s string, overrideMacros *interfaces.Macros, config ParserConfig) string {
	t.Helper()
	testHandler := &main_handler.MainHandler{}
	p := New(overrideMacros, config)

	in := input.NewInput("test.txt.◊", s)
	output, e := p.Generate(testHandler, in)
	if e != nil {
		t.Fatalf("generate error: %q", e)
	}
	return output
}
func execAndReturnStdOut(t testing.TB, name, code string) string {
	tmpDir, err := os.MkdirTemp("", name)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(tmpDir)
	})

	mainFname := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(mainFname, []byte(code), 0600)
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = tmpDir
	var sysout, syserr bytes.Buffer
	cmd.Stdout = &sysout
	cmd.Stderr = &syserr
	err = cmd.Run()
	if syserr.String() != "" || err != nil {
		t.Fatalf(`
err: %q
--------------------
syserr:
--------------------
%s
--------------------
stdout:
--------------------
%s`[1:], err, syserr.String(), sysout.String())
	}

	return sysout.String()
}

// ◊.LogValue(1 + 2) => "(1 + 2) = 3"
type LogValue struct{}

func (m LogValue) Name() string {
	return "LogValue"
}

func (m LogValue) NextTokens(ct interfaces.ContentTokenizer, in *input.Input) (toks []*token.Token, err error) {
	_, _ = in.ConsumeString(m.Name())
	var valToks []*token.Token
	valToks, err = ct.ParseGoCodeFromTo(in, token.TTcodeLocalExpr, '(', ')', true)
	if err != nil {
		return nil, err
	}
	if len(valToks) != 1 {
		return nil, fmt.Errorf("got len(%d) want len(1) of valtoks\n%v", len(valToks), valToks)
	}
	valTok := valToks[0]
	contentSlc := input.NewSlice("test", fmt.Sprintf("%s = ", valTok.Slc.S), valTok.Slc.Start, valTok.Slc.End)
	contentToken := token.NewToken(token.TTcontent, contentSlc)
	return []*token.Token{contentToken, valTok}, nil
}

func (m LogValue) Parse(_ interfaces.TemplateHandler, toks []*token.Token) (rest []*token.Token, err error) {
	return toks, nil
}
