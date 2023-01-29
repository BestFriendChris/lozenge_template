package macro_for

import (
	"github.com/BestFriendChris/lozenge_template/input"
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

func (m MacroFor) NextTokens(ct interfaces.ContentTokenizer, in *input.Input) (toks []*token.Token, err error) {
	tokens := make([]*token.Token, 0)

	var tok *token.Token
	tok, err = ct.NextTokenCodeUntilOpenBraceLoz(in)
	if err != nil {
		return nil, err
	}
	tokens = append(tokens, tok)

	var subTokens []*token.Token
	subTokens, err = ct.ReadTokensUntil(in, "◊}")
	if err != nil {
		return nil, err
	}
	for _, subToken := range subTokens {
		tokens = append(tokens, subToken)
	}
	in.Shift('◊')
	in.Shift('}')

	tokens = append(tokens, token.NewToken(token.TTcodeLocalBlock, "}"))

	return tokens, nil
}

func (m MacroFor) Parse(_ interfaces.TemplateHandler, toks []*token.Token) (rest []*token.Token, err error) {
	// No extra work to do
	return toks, nil
}
