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

	in.Shift('}')
	tokens = append(tokens, token.NewToken(token.TTcodeLocalBlock, "}"))

	return tokens, nil
}

func (m *MacroIf) Parse(_ interfaces.TemplateHandler, toks []*token.Token) (rest []*token.Token, err error) {
	// No extra work to do
	return toks, nil
}
