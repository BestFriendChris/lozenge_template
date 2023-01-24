package macro_for

import (
	"fmt"
	"strings"

	"github.com/BestFriendChris/lozenge_template/interfaces"
	"github.com/BestFriendChris/lozenge_template/internal/logic/token"
)

func New() *MacroFor {
	return &MacroFor{}
}

type MacroFor struct{}

func (m MacroFor) Name() string {
	return "for"
}

func (m MacroFor) NextTokens(ct interfaces.ContentTokenizer, input string) (toks []*token.Token, rest string, err error) {
	rest = input
	tokens := make([]*token.Token, 0)

	var tok *token.Token
	tok, rest, err = ct.NextTokenCodeUntilOpenBraceLoz(rest)
	if err != nil {
		return nil, "", fmt.Errorf("macro(for): %w", err)
	}
	tokens = append(tokens, tok)

	var subTokens []*token.Token
	subTokens, rest, err = ct.ReadTokensUntil(rest, "◊}")
	if err != nil {
		return nil, "", fmt.Errorf("macro(for): %w", err)
	}
	rest = strings.TrimPrefix(rest, "◊")
	for _, subToken := range subTokens {
		tokens = append(tokens, subToken)
	}

	tokens = append(tokens, token.NewToken(token.TTcodeLocalBlock, "}"))

	rest = strings.TrimPrefix(rest, "}")

	return tokens, rest, nil
}

func (m MacroFor) Parse(_ interfaces.TemplateHandler, toks []*token.Token) (rest []*token.Token, err error) {
	// No extra work to do
	return toks, nil
}
