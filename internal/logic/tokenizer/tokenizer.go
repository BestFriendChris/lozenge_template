package tokenizer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/BestFriendChris/lozenge/interfaces"
	"github.com/BestFriendChris/lozenge/internal/logic/token"
)

type ContentTokenizer struct {
	loz    rune
	macros *interfaces.Macros
}

func NewDefault(macros *interfaces.Macros) *ContentTokenizer {
	return New('◊', macros)
}

func New(loz rune, macros *interfaces.Macros) *ContentTokenizer {
	return &ContentTokenizer{
		loz:    loz,
		macros: macros,
	}
}

func (ct *ContentTokenizer) ReadAll(input string) ([]*token.Token, error) {
	toks, _, err := ct.ReadTokensUntil(input, "")
	return toks, err
}

func (ct *ContentTokenizer) ReadTokensUntil(input, stopAt string) (tokens []*token.Token, rest string, err error) {
	rest = input
	tokens = make([]*token.Token, 0)
	var toks []*token.Token
	for {
		if len(rest) == 0 {
			break
		}
		if stopAt != "" && strings.HasPrefix(rest, stopAt) {
			return tokens, rest, nil
		}

		toks, rest, err = ct.NextTokens(rest)
		if err != nil {
			return toks, rest, err
		}
		for _, tok := range toks {
			tokens = append(tokens, tok)
		}
	}
	if stopAt != "" {
		return nil, "", fmt.Errorf("did not find %q", stopAt)
	} else {
		return tokens, rest, nil
	}
}

func (ct *ContentTokenizer) NextTokens(s string) ([]*token.Token, string, error) {
	var tt token.TokenType
	var endIdx int
	runes := []rune(s)
loop:
	for i := 0; i < len(runes)+1; i++ {
		endIdx = i
		if i == len(runes) {
			// EOL
			break
		}
		c := runes[i]
		switch c {
		case ' ', '\t':
			if tt == token.TTunknown {
				tt = token.TTws
			}
			if tt != token.TTws {
				break loop
			}
		case '\n':
			if tt == token.TTunknown {
				tt = token.TTnl
			} else {
				break loop
			}
		case ct.loz:
			if tt == token.TTunknown {
				return ct.parseLozenge(runes[i+1:])
			} else {
				break loop
			}
		default:
			if tt == token.TTunknown {
				tt = token.TTcontent
			}
			if tt != token.TTcontent {
				break loop
			}
		}
	}
	return []*token.Token{token.NewToken(tt, s[:endIdx])}, s[endIdx:], nil
}

func (ct *ContentTokenizer) NextTokenCodeUntilOpenBraceLoz(s string) (*token.Token, string, error) {
	runes := []rune(s)
	var endIdx int
	var inString, inBackQuotes, escapeInQuote bool
	for i, r := range runes {
		if r == '{' {
			if inString || inBackQuotes {
				continue
			}

			if i+1 < len(runes) && runes[i+1] == ct.loz {
				endIdx = i + 1
				break
			}

		} else if r == ct.loz {
			if inString || inBackQuotes {
				continue
			}
			break
		} else if r == '\\' {
			if inString {
				if !escapeInQuote {
					escapeInQuote = true
					continue
				}
			}
		} else if r == '"' {
			if inBackQuotes {
				continue
			}
			if !escapeInQuote {
				inString = !inString
			}
		} else if r == '`' {
			if inString {
				continue
			}
			inBackQuotes = !inBackQuotes
		}
		if escapeInQuote {
			escapeInQuote = false
		}
	}
	if endIdx == 0 {
		return nil, "", fmt.Errorf("no open brace found")
	} else {
		return token.NewToken(token.TTcodeLocalBlock, string(runes[:endIdx])), string(runes[endIdx+1:]), nil
	}
}

func (ct *ContentTokenizer) lozengeFallback(runes []rune) (*token.Token, string, error) {
	return token.NewToken(token.TTcontent, string(ct.loz)), string(runes), nil
}

