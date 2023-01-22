package interfaces

import (
	"github.com/BestFriendChris/lozenge/internal/logic/token"
)

type Parser interface {
	Parse(TemplateHandler, []*token.Token) (rest []*token.Token, err error)
}
