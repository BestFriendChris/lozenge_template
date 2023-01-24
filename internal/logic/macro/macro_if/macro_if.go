package macro_if

import (
	"fmt"
	"regexp"
	"strings"

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

func (m *MacroIf) NextTokens(ct interfaces.ContentTokenizer, input string) ([]*token.Token, string, error) {
	rest := input
	tokens := make([]*token.Token, 0)

	var tok *token.Token
	var err error
	tok, rest, err = ct.NextTokenCodeUntilOpenBraceLoz(rest)
	if err != nil {
		return nil, "", fmt.Errorf("macro(if): %w", err)
	}
	tokens = append(tokens, tok)

	var subTokens []*token.Token
	subTokens, rest, err = ct.ReadTokensUntil(rest, "◊}")
	if err != nil {
		return nil, "", fmt.Errorf("macro(if): %w", err)
	}
	rest = strings.TrimPrefix(rest, "◊")
	for _, subToken := range subTokens {
		tokens = append(tokens, subToken)
	}

	for {
		found := elseIfRegex.FindIndex([]byte(rest))
		if found != nil && found[0] == 0 {
			tok, rest, err = ct.NextTokenCodeUntilOpenBraceLoz(rest)
			if err != nil {
				return nil, "", fmt.Errorf("macro(if): %w", err)
			}
			tokens = append(tokens, tok)

			subTokens, rest, err = ct.ReadTokensUntil(rest, "◊}")
			if err != nil {
				return nil, "", fmt.Errorf("macro(if): %w", err)
			}
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
		tok, rest, err = ct.NextTokenCodeUntilOpenBraceLoz(rest)
		if err != nil {
			return nil, "", fmt.Errorf("macro(if): %w", err)
		}
		tokens = append(tokens, tok)

		subTokens, rest, err = ct.ReadTokensUntil(rest, "◊}")
		if err != nil {
			return nil, "", fmt.Errorf("macro(if): %w", err)
		}
		rest = strings.TrimPrefix(rest, "◊")
		for _, subToken := range subTokens {
			tokens = append(tokens, subToken)
		}
	}

	tokens = append(tokens, token.NewToken(token.TTcodeLocalBlock, "}"))

	rest = strings.TrimPrefix(rest, "}")

	return tokens, rest, nil
}

func (m *MacroIf) Parse(_ interfaces.TemplateHandler, toks []*token.Token) (rest []*token.Token, err error) {
	// No extra work to do
	return toks, nil
}
