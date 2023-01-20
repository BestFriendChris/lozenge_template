package macro_if

import (
	"reflect"
	"testing"

	"github.com/BestFriendChris/lozenge/internal/logic/macro"
	"github.com/BestFriendChris/lozenge/internal/logic/token"
	"github.com/BestFriendChris/lozenge/internal/logic/tokenizer"
)

func TestMacroIf_NextToken(t *testing.T) {
	t.Run("basic if", func(t *testing.T) {
		ct := tokenizer.NewDefault(macro.New())

		rest := `
if reflect.DeepEqual(val, []string{"foo"}) {◊
  hi
◊}bar`[1:]

		macroIf := New(ct)

		var tokens []*token.Token
		tokens, rest = macroIf.NextTokens(rest)

		if rest != "bar" {
			t.Errorf("\n got: %q\nwant: %q", rest, "bar")
		}

		tests := []struct {
			tt token.TokenType
			ts string
		}{
			{TTifBlock, ""},
			{token.TTgoCodeLocalBlock, "if reflect.DeepEqual(val, []string{\"foo\"}) {"},
			{token.TTnl, "\n"},
			{token.TTws, "  "},
			{token.TTcontent, "hi"},
			{token.TTnl, "\n"},
			{TTifEnd, ""},
			{token.TTgoCodeLocalBlock, "}"},
		}

		if len(tokens) != len(tests) {
			wantTokens := make([]*token.Token, len(tests))
			for i, test := range tests {
				wantTokens[i] = token.NewToken(test.tt, test.ts)
			}
			t.Fatalf("\n got: %v\nwant: %v", tokens, wantTokens)
		}

		for i, test := range tests {
			got := tokens[i]
			want := token.NewToken(test.tt, test.ts)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("\n got: %v\nwant: %v", got, want)
			}
		}
	})
	t.Run("basic if else", func(t *testing.T) {
		ct := tokenizer.NewDefault(macro.New())

		rest := `
if true {◊
  foo
◊}  else  {◊
  bar
◊}baz
`[1:]

		macroIf := New(ct)

		var tokens []*token.Token
		tokens, rest = macroIf.NextTokens(rest)

		if rest != "baz\n" {
			t.Errorf("\n got: %q\nwant: %q", rest, "baz\n")
		}

		tests := []struct {
			tt token.TokenType
			ts string
		}{
			{TTifBlock, ""},
			{token.TTgoCodeLocalBlock, "if true {"},
			{token.TTnl, "\n"},
			{token.TTws, "  "},
			{token.TTcontent, "foo"},
			{token.TTnl, "\n"},
			{TTifElseBlock, ""},
			{token.TTgoCodeLocalBlock, "}  else  {"},
			{token.TTnl, "\n"},
			{token.TTws, "  "},
			{token.TTcontent, "bar"},
			{token.TTnl, "\n"},
			{TTifEnd, ""},
			{token.TTgoCodeLocalBlock, "}"},
		}

		wantTokens := make([]*token.Token, len(tests))
		for i, test := range tests {
			wantTokens[i] = token.NewToken(test.tt, test.ts)
		}
		if len(tokens) != len(tests) {
			t.Fatalf("\n got: %v\nwant: %v", tokens, wantTokens)
		}

		for i, test := range tests {
			got := tokens[i]
			want := token.NewToken(test.tt, test.ts)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("\n got: %v\nwant: %v", got, want)
			}
		}
		if t.Failed() {
			t.Fatalf("\n got: %v\nwant: %v", tokens, wantTokens)
		}
	})
	t.Run("basic if else if", func(t *testing.T) {
		ct := tokenizer.NewDefault(macro.New())

		rest := `
if v == 1 {◊
  one
◊}  else  if v == 2 {◊
  two
◊}  else  if  v == 3 {◊
  three
◊}  else {◊
  four
◊}baz
`[1:]

		macroIf := New(ct)

		var tokens []*token.Token
		tokens, rest = macroIf.NextTokens(rest)

		if rest != "baz\n" {
			t.Errorf("\n got: %q\nwant: %q", rest, "baz\n")
		}

		tests := []struct {
			tt token.TokenType
			ts string
		}{
			{TTifBlock, ""},
			{token.TTgoCodeLocalBlock, "if v == 1 {"},
			{token.TTnl, "\n"},
			{token.TTws, "  "},
			{token.TTcontent, "one"},
			{token.TTnl, "\n"},
			{TTifElseIfBlock, ""},
			{token.TTgoCodeLocalBlock, "}  else  if v == 2 {"},
			{token.TTnl, "\n"},
			{token.TTws, "  "},
			{token.TTcontent, "two"},
			{token.TTnl, "\n"},
			{TTifElseIfBlock, ""},
			{token.TTgoCodeLocalBlock, "}  else  if  v == 3 {"},
			{token.TTnl, "\n"},
			{token.TTws, "  "},
			{token.TTcontent, "three"},
			{token.TTnl, "\n"},
			{TTifElseBlock, ""},
			{token.TTgoCodeLocalBlock, "}  else {"},
			{token.TTnl, "\n"},
			{token.TTws, "  "},
			{token.TTcontent, "four"},
			{token.TTnl, "\n"},
			{TTifEnd, ""},
			{token.TTgoCodeLocalBlock, "}"},
		}

		wantTokens := make([]*token.Token, len(tests))
		for i, test := range tests {
			wantTokens[i] = token.NewToken(test.tt, test.ts)
		}
		if len(tokens) != len(tests) {
			t.Fatalf("\n got: %v\nwant: %v", tokens, wantTokens)
		}

		for i, test := range tests {
			got := tokens[i]
			want := token.NewToken(test.tt, test.ts)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("\n got: %v\nwant: %v", got, want)
			}
		}
		if t.Failed() {
			t.Fatalf("\n got: %v\nwant: %v", tokens, wantTokens)
		}
	})
}

func TestMacroIf_NextToken_errorCases(t *testing.T) {
	tests := []struct {
		Name, input string
	}{
		{
			"no closing brace with if",
			`if true `,
		},
		{
			"no closing brace with else",
			`
if v == 1 {◊
  one
◊}  else 
  four
◊}baz
`[1:],
		},
		{
			"no closing brace with else if",
			`
if v == 1 {◊
  one
◊}  else  if v == 2 {◊
  two
◊}  else  if  v == 3 
  three
◊}  else {◊
  four
◊}baz
`[1:],
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ct := tokenizer.NewDefault(macro.New())
			macroIf := New(ct)

			tokens, rest := macroIf.NextTokens(test.input)

			if len(tokens) != 0 {
				t.Fatalf("got %d tokens, want 0\ntokens: %v", len(tokens), tokens)
			}

			want := test.input
			if rest != want {
				t.Errorf("\n got: %q\nwant: %q", rest, want)
			}
		})
	}
}
