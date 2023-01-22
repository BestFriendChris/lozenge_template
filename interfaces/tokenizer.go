package interfaces

import "github.com/BestFriendChris/lozenge/internal/logic/token"

type Tokenizer interface {
	NextTokens(input string) (toks []*token.Token, rest string)
}
