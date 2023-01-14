package parser_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	th "github.com/BestFriendChris/lozenge/handler/test_handler"
	. "github.com/BestFriendChris/lozenge/parser"
	"github.com/sebdah/goldie/v2"
)

func enableDebug(t testing.TB) {
	DEBUGLOG = true
	t.Cleanup(func() { DEBUGLOG = false })
}

func TestParse_Simple(t *testing.T) {
	input := `hi
there`

	output := ParseWithTestHandler(t, input)

	testGoCodeGeneration(t, output)
	testGoCodeExecution(t, output)
}

func TestParse_LozengeSimple(t *testing.T) {
	input := `◊{foo := 1;baz_123 := 2}hi ◊foo bar
<span>◊baz_123</span>there
Loz-space is ignored "◊ "
Loz-newline is also ignored ◊
Loz-Loz is also ignored "◊◊"`

	output := ParseWithTestHandler(t, input)

	testGoCodeGeneration(t, output)
	testGoCodeExecution(t, output)
}

func TestParse_LozengeExpr(t *testing.T) {
	input := `foo ◊(1 + 2) bar`

	output := ParseWithTestHandler(t, input)

	testGoCodeGeneration(t, output)
	testGoCodeExecution(t, output)
}

func TestParse_LozengeBlock(t *testing.T) {
	input := `◊{ foo := "Chris" }
Hello ◊foo`

	output := ParseWithTestHandler(t, input)

	testGoCodeGeneration(t, output)
	testGoCodeExecution(t, output)
}

func TestParse_LozengeGlobal(t *testing.T) {
	input := `◊^{import "strings"
	func myName() string {
		return "chris"
	}
}
◊{ foo := strings.ToUpper(myName()) }
Hello ◊foo`

	output := ParseWithTestHandler(t, input)

	testGoCodeGeneration(t, output)
	testGoCodeExecution(t, output)
}

func TestParse_LozengeMacro_if(t *testing.T) {
	input := `Try:
◊{val := "hi"}
◊.if val != "" {
	<span>◊val</span>
◊} else if 1 == 0 {
	<span>impossible</span>
◊} else {
	<span>default</span>
◊}
DONE`

	output := ParseWithTestHandler(t, input)

	testGoCodeGeneration(t, output)
	testGoCodeExecution(t, output)
}

func TestParse_LozengeMacro_for(t *testing.T) {
	input := `Try:
◊{vals := []string{"a", "b"}}
◊.for _, v := range vals {
	<span>◊v</span>
◊}
DONE`

	output := ParseWithTestHandler(t, input)

	testGoCodeGeneration(t, output)
	testGoCodeExecution(t, output)
}

func TestParse_LozengeMacro_Custom(t *testing.T) {
	input := `Try:
◊.CustomFoo(bar)
DONE`

	macros := map[string]Macro{
		"CustomFoo": func(p Parser, s string) (string, error) {
			before, after, _ := strings.Cut(s, ")")
			p.H.WriteContent(fmt.Sprintf("CUSTOM %s", before[1:]))
			return after, nil
		},
	}
	output := ParseWithTestHandlerWithMacros(t, input, macros)

	testGoCodeGeneration(t, output)
	testGoCodeExecution(t, output)
}

func TestParse_ComplexExample(t *testing.T) {
	t.Skip()
	enableDebug(t)
	input := `Try:
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
`

	macros := map[string]Macro{
		"CustomFoo": func(p Parser, s string) (string, error) {
			before, after, _ := strings.Cut(s, ")")
			p.H.WriteContent(fmt.Sprintf("CUSTOM %s", before[1:]))
			return after, nil
		},
	}
	config := NewParserConfig().WithTrimSpaces(true)
	output := ParseWithTestHandlerWithMacrosWithConfig(t, input, macros, config)

	testGoCodeGeneration(t, output)
	testGoCodeExecution(t, output)
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
	testHandler := th.TestHandler{}
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
		t.Fatalf(`err: %q
--------------------
syserr:
--------------------
%s
--------------------
stdout:
--------------------
%s`, err, syserr.String(), sysout.String())
	}

	return sysout.String()
}

func newGoldie(t *testing.T, suffix string) *goldie.Goldie {
	return goldie.New(t,
		goldie.WithTestNameForDir(true),
		goldie.WithNameSuffix(fmt.Sprintf(".golden%s", suffix)),
	)
}

func testGoCodeGeneration(t *testing.T, output string) {
	t.Helper()
	t.Run("generate go", func(t *testing.T) {
		t.Helper()
		g := newGoldie(t, ".go")
		g.Assert(t, "template", []byte(output))
	})
}

func testGoCodeExecution(t *testing.T, output string) {
	t.Helper()
	t.Run("compile and run", func(t *testing.T) {
		t.Helper()
		if testing.Short() {
			t.Skip()
		}
		stdout := execAndReturnStdOut(t, "simple", output)
		g := newGoldie(t, "")
		g.Assert(t, "output", []byte(stdout))
	})
}
