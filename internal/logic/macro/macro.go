package macro

import (
	"github.com/BestFriendChris/lozenge/internal/logic/token"
)

type Macro interface {
	NextTokens(s string) ([]*token.Token, string)
}

type MapMacros = map[string]Macro

func New() MapMacros {
	return make(MapMacros)
}
