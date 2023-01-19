package tokenizer

import "fmt"

type Token struct {
	TT TokenType
	S  string
	E  *any
}

func NewToken(tt TokenType, s string) *Token {
	return &Token{TT: tt, S: s}
}

func (t *Token) String() string {
	var extra string
	if t.E != nil {
		if stringer, ok := (*t.E).(fmt.Stringer); ok {
			extra = stringer.String()
		}
	}
	return fmt.Sprintf("%s(%q)%s", t.TT, t.S, extra)
}
