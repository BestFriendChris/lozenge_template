package parser_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
	th "github.com/BestFriendChris/lozenge/handler/main_handler"
	. "github.com/BestFriendChris/lozenge/parser"
)

func enableDebug(t testing.TB) {
	DEBUGLOG = true
	t.Cleanup(func() { DEBUGLOG = false })
}

func TestParse_Simple(t *testing.T) {
	input := `
hi
there`[1:]

	output := ParseWithTestHandler(t, input)

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
				buf.WriteString("hi\n")
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
}

func TestParse_LozengeSimple(t *testing.T) {
	input := `
◊{foo := 1;baz_123 := 2}hi ◊foo bar
<span>◊baz_123</span>there
Loz-space is ignored "◊ "
Loz-newline is also ignored ◊
Loz-Loz is also ignored "◊◊"`[1:]

	output := ParseWithTestHandler(t, input)

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
			Loz-Loz is also ignored "◊"`)
	})
}

func TestParse_LozengeExpr(t *testing.T) {
	input := `foo ◊(1 + 2) bar`

	output := ParseWithTestHandler(t, input)

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
				buf.WriteString("foo ")
				buf.WriteString(fmt.Sprintf("%v", (1 + 2)))
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
}

func TestParse_LozengeBlock(t *testing.T) {
	input := `
◊{ foo := "Chris" }
Hello ◊foo`[1:]

	output := ParseWithTestHandler(t, input)

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
				foo := "Chris"
				buf.WriteString("\n")
				buf.WriteString("Hello ")
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
}

func TestParse_LozengeGlobal(t *testing.T) {
	input := `
◊^{import "strings"
	func myName() string {
		return "chris"
	}
}
◊{ foo := strings.ToUpper(myName()) }
Hello ◊foo`[1:]

	output := ParseWithTestHandler(t, input)

	t.Run("generate go", func(t *testing.T) {
		c := ic.New(t)
		c.Print(output)
		c.Expect(`
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
}

func TestParse_LozengeMacro_if(t *testing.T) {
	input := `
Try:
◊{val := "hi"}
◊.if val != "" {
	<span>◊val</span>
◊} else if 1 == 0 {
	<span>impossible</span>
◊} else {
	<span>default</span>
◊}
DONE`[1:]

	output := ParseWithTestHandler(t, input)

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
}

func TestParse_LozengeMacro_for(t *testing.T) {
	input := `
Try:
◊{vals := []string{"a", "b"}}
◊.for _, v := range vals {
	<span>◊v</span>
◊}
DONE`[1:]

	output := ParseWithTestHandler(t, input)

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
}

func TestParse_LozengeMacro_Custom(t *testing.T) {
	input := `
Try:
◊.CustomFoo(bar)
DONE`[1:]

	macros := map[string]Macro{
		"CustomFoo": func(p Parser, s string) (string, error) {
			before, after, _ := strings.Cut(s, ")")
			p.Handler().WriteContent(fmt.Sprintf("CUSTOM %s", before[1:]))
			return after, nil
		},
	}
	output := ParseWithTestHandlerWithMacros(t, input, macros)

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
				buf.WriteString("Try:\n")
				buf.WriteString("CUSTOM bar")
				buf.WriteString("\n")
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
			CUSTOM bar
			DONE`)
	})
}

func TestParse_ComplexExample(t *testing.T) {
	t.Skip()
	enableDebug(t)
	input := `
Try:
◊{ vals := []string{"a", "b", "c", "d"} }
◊.for _, v := range vals {
	◊.if v != "c" && v != "d" {
		◊.if v == "a" {
FOUND A: ◊v
		◊} else {
FOUND B: ◊v
		◊}
	◊} else if v == "c" {
FOUND C: ◊v
	◊} else {
FOUND D: ◊v
	◊}
◊}

◊.CustomFoo(bar) bar
DONE
`[1:]

	macros := map[string]Macro{
		"CustomFoo": func(p Parser, s string) (string, error) {
			before, after, _ := strings.Cut(s, ")")
			p.Handler().WriteContent(fmt.Sprintf("CUSTOM %s", before[1:]))
			return after, nil
		},
	}
	config := NewParserConfig().WithTrimSpaces(true)
	output := ParseWithTestHandlerWithMacrosWithConfig(t, input, macros, config)

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
				buf.WriteString("Try:")
				vals := []string{"a", "b", "c", "d"}
				buf.WriteString("\n")
				for _, v := range vals {
					buf.WriteString("\n")
					if v != "c" && v != "d" {
						buf.WriteString("\n")
						if v == "a" {
							buf.WriteString("\n")
							buf.WriteString("FOUND A:")
							buf.WriteString(fmt.Sprintf("%v", v))
							buf.WriteString("\n")
						} else {
							buf.WriteString("\n")
							buf.WriteString("FOUND B:")
							buf.WriteString(fmt.Sprintf("%v", v))
							buf.WriteString("\n")
						}
						buf.WriteString("\n")
					} else if v == "c" {
						buf.WriteString("\n")
						buf.WriteString("FOUND C:")
						buf.WriteString(fmt.Sprintf("%v", v))
						buf.WriteString("\n")
					} else {
						buf.WriteString("\n")
						buf.WriteString("FOUND D:")
						buf.WriteString(fmt.Sprintf("%v", v))
						buf.WriteString("\n")
					}
					buf.WriteString("\n")
				}
				buf.WriteString("\n")
				buf.WriteString("\n")
				buf.WriteString("CUSTOM bar")
				buf.WriteString("bar")
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
			
			
			
			FOUND A:a
			
			
			
			
			
			FOUND B:b
			
			
			
			
			FOUND C:c
			
			
			
			FOUND D:d
			
			
			
			CUSTOM barbarDONE`)
	})
}

func ParseWithTestHandler(t testing.TB, s string) string {
	t.Helper()
	return ParseWithTestHandlerWithMacros(t, s, nil)
}

func ParseWithTestHandlerWithMacros(t testing.TB, s string, macros map[string]Macro) string {
	t.Helper()
	defaultConfig := NewParserConfig()
	return ParseWithTestHandlerWithMacrosWithConfig(t, s, macros, defaultConfig)
}

func ParseWithTestHandlerWithMacrosWithConfig(t testing.TB, s string, macros map[string]Macro, config ParserConfig) string {
	t.Helper()
	testHandler := th.MainHandler{}
	p := New(&testHandler, macros, config)

	output, e := p.Parse(s)
	if e != nil {
		t.Fatalf("parse error: %q", e)
	}
	return output
}

func execAndReturnStdOut(t testing.TB, name, code string) string {
	var (
		tmpDir string
		err    error
	)
	tmpDir, err = os.MkdirTemp("", name)
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
