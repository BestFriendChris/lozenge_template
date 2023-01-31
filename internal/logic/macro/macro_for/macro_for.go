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
	in.ShiftSlc('◊')

	tok = token.NewTokenSlice(token.TTcodeLocalBlock, in.ShiftSlc('}')).ToToken()
	tokens = append(tokens, tok)

	return tokens, nil
}

func (m MacroFor) NextTokensSlc(ct interfaces.ContentTokenizerSlc, in *input.Input) (toks []*token.TokenSlice, err error) {
	tokens := make([]*token.TokenSlice, 0)

	var tok *token.TokenSlice
	tok, err = ct.NextTokenCodeUntilOpenBraceLozSlc(in)
	if err != nil {
		return nil, err
	}
	tokens = append(tokens, tok)

	var subTokens []*token.TokenSlice
	subTokens, err = ct.ReadTokensUntilSlc(in, "◊}")
	if err != nil {
		return nil, err
	}
	for _, subToken := range subTokens {
		tokens = append(tokens, subToken)
	}
	in.ShiftSlc('◊')

	tok = token.NewTokenSlice(token.TTcodeLocalBlock, in.ShiftSlc('}'))
	tokens = append(tokens, tok)

	return tokens, nil
}

func (m MacroFor) Parse(_ interfaces.TemplateHandler, toks []*token.Token) (rest []*token.Token, err error) {
	// No extra work to do
	return toks, nil
}

func (m MacroFor) ParseSlc(_ interfaces.TemplateHandler, toks []*token.TokenSlice) (rest []*token.TokenSlice, err error) {
	return toks, nil
}
