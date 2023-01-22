package token

import (
	"testing"

	"github.com/BestFriendChris/go-ic/ic"
)

func TestToken_String(t *testing.T) {
	c := ic.New(t)

	var data any = "extra-data"
	c.PT([]struct {
		Name string
		Tok  *Token
	}{
		{"with S no E", &Token{TTnl, "\n", nil}},
		{"no S no E", &Token{TTcustom, "", nil}},
		{"with S with E", &Token{TTcustom, "foo", &data}},
		{"no S with E", &Token{TTcustom, "", &data}},
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

func TestNewTokenMarker(t *testing.T) {
	c := ic.New(t)
	c.PrintSection("NewTokenMarker has empty .S field")
	c.PV(NewTokenMarker(TTcustom))
	c.Expect(`
		################################################################################
		# NewTokenMarker has empty .S field
		################################################################################
		Token.TT: TT.Custom
		Token.S: ""
		Token.E: 
		`)
}
