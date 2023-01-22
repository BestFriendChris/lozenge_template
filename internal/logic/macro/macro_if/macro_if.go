package macro_if

import (
	"github.com/BestFriendChris/lozenge/internal/logic/token"
	"github.com/BestFriendChris/lozenge/internal/logic/tokenizer"
	"regexp"
	"strings"
)

func New(ct *tokenizer.ContentTokenizer) *MacroIf {
	return &MacroIf{
		ct: ct,
	}
}

type MacroIf struct {
	ct *tokenizer.ContentTokenizer
}

var (
	elseIfRegex = regexp.MustCompile(`}\s*else\s*if\s`)
	elseRegex   = regexp.MustCompile(`}\s*else\s*`)
)

func (m *MacroIf) NextTokens(input string) ([]*token.Token, string) {
	rest := input
	tokens := make([]*token.Token, 0)

	var tok *token.Token
	tok, rest = m.ct.NextTokenCodeUntilOpenBraceLoz(rest)
	if tok == nil {
		return make([]*token.Token, 0), input
	}
	tokens = append(tokens, tok)

	var subTokens []*token.Token
	subTokens, rest = m.ct.ReadTokensUntil(rest, "◊}")
	rest = strings.TrimPrefix(rest, "◊")
	for _, subToken := range subTokens {
		tokens = append(tokens, subToken)
	}

	for {
		found := elseIfRegex.FindIndex([]byte(rest))
		if found != nil && found[0] == 0 {
			tok, rest = m.ct.NextTokenCodeUntilOpenBraceLoz(rest)
			if tok == nil {
				return make([]*token.Token, 0), input
			}
			tokens = append(tokens, tok)

			subTokens, rest = m.ct.ReadTokensUntil(rest, "◊}")
			rest = strings.TrimPrefix(rest, "◊")
			for _, subToken := range subTokens {
				tokens = append(tokens, subToken)
			}
		} else {
			break
		}
	}
	found := elseRegex.FindIndex([]byte(rest))
	if found != nil && found[0] == 0 {
		tok, rest = m.ct.NextTokenCodeUntilOpenBraceLoz(rest)
		if tok == nil {
			return make([]*token.Token, 0), input
		}
		tokens = append(tokens, tok)

		subTokens, rest = m.ct.ReadTokensUntil(rest, "◊}")
		rest = strings.TrimPrefix(rest, "◊")
		for _, subToken := range subTokens {
			tokens = append(tokens, subToken)
		}
	}

	tokens = append(tokens, token.NewToken(token.TTcodeLocalBlock, "}"))

	rest = strings.TrimPrefix(rest, "}")

	return tokens, rest
}
