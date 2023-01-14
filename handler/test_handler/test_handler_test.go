package test_handler

import (
	"reflect"
	"testing"

	"github.com/andreyvit/diff"
	"mvdan.cc/gofumpt/format"
)

func Test_WriteContent(t *testing.T) {
	th := TestHandler{}
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
		got, _ := th.Done(format.Options{})
		want := wantFormatted(
			``,
			`buf.WriteString("foo")
buf.WriteString("bar")`,
		)
		assertEqualDiff(t, got, want)
	}
}

func Test_WriteCodeInline(t *testing.T) {
	th := TestHandler{}
	th.WriteCodeInline("foo")
	th.WriteCodeInline("bar")

	got, _ := th.Done(format.Options{})
	want := wantFormatted(
		``,
		`buf.WriteString(fmt.Sprintf("%v", foo))
buf.WriteString(fmt.Sprintf("%v", bar))`,
	)
	assertEqualDiff(t, got, want)
}

func Test_WriteCodeBlock(t *testing.T) {
	th := TestHandler{}
	th.WriteCodeBlock("foo := 1")
	th.WriteCodeBlock("bar := 2")

	got, _ := th.Done(format.Options{})
	want := wantFormatted(
		``,
		`foo := 1
bar := 2`,
	)
	assertEqualDiff(t, got, want)
}

func Test_WriteCodeGlobalBlock(t *testing.T) {
	th := TestHandler{}
	th.WriteCodeGlobalBlock("var foo = 1")
	th.WriteCodeGlobalBlock("var bar = 2")

	got, _ := th.Done(format.Options{})
	want := wantFormatted(
		`var foo = 1
var bar = 2`,
		``,
	)
	assertEqualDiff(t, got, want)
}

func wantFormatted(global, inline string) string {
	s := STATIC[0] + global + STATIC[1] + inline + STATIC[2]

	formatted, err := format.Source([]byte(s), format.Options{})
	if err != nil {
		panic(err)
	}
	return string(formatted)
}

func assertEqualDiff(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("diff:\n-got\n+want\n=====\n%v", diff.LineDiff(got, want))
	}
}
