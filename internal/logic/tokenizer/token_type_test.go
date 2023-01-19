package tokenizer

import (
	"fmt"
	"strings"
	"testing"
)

func TestRegisterCustom(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		myCustom := RegisterCustomTokenType("MyCustom")
		got := myCustom.String()
		want := "TT.MyCustom"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func TestTokenType_String(t *testing.T) {
	const maxTokenTypeId = 7
	t.Run("known cases", func(t *testing.T) {
		knownIds := []TokenType{
			TTcustom,
		}
		for i := 0; i <= maxTokenTypeId; i++ {
			knownIds = append(knownIds, TokenType(i))
		}
		for _, tokenType := range knownIds {
			got := tokenType.String()
			if !strings.HasPrefix(got, "TT.") {
				t.Errorf("unexpected token type %q", got)
			}
		}
	})
	t.Run("unset token type", func(t *testing.T) {
		invalidTokenType := TokenType(maxTokenTypeId + 1)
		wantMessage := fmt.Sprintf("unrecognized TokenType with value %d", invalidTokenType)
		assertPanicsWithMessage(t, wantMessage, func() {
			s := invalidTokenType.String()
			t.Errorf("got %q want to panic", s)
		})
	})
	t.Run("custom token type", func(t *testing.T) {
		customTokenType := RegisterCustomTokenType("MyCustom")
		got := customTokenType.String()
		want := "TT.MyCustom"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func assertPanicsWithMessage(t *testing.T, msg string, f func()) {
	t.Helper()
	defer func() {
		t.Helper()
		r := recover()
		if r == nil {
			t.Fatalf("The code did not panic")
		}
		if r != msg {
			t.Fatalf("\n got panic %q\nwant panic %q", r, msg)
		}
	}()
	f()
}
