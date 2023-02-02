package token

import (
	"github.com/BestFriendChris/lozenge_template/input"
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
)

func TestToken_String(t *testing.T) {
	c := ic.New(t)

	mkSlc := func(s string) input.Slice {
		return input.Slice{S: s}
	}
	var data any = "extra-data"
	c.PT([]struct {
		Name string
		Tok  *Token
	}{
		{"with S no E", &Token{TTnl, mkSlc("\n"), nil}},
		{"no S no E", &Token{TTcustom, mkSlc(""), nil}},
		{"with S with E", &Token{TTcustom, mkSlc("foo"), &data}},
		{"no S with E", &Token{TTcustom, mkSlc(""), &data}},
	})

	c.Expect(`
		   | Name            | Tok                            |
		---+-----------------+--------------------------------+
		 1 | "with S no E"   | TT.NL("\n")                    |
		---+-----------------+--------------------------------+
		 2 | "no S no E"     | TT.Custom                      |
		---+-----------------+--------------------------------+
		 3 | "with S with E" | TT.Custom("foo")["extra-data"] |
		---+-----------------+--------------------------------+
		 4 | "no S with E"   | TT.Custom["extra-data"]        |
		---+-----------------+--------------------------------+
		`)
}
