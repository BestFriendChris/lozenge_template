package main_handler

import (
	"reflect"
	"testing"

	"github.com/BestFriendChris/lozenge/internal/infra/go_format"
	"github.com/andreyvit/diff"
)

func Test_WriteContent(t *testing.T) {
	th := MainHandler{}
	th.WriteContent("foo")
	th.WriteContent("bar")

	{
		got := th.Content
		want := []string{"foo", "bar"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}

	{
		got, _ := th.Done()
		got = formatCode(t, got)
		want := wantFormatted(t, ``, `
buf.WriteString("foo")
buf.WriteString("bar")`[1:])
		assertEqualDiff(t, got, want)
	}
}

func Test_WriteCodeExpression(t *testing.T) {
	th := MainHandler{}
	th.WriteCodeExpression("foo")
	th.WriteCodeExpression("bar")

	got, _ := th.Done()
	got = formatCode(t, got)
	want := wantFormatted(t, ``, `
buf.WriteString(fmt.Sprintf("%v", foo))
buf.WriteString(fmt.Sprintf("%v", bar))`[1:])
	assertEqualDiff(t, got, want)
}

func Test_WriteCodeBlock(t *testing.T) {
	th := MainHandler{}
	th.WriteCodeBlock("foo := 1")
	th.WriteCodeBlock("bar := 2")

	got, _ := th.Done()
	got = formatCode(t, got)
	want := wantFormatted(t, ``, `
foo := 1
bar := 2`[1:])
	assertEqualDiff(t, got, want)
}

func Test_WriteCodeGlobalBlock(t *testing.T) {
	th := MainHandler{}
	th.WriteCodeGlobalBlock("var foo = 1")
	th.WriteCodeGlobalBlock("var bar = 2")

	got, _ := th.Done()
	got = formatCode(t, got)
	want := wantFormatted(
		t,
		`
var foo = 1
var bar = 2`[1:],
		``)
	assertEqualDiff(t, got, want)
}

func wantFormatted(t *testing.T, global, inline string) string {
	s := STATIC[0] + global + STATIC[1] + inline + STATIC[2]

	return formatCode(t, s)
}

func formatCode(t *testing.T, s string) string {
	formatted, err := go_format.Format(s)
	if err != nil {
		t.Fatal(err)
	}
	return formatted
}

func assertEqualDiff(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("diff:\n-got\n+want\n=====\n%v", diff.LineDiff(got, want))
	}
}
