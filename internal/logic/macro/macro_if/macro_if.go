package macro_if

import (
	"regexp"

	"github.com/BestFriendChris/lozenge_template/input"
	"github.com/BestFriendChris/lozenge_template/interfaces"
	"github.com/BestFriendChris/lozenge_template/internal/logic/token"
)

func New() *MacroIf {
	return &MacroIf{}
}

type MacroIf struct {
}

var (
	elseIfRegex = regexp.MustCompile(`}\s*else\s*if\s`)
	elseRegex   = regexp.MustCompile(`}\s*else\s*`)
)

func (m *MacroIf) Name() string {
	return "if"
}

func (m *MacroIf) NextTokens(ct interfaces.ContentTokenizer, in *input.Input) (toks []*token.Token, err error) {
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
	for {
		found := in.HasPrefixRegexp(elseIfRegex)
		if !found {
			break
		}
		tok, err = ct.NextTokenCodeUntilOpenBraceLoz(in)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)

		subTokens, err = ct.ReadTokensUntil(in, "◊}")
		if err != nil {
			return nil, err
		}
		for _, subToken := range subTokens {
			tokens = append(tokens, subToken)
		}
		in.Shift('◊')
	}

	found := in.HasPrefixRegexp(elseRegex)
	if found {
		tok, err = ct.NextTokenCodeUntilOpenBraceLoz(in)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)

		subTokens, err = ct.ReadTokensUntil(in, "◊}")
		if err != nil {
			return nil, err
		}
		for _, subToken := range subTokens {
			tokens = append(tokens, subToken)
		}
		in.Shift('◊')
	}

	tok = token.NewTokenSlice(token.TTcodeLocalBlock, in.ShiftSlc('}')).ToToken()
	tokens = append(tokens, tok)

	return tokens, nil
}

func (m *MacroIf) NextTokensSlc(ct interfaces.ContentTokenizerSlc, in *input.Input) (toks []*token.TokenSlice, err error) {
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
	in.Shift('◊')
	for {
		found := in.HasPrefixRegexp(elseIfRegex)
		if !found {
			break
		}
		tok, err = ct.NextTokenCodeUntilOpenBraceLozSlc(in)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)

		subTokens, err = ct.ReadTokensUntilSlc(in, "◊}")
		if err != nil {
			return nil, err
		}
		for _, subToken := range subTokens {
			tokens = append(tokens, subToken)
		}
		in.Shift('◊')
	}

	found := in.HasPrefixRegexp(elseRegex)
	if found {
		tok, err = ct.NextTokenCodeUntilOpenBraceLozSlc(in)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)

		subTokens, err = ct.ReadTokensUntilSlc(in, "◊}")
		if err != nil {
			return nil, err
		}
		for _, subToken := range subTokens {
			tokens = append(tokens, subToken)
		}
		in.Shift('◊')
	}

	tok = token.NewTokenSlice(token.TTcodeLocalBlock, in.ShiftSlc('}'))
	tokens = append(tokens, tok)

	return tokens, nil
}

func (m *MacroIf) Parse(_ interfaces.TemplateHandler, toks []*token.Token) (rest []*token.Token, err error) {
	// No extra work to do
	return toks, nil
}

func (m *MacroIf) ParseSlc(_ interfaces.TemplateHandler, toks []*token.TokenSlice) (rest []*token.TokenSlice, err error) {
	return toks, nil
}
