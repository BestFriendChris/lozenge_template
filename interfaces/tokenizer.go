package interfaces

import "github.com/BestFriendChris/lozenge/internal/logic/token"

type Tokenizer interface {
	NextTokens(ct ContentTokenizer, input string) (toks []*token.Token, rest string, err error)
}

type ContentTokenizer interface {
	NextTokenCodeUntilOpenBraceLoz(s string) (*token.Token, string, error)
	ReadTokensUntil(input, stopAt string) ([]*token.Token, string, error)
	ParseGoCodeFromTo(runes []rune, tt token.TokenType, open, close rune, keep bool) (*token.Token, string, error)
	ParseGoToClosingBrace(runes []rune) (*token.Token, string, error)
}
