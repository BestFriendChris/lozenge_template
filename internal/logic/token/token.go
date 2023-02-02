package token

import (
	"fmt"

	"github.com/BestFriendChris/lozenge_template/input"
)

type Token struct {
	TT  TokenType
	Slc input.Slice
	E   *any
}

func NewToken(tt TokenType, s input.Slice) *Token {
	return &Token{TT: tt, Slc: s}
}

func (t Token) String() string {
	var str, extra string
	if t.Slc.S != "" {
		str = fmt.Sprintf("(%q)", t.Slc.S)
	}
	if t.E != nil {
		if stringer, ok := (*t.E).(fmt.Stringer); ok {
			extra = fmt.Sprintf("[%s]", stringer.String())
		} else {
			extra = fmt.Sprintf("[%#v]", *t.E)
		}
	}
	return fmt.Sprintf("%s%s%s", t.TT, str, extra)
}
