package token

import (
	"fmt"
)

type Token struct {
	TT TokenType
	S  string
	E  *any
}

func NewToken(tt TokenType, s string) *Token {
	return &Token{TT: tt, S: s}
}

func (t *Token) String() string {
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