func (ct *ContentTokenizer) parseLozenge(runes []rune) ([]*token.Token, string, error) {
	singletonToks := func(t *token.Token, rest string, err error) ([]*token.Token, string, error) {
		return []*token.Token{t}, rest, err
	}
	if len(runes) == 0 {
		return singletonToks(ct.lozengeFallback(runes))
	}
	switch runes[0] {
	case ' ', '\n':
		return singletonToks(ct.lozengeFallback(runes))
	case '{':
		return singletonToks(ct.ParseGoToClosingBrace(runes))
	case '(':
		return singletonToks(ct.ParseGoCodeFromTo(runes, token.TTcodeLocalExpr, '(', ')', true))
	case '.':
		return ct.parseMacroIdentifier(runes)
	case '^':
		if runes[1] == '{' {
			return singletonToks(ct.ParseGoToClosingBrace(runes))
		} else {
			return singletonToks(ct.lozengeFallback(runes))
		}
	case ct.loz:
		return singletonToks(ct.lozengeFallback(runes[1:]))
	default:
		identifier, rest := ct.readIdentifier(runes)
		if len(identifier) == 0 {
			return singletonToks(ct.lozengeFallback(runes))
		} else {
			return singletonToks(token.NewToken(token.TTcodeLocalExpr, identifier), rest, nil)
		}
	}
}

func (ct *ContentTokenizer) parseMacroIdentifier(runes []rune) (tokens []*token.Token, rest string, err error) {
	runesSkipDot := runes[1:]
	identifier, _ := ct.readIdentifier(runesSkipDot)
	if len(identifier) == 0 {
		tok, s, _ := ct.lozengeFallback(runes)
		return []*token.Token{tok}, s, nil
	}
	m, found := ct.macros.Get(identifier)
	if found {
		tokens = []*token.Token{token.NewToken(token.TTmacro, identifier)}
		var nextTokens []*token.Token
		nextTokens, rest, _ = m.NextTokens(ct, string(runes[1:]))
		for _, nextToken := range nextTokens {
			tokens = append(tokens, nextToken)
		}
	} else {
		return make([]*token.Token, 0), "", fmt.Errorf("unknown macro %q", identifier)
	}
	return
}

func (ct *ContentTokenizer) ParseGoToClosingBrace(runes []rune) (*token.Token, string, error) {
	var tt token.TokenType
	if runes[0] == '^' {
		tt = token.TTcodeGlobalBlock
		runes = runes[1:]
	} else {
		tt = token.TTcodeLocalBlock
	}
	return ct.ParseGoCodeFromTo(runes, tt, '{', '}', false)
}

func (ct *ContentTokenizer) ParseGoCodeFromTo(runes []rune, tt token.TokenType, open, close rune, keep bool) (*token.Token, string, error) {
	var openCloseCount, endIdx int
	var inString, inBackQuotes, inChar, escapeInQuote bool
	for i, r := range runes {
		if r == open {
			if inString || inBackQuotes || inChar {
				continue
			}

			openCloseCount += 1
		} else if r == close {
			if inString || inBackQuotes || inChar {
				continue
			}
			openCloseCount -= 1
			if openCloseCount == 0 {
				endIdx = i + 1
				break
			}
		} else if r == '\\' {
			if inString || inChar {
				if !escapeInQuote {
					escapeInQuote = true
					continue
				}
			}
		} else if r == '"' {
			if inBackQuotes {
				continue
			}
			if !escapeInQuote {
				inString = !inString
			}
		} else if r == '`' {
			if inString {
				continue
			}
			inBackQuotes = !inBackQuotes
		} else if r == '\'' {
			if inString || inBackQuotes {
				continue
			}

			if !escapeInQuote {
				inChar = !inChar
			}
		}
		if escapeInQuote {
			escapeInQuote = false
		}
	}
	var s, rest []rune
	if openCloseCount > 0 {
		return nil, "", fmt.Errorf("did not find matched '%c'", close)
	} else if keep {
		s, rest = runes[0:endIdx], runes[endIdx:]
	} else {
		s, rest = runes[1:endIdx-1], runes[endIdx:]
	}
	return token.NewToken(tt, string(s)), string(rest), nil
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func (ct *ContentTokenizer) readIdentifier(runes []rune) (string, string) {
	if !isLetter(runes[0]) {
		return "", string(runes)
	}
	var endIdx int
	for i, r := range runes {
		if isLetter(r) || unicode.IsNumber(r) {
			endIdx = i + 1
		} else {
			break
		}
	}
	return string(runes[:endIdx]), string(runes[endIdx:])
}
