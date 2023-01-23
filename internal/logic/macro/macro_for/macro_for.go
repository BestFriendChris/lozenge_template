package macro_for

import (
	"strings"

	"github.com/BestFriendChris/lozenge/interfaces"
	"github.com/BestFriendChris/lozenge/internal/logic/token"
)

func New() *MacroFor {
	return &MacroFor{}
}

type MacroFor struct{}

func (m MacroFor) Name() string {
	return "for"
}

func (m MacroFor) NextTokens(ct interfaces.ContentTokenizer, input string) (toks []*token.Token, rest string) {
	rest = input
	tokens := make([]*token.Token, 0)

	var tok *token.Token
	tok, rest = ct.NextTokenCodeUntilOpenBraceLoz(rest)
	if tok == nil {
		return make([]*token.Token, 0), input
	}
	tokens = append(tokens, tok)

	var subTokens []*token.Token
	subTokens, rest = ct.ReadTokensUntil(rest, "◊}")
	rest = strings.TrimPrefix(rest, "◊")
	for _, subToken := range subTokens {
		tokens = append(tokens, subToken)
	}

	tokens = append(tokens, token.NewToken(token.TTcodeLocalBlock, "}"))

	rest = strings.TrimPrefix(rest, "}")

	return tokens, rest
}

func (m MacroFor) Parse(_ interfaces.TemplateHandler, toks []*token.Token) (rest []*token.Token, err error) {
	// No extra work to do
	return toks, nil
}
