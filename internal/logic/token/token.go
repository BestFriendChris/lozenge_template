package token

import (
	"fmt"

	"github.com/BestFriendChris/lozenge_template/input"
)

type Token struct {
	TT TokenType
	S  string
	E  *any
}

func NewToken(tt TokenType, s string) *Token {
	return &Token{TT: tt, S: s}
}

func NewTokenMarker(tt TokenType) *Token {
	return &Token{TT: tt}
}

func (t Token) String() string {
	var str, extra string
	if t.S != "" {
		str = fmt.Sprintf("(%q)", t.S)
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

type TokenSlice struct {
	TT  TokenType
	Slc input.Slice
	E   *any
}

func NewTokenSlice(tt TokenType, s input.Slice) *TokenSlice {
	return &TokenSlice{TT: tt, Slc: s}
}

func (t TokenSlice) String() string {
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

func (t TokenSlice) ToToken() *Token {
	return NewToken(t.TT, t.Slc.S)
}
