package tokenizer

import (
	"fmt"
	"github.com/BestFriendChris/lozenge/internal/logic/token"
	"strings"
	"unicode"

	"github.com/BestFriendChris/lozenge/internal/logic/macro"
)

type ContentTokenizer struct {
	loz    rune
	macros macro.MapMacros
}

func NewDefault(macros macro.MapMacros) *ContentTokenizer {
	return New('◊', macros)
}

func New(loz rune, macros macro.MapMacros) *ContentTokenizer {
	return &ContentTokenizer{
		loz:    loz,
		macros: macros,
	}
}

func (ct *ContentTokenizer) ReadTokensUntil(rest, stopAt string) ([]*token.Token, string) {
	tokens := make([]*token.Token, 0)
	var tok *token.Token
	for {
		if len(rest) == 0 {
			break
		}
		if stopAt != "" && strings.HasPrefix(rest, stopAt) {
			break
		}

		prev := rest
		tok, rest = ct.NextToken(rest)
		if tok.TT == token.TTmacro {
			m, found := ct.macros[tok.S]
			if found {
				tokens = append(tokens, tok)
				var nextTokens []*token.Token
				nextTokens, rest = m.NextTokens(rest)
				for _, nextToken := range nextTokens {
					tokens = append(tokens, nextToken)
				}
			} else {

				runes := []rune(prev)[1:]
				tok, rest = ct.lozengeFallback(runes)
				tokens = append(tokens, tok)
				//rest = string(runes[1:])
				fmt.Printf("new rest: %q\n", rest)
				continue
			}
		} else {
			tokens = append(tokens, tok)
		}
	}
	return tokens, rest
}

func (ct *ContentTokenizer) NextToken(s string) (*token.Token, string) {
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
	return token.NewToken(tt, s[:endIdx]), s[endIdx:]
}

func (ct *ContentTokenizer) NextTokenGoCodeUntilOpenBraceLoz(s string) (*token.Token, string) {
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
		return nil, s
	} else {
		return token.NewToken(token.TTgoCodeLocalBlock, string(runes[:endIdx])), string(runes[endIdx+1:])
	}
}

func (ct *ContentTokenizer) lozengeFallback(runes []rune) (*token.Token, string) {
	return token.NewToken(token.TTcontent, string(ct.loz)), string(runes)
}

func (ct *ContentTokenizer) parseLozenge(runes []rune) (*token.Token, string) {
	if len(runes) == 0 {
		return ct.lozengeFallback(runes)
	}
	switch runes[0] {
	case ' ', '\n':
		return ct.lozengeFallback(runes)
	case '{':
		return ct.ParseGoToClosingBrace(runes)
	case '(':
		return ct.ParseGoCodeFromTo(runes, token.TTgoCodeExpr, '(', ')', true)
	case '.':
		return ct.parseMacroIdentifier(runes)
	case '^':
		if runes[1] == '{' {
			return ct.ParseGoToClosingBrace(runes)
		} else {
			return ct.lozengeFallback(runes)
		}
	case ct.loz:
		return ct.lozengeFallback(runes[1:])
	default:
		identifier, rest := ct.readGoIdentifier(runes)
		if len(identifier) == 0 {
			return ct.lozengeFallback(runes)
		} else {
			return token.NewToken(token.TTgoCodeExpr, identifier), rest
		}
	}
}

func (ct *ContentTokenizer) parseMacroIdentifier(runes []rune) (*token.Token, string) {
	runesSkipDot := runes[1:]
	identifier, _ := ct.readGoIdentifier(runesSkipDot)
	if len(identifier) == 0 {
		return ct.lozengeFallback(runes)
	}
	return token.NewToken(token.TTmacro, identifier), string(runesSkipDot)
}

func (ct *ContentTokenizer) ParseGoToClosingBrace(runes []rune) (*token.Token, string) {
	var tt token.TokenType
	if runes[0] == '^' {
		tt = token.TTgoCodeGlobalBlock
		runes = runes[1:]
	} else {
		tt = token.TTgoCodeLocalBlock
	}
	return ct.ParseGoCodeFromTo(runes, tt, '{', '}', false)
}

func (ct *ContentTokenizer) ParseGoCodeFromTo(runes []rune, tt token.TokenType, open, close rune, keep bool) (*token.Token, string) {
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
		return ct.lozengeFallback(runes)
	} else if keep {
		s, rest = runes[0:endIdx], runes[endIdx:]
	} else {
		s, rest = runes[1:endIdx-1], runes[endIdx:]
	}
	return token.NewToken(tt, string(s)), string(rest)
}

func isGoLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func (ct *ContentTokenizer) readGoIdentifier(runes []rune) (string, string) {
	if !isGoLetter(runes[0]) {
		return "", string(runes)
	}
	var endIdx int
	for i, r := range runes {
		if isGoLetter(r) || unicode.IsNumber(r) {
			endIdx = i + 1
		} else {
			break
		}
	}
	return string(runes[:endIdx]), string(runes[endIdx:])
}
