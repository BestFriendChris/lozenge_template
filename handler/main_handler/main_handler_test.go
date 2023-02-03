package main_handler

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
	"github.com/BestFriendChris/lozenge_template/input"
	"github.com/BestFriendChris/lozenge_template/internal/infra/go_format"
)

func Test_WriteContent(t *testing.T) {
	th := MainHandler{}
	i := input.NewInput("test", "foo\nbar")
	th.WriteTextContent(nextSlice(i, "foo"))
	nextSlice(i, "\n")
	th.WriteTextContent(nextSlice(i, "bar"))

	{
		got := th.Content
		want := []string{`test:1 - "foo"`, `test:2 - "bar"`}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}

	got, _ := th.Done()
	got = formatCode(t, got)

	c := ic.New(t)
	c.PrintSection("Formatted code")
	c.Println(got)
	c.Expect(`
		################################################################################
		# Formatted code
		################################################################################
		package main
		
		import (
			"bytes"
			"fmt"
		)
		
		func main() {
			buf := new(bytes.Buffer)
		//line test:1
			buf.WriteString("foo")
		//line test:2
			buf.WriteString("bar")
			fmt.Print(buf.String())
		}
		
		`)
}

func Test_WriteCodeExpression(t *testing.T) {
	th := MainHandler{}
	i := input.NewInput("test", "foo\nbar")
	th.WriteCodeLocalExpression(nextSlice(i, "foo"))
	_ = nextSlice(i, "\n")
	th.WriteCodeLocalExpression(nextSlice(i, "bar"))

	got, _ := th.Done()
	got = formatCode(t, got)

	c := ic.New(t)
	c.PrintSection("Formatted code")
	c.Println(got)
	c.Expect(`
		################################################################################
		# Formatted code
		################################################################################
		package main
		
		import (
			"bytes"
			"fmt"
		)
		
		func main() {
			buf := new(bytes.Buffer)
		//line test:1
			buf.WriteString(fmt.Sprintf("%v", foo))
		//line test:2
			buf.WriteString(fmt.Sprintf("%v", bar))
			fmt.Print(buf.String())
		}
		
		`)
}

func Test_WriteCodeBlock(t *testing.T) {
	th := MainHandler{}
	i := input.NewInput("test", "foo := 1\nbar := 2")
	th.WriteCodeLocalBlock(nextSlice(i, "foo := 1"))
	nextSlice(i, "\n")
	th.WriteCodeLocalBlock(nextSlice(i, "bar := 2"))

	got, _ := th.Done()
	got = formatCode(t, got)

	c := ic.New(t)
	c.PrintSection("Formatted code")
	c.Println(got)
	c.Expect(`
		################################################################################
		# Formatted code
		################################################################################
		package main
		
		import (
			"bytes"
			"fmt"
		)
		
		func main() {
			buf := new(bytes.Buffer)
		//line test:1
			foo := 1
		//line test:2
			bar := 2
			fmt.Print(buf.String())
		}
		
		`)
}

func Test_WriteCodeGlobalBlock(t *testing.T) {
	th := MainHandler{}
	i := input.NewInput("test", "var foo = 1\nvar bar = 2")
	th.WriteCodeGlobalBlock(nextSlice(i, "var foo = 1"))
	nextSlice(i, "\n")
	th.WriteCodeGlobalBlock(nextSlice(i, "var bar = 2"))

	got, _ := th.Done()
	got = formatCode(t, got)

	c := ic.New(t)
	c.PrintSection("Formatted code")
	c.Println(got)
	c.Expect(`
		################################################################################
		# Formatted code
		################################################################################
		package main
		
		import (
			"bytes"
			"fmt"
		)
		
		//line test:1
		var (
			foo = 1
		//line test:2
			bar = 2
		)
		
		func main() {
			buf := new(bytes.Buffer)
		
			fmt.Print(buf.String())
		}
		
		`)
}

func formatCode(t *testing.T, s string) string {
	formatted, err := go_format.Format(s)
	if err != nil {
		t.Fatal(err)
	}
	return formatted
}

func nextSlice(i *input.Input, prefix string) input.Slice {
	slice, found := i.ConsumeString(prefix)
	if !found {
		panic(fmt.Sprintf("incorrect prefix: %q\ninput: %q", prefix, i.Rest()))
	}
	return slice
}
