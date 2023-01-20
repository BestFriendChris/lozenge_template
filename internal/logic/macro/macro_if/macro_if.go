package macro_if

import (
	"fmt"
	"github.com/BestFriendChris/lozenge/internal/logic/token"
	"github.com/BestFriendChris/lozenge/internal/logic/tokenizer"
	"regexp"
	"strings"
)

var (
	TTifBlock       = token.RegisterCustomTokenType("If.Block")
	TTifElseIfBlock = token.RegisterCustomTokenType("If.ElseIfBlock")
	TTifElseBlock   = token.RegisterCustomTokenType("If.ElseBlock")
	TTifEnd         = token.RegisterCustomTokenType("If.End")
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
	var tokens, subTokens []*token.Token
	tokens = []*token.Token{
		token.NewToken(TTifBlock, ""),
	}

	var tok *token.Token
	tok, rest = m.ct.NextTokenGoCodeUntilOpenBraceLoz(rest)
	if tok == nil {
		return make([]*token.Token, 0), input
	}
	tokens = append(tokens, tok)

	subTokens, rest = m.ct.ReadTokensUntil(rest, "◊}")
	rest = strings.TrimPrefix(rest, "◊")
	for _, subToken := range subTokens {
		tokens = append(tokens, subToken)
	}

	for {
		found := elseIfRegex.FindIndex([]byte(rest))
		if found != nil && found[0] == 0 {
			tokens = append(tokens, token.NewToken(TTifElseIfBlock, ""))
			tok, rest = m.ct.NextTokenGoCodeUntilOpenBraceLoz(rest)
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
		tokens = append(tokens, token.NewToken(TTifElseBlock, ""))
		tok, rest = m.ct.NextTokenGoCodeUntilOpenBraceLoz(rest)
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

	tokens = append(tokens, token.NewToken(TTifEnd, ""))
	tokens = append(tokens, token.NewToken(token.TTgoCodeLocalBlock, "}"))
	fmt.Printf("END rest: %q\n", rest)
	rest = strings.TrimPrefix(rest, "}")

	return tokens, rest
}
