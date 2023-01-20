package tokenizer

import "github.com/BestFriendChris/lozenge/internal/logic/token"

type Tokenizer interface {
	NextToken(s string) (token.Token, string)
}
