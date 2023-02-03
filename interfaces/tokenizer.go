package interfaces

import (
	"github.com/BestFriendChris/lozenge_template/input"
	"github.com/BestFriendChris/lozenge_template/internal/logic/token"
)

type Tokenizer interface {
	NextTokens(ct ContentTokenizer, in *input.Input) (toks []*token.Token, err error)
}

type ContentTokenizer interface {
	NextTokenCodeUntilOpenBraceLoz(in *input.Input) (*token.Token, error)
	ReadTokensUntil(in *input.Input, stopAt string) ([]*token.Token, error)
	ParseGoCodeFromTo(in *input.Input, tt token.TokenType, open, close rune, keep bool) ([]*token.Token, error)
	ParseGoToClosingBrace(in *input.Input) ([]*token.Token, error)
}
