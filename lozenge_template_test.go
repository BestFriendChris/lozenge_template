package lozenge_template

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
	"github.com/BestFriendChris/lozenge_template/handler/main_handler"
	"github.com/BestFriendChris/lozenge_template/interfaces"
	"github.com/BestFriendChris/lozenge_template/internal/logic/token"
)

func TestGenerate(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		input := `
hi
there`[1:]
		output := GenerateWithTestHandler(t, input)

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
					buf.WriteString("hi\nthere")
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
		input := `
◊{
foo := 1
baz_123 := 2}hi ◊foo bar
<span>◊baz_123</span>there
Loz-space is ignored "◊ "
Loz-newline is also ignored ◊
Loz-Loz is also ignored "◊◊"
Loz-EOL is also ignored ◊`[1:]
		output := GenerateWithTestHandler(t, input)

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
					buf.WriteString(" bar\n<span>")
					buf.WriteString(fmt.Sprintf("%v", baz_123))
					buf.WriteString("</span>there\nLoz-space is ignored \"◊ \"\nLoz-newline is also")
					buf.WriteString(" ignored ◊\nLoz-Loz is also ignored \"◊\"\nLoz-EOL is also ")
					buf.WriteString("ignored ◊")
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
		input := `
∆{
foo := 1
baz_123 := 2}hi ∆foo bar
<span>∆baz_123</span>there
Loz-space is ignored "∆ "
Loz-newline is also ignored ∆
Loz-Loz is also ignored "∆∆"
Loz-EOL is also ignored ∆`[1:]

		config := NewParserConfig().WithMarker('∆')
		output := GenerateWithTestHandlerWithMacrosWithConfig(t, input, nil, config)

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
					buf.WriteString(" bar\n<span>")
					buf.WriteString(fmt.Sprintf("%v", baz_123))
					buf.WriteString("</span>there\nLoz-space is ignored \"∆ \"\nLoz-newline is also")
					buf.WriteString(" ignored ∆\nLoz-Loz is also ignored \"∆\"\nLoz-EOL is also ")
					buf.WriteString("ignored ∆")
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
		input := `foo ◊(1 + 2) bar`

		output := GenerateWithTestHandler(t, input)

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

	})
	t.Run("lozenge local block", func(t *testing.T) {
		input := `
◊{ foo := "Chris" }
Hello ◊foo`[1:]

		output := GenerateWithTestHandler(t, input)

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
					buf.WriteString("\nHello ")
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
		input := `
◊^{import "strings"
	func myName() string {
		return "chris"
	}
}
◊{ foo := strings.ToUpper(myName()) }
Hello ◊foo`[1:]

		output := GenerateWithTestHandler(t, input)
		config := NewParserConfig().WithTrimSpaces()
		outputTrimSpaces := GenerateWithTestHandlerWithMacrosWithConfig(t, input, nil, config)

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
					buf.WriteString("\nHello ")
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
					"strings"
				)
				
				func myName() string {
					return "chris"
				}
				
				func main() {
					buf := new(bytes.Buffer)
					foo := strings.ToUpper(myName())
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
		input := `
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

		output := GenerateWithTestHandler(t, input)
		config := NewParserConfig().WithTrimSpaces()
		outputTrimSpaces := GenerateWithTestHandlerWithMacrosWithConfig(t, input, nil, config)

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
						buf.WriteString("\n\t<span>")
						buf.WriteString(fmt.Sprintf("%v", val))
						buf.WriteString("</span>\n")
					} else if 1 == 0 {
						buf.WriteString("\n\t<span>impossible</span>\n")
					} else {
						buf.WriteString("\n\t<span>default</span>\n")
					}
					buf.WriteString("\nDONE")
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
					buf.WriteString("Try:\n")
					val := "hi"
					if val != "" {
						buf.WriteString("\t<span>")
						buf.WriteString(fmt.Sprintf("%v", val))
						buf.WriteString("</span>\n")
					} else if 1 == 0 {
						buf.WriteString("\t<span>impossible</span>\n")
					} else {
						buf.WriteString("\t<span>default</span>\n")
					}
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
		input := `
Try:
◊{vals := []string{"a", "b"}}
◊.for _, v := range vals {◊
	<span>◊v</span>
◊}
DONE`[1:]

		output := GenerateWithTestHandler(t, input)
		config := NewParserConfig().WithTrimSpaces()
		outputTrimSpaces := GenerateWithTestHandlerWithMacrosWithConfig(t, input, nil, config)

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
						buf.WriteString("\n\t<span>")
						buf.WriteString(fmt.Sprintf("%v", v))
						buf.WriteString("</span>\n")
					}
					buf.WriteString("\nDONE")
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
					buf.WriteString("Try:\n")
					vals := []string{"a", "b"}
					for _, v := range vals {
						buf.WriteString("\t<span>")
						buf.WriteString(fmt.Sprintf("%v", v))
						buf.WriteString("</span>\n")
					}
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
		input := `
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

		output := GenerateWithTestHandlerWithMacrosWithConfig(t, input, macros, config)
		outputTrimSpaces := GenerateWithTestHandlerWithMacrosWithConfig(t, input, macros, config.WithTrimSpaces())

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
					vals := []string{"a", "b", "c", "d"}
					buf.WriteString("\n")
					for _, v := range vals {
						buf.WriteString("\n\t")
						if v != "c" && v != "d" {
							buf.WriteString("\n\t\t")
							if v == "a" {
								buf.WriteString("\nFOUND A: ")
								buf.WriteString(fmt.Sprintf("%v", v))
								buf.WriteString("\n\t\t")
							} else {
								buf.WriteString("\nFOUND B: ")
								buf.WriteString(fmt.Sprintf("%v", v))
								buf.WriteString("\n\t\t")
							}
							buf.WriteString("\n\t")
						} else if v == "c" {
							buf.WriteString("\nFOUND C: ")
							buf.WriteString(fmt.Sprintf("%v", v))
							buf.WriteString("\n\t")
						} else {
							buf.WriteString("\nFOUND D: ")
							buf.WriteString(fmt.Sprintf("%v", v))
							buf.WriteString("\n\t")
						}
						buf.WriteString("\n")
					}
					buf.WriteString("\n\n")
					buf.WriteString("(1 + 2) = ")
					buf.WriteString(fmt.Sprintf("%v", (1 + 2)))
					buf.WriteString(" bar\nDONE\n")
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
					buf.WriteString("Try:\n")
					vals := []string{"a", "b", "c", "d"}
					for _, v := range vals {
						if v != "c" && v != "d" {
							if v == "a" {
								buf.WriteString("FOUND A: ")
								buf.WriteString(fmt.Sprintf("%v", v))
								buf.WriteString("\n")
							} else {
								buf.WriteString("FOUND B: ")
								buf.WriteString(fmt.Sprintf("%v", v))
								buf.WriteString("\n")
							}
						} else if v == "c" {
							buf.WriteString("FOUND C: ")
							buf.WriteString(fmt.Sprintf("%v", v))
							buf.WriteString("\n")
						} else {
							buf.WriteString("FOUND D: ")
							buf.WriteString(fmt.Sprintf("%v", v))
							buf.WriteString("\n")
						}
					}
					buf.WriteString("\n")
					buf.WriteString("(1 + 2) = ")
					buf.WriteString(fmt.Sprintf("%v", (1 + 2)))
					buf.WriteString(" bar\nDONE\n")
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

	output, e := p.Generate(testHandler, s)
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

func (m LogValue) NextTokens(ct interfaces.ContentTokenizer, rest string) ([]*token.Token, string, error) {
	rest = strings.TrimPrefix(rest, m.Name())
	fmt.Printf("Logvalue with %q\n", rest)

	runes := []rune(rest)
	var tok *token.Token
	tok, rest, _ = ct.ParseGoCodeFromTo(runes, token.TTcodeLocalExpr, '(', ')', true)
	contentToken := token.NewToken(token.TTcontent, fmt.Sprintf("%s = ", tok.S))
	return []*token.Token{contentToken, tok}, rest, nil
}

func (m LogValue) Parse(_ interfaces.TemplateHandler, toks []*token.Token) (rest []*token.Token, err error) {
	return toks, nil
}
