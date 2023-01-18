package tokenizer

import (
	"unicode"
)

type ContentTokenizer struct {
	loz rune
}

func NewDefault() *ContentTokenizer {
	return New('◊')
}

func New(loz rune) *ContentTokenizer {
	return &ContentTokenizer{loz}
}

func (ct *ContentTokenizer) NextToken(s string) (*Token, string) {
	var tt TokenType
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
			if tt == Unknown {
				tt = WS
			}
			if tt != WS {
				break loop
			}
		case '\n':
			if tt == Unknown {
				tt = NL
			} else {
				break loop
			}
		case ct.loz:
			if tt == Unknown {
				return ct.parseLozenge(runes[i+1:])
			} else {
				break loop
			}
		default:
			if tt == Unknown {
				tt = Content
			}
			if tt != Content {
				break loop
			}
		}
	}
	return NewToken(tt, s[:endIdx]), s[endIdx:]
}

func (ct *ContentTokenizer) lozengeFallback(runes []rune) (*Token, string) {
	return NewToken(Content, string(ct.loz)), string(runes)
}

func (ct *ContentTokenizer) parseLozenge(runes []rune) (*Token, string) {
	if len(runes) == 0 {
		return ct.lozengeFallback(runes)
	}
	switch runes[0] {
	case ' ', '\n':
		return ct.lozengeFallback(runes)
	case '{':
		return ct.ParseGoToClosingBrace(runes)
	case '(':
		return ct.ParseGoCodeTo(runes, GoCodeExpr, '(', ')', true)
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
			return NewToken(GoCodeExpr, identifier), rest
		}
	}
}

func (ct *ContentTokenizer) parseMacroIdentifier(runes []rune) (*Token, string) {
	identifier, rest := ct.readGoIdentifier(runes[1:])
	if len(identifier) == 0 {
		return ct.lozengeFallback(runes)
	}
	return NewToken(Macro, identifier), rest
}

func (ct *ContentTokenizer) ParseGoToClosingBrace(runes []rune) (*Token, string) {
	var tt TokenType
	if runes[0] == '^' {
		tt = GoCodeGlobalBlock
		runes = runes[1:]
	} else {
		tt = GoCodeLocalBlock
	}
	return ct.ParseGoCodeTo(runes, tt, '{', '}', false)
}

func (ct *ContentTokenizer) ParseGoCodeTo(runes []rune, tt TokenType, open, close rune, keep bool) (*Token, string) {
	var openCloseCount, endIdx int
	var inString, inBackQuotes, escapeInQuote bool
	for i, r := range runes {
		if r == open {
			if inString || inBackQuotes {
				continue
			}

			openCloseCount += 1
		} else if r == close {
			if inString || inBackQuotes {
				continue
			}
			openCloseCount -= 1
			if openCloseCount == 0 {
				endIdx = i + 1
				break
			}
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
	var s, rest []rune
	if openCloseCount > 0 {
		return ct.lozengeFallback(runes)
	} else if keep {
		s, rest = runes[0:endIdx], runes[endIdx:]
	} else {
		s, rest = runes[1:endIdx-1], runes[endIdx:]
	}
	return NewToken(tt, string(s)), string(rest)
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
