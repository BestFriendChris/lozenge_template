package tokenizer

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/BestFriendChris/lozenge_template/input"
	"github.com/BestFriendChris/lozenge_template/interfaces"
	"github.com/BestFriendChris/lozenge_template/internal/logic/token"
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

func (ct *ContentTokenizer) ReadAll(in *input.Input) ([]*token.Token, error) {
	return ct.ReadTokensUntil(in, "")
}

func (ct *ContentTokenizer) ReadTokensUntil(in *input.Input, stopAt string) (tokens []*token.Token, err error) {
	tokens = make([]*token.Token, 0)
	var toks []*token.Token
	for {
		if in.Consumed() {
			break
		}

		if stopAt != "" && in.HasPrefix(stopAt) {
			return tokens, nil
		}

		toks, err = ct.NextTokens(in)
		if err != nil {
			return nil, err
		}
		for _, tok := range toks {
			tokens = append(tokens, tok)
		}
	}
	if stopAt != "" {
		return nil, in.ErrorHere(fmt.Errorf("did not find %q", stopAt))
	} else {
		return tokens, nil
	}
}

func (ct *ContentTokenizer) NextTokens(in *input.Input) ([]*token.Token, error) {
	var tt token.TokenType
	var foundLoz bool
	s, _ := in.TryReadWhile(func(r rune, last bool) (bool, error) {
		switch r {
		case ' ', '\t':
			if tt == token.TTunknown {
				tt = token.TTws
			}
			if tt != token.TTws {
				return false, nil
			}
		case '\n':
			if tt == token.TTunknown {
				tt = token.TTnl
			} else {
				return false, nil
			}
		case ct.loz:
			if tt == token.TTunknown {
				foundLoz = true
				in.Shift(r)
				return false, nil
			} else {
				return false, nil
			}
		default:
			if tt == token.TTunknown {
				tt = token.TTcontent
			}
			if tt != token.TTcontent {
				return false, nil
			}
		}
		return true, nil
	})

	if foundLoz {
		return ct.parseLozenge(in)
	} else {
		return []*token.Token{token.NewToken(tt, s)}, nil
	}
}

func (ct *ContentTokenizer) NextTokenCodeUntilOpenBraceLoz(in *input.Input) (*token.Token, error) {
	var inString, inBackQuotes, escapeInQuote, foundOpenBrace, foundLoz bool
	goCode, err := in.TryReadWhile(func(r rune, last bool) (bool, error) {
		if foundLoz {
			return false, nil
		}

		if r == '{' && !(inString || inBackQuotes) {
			foundOpenBrace = true
			return true, nil
		} else if r == ct.loz && !(inString || inBackQuotes) {
			if foundOpenBrace {
				if last {
					return false, nil
				}
				foundLoz = true
				return true, nil
			} else {
				return false, fmt.Errorf("no open brace found")
			}
		} else if r == '\\' && inString {
			if !escapeInQuote {
				escapeInQuote = true
				return true, nil
			}
		} else if r == '"' && !inBackQuotes {
			if !escapeInQuote {
				inString = !inString
			}
		} else if r == '`' && !inString {
			inBackQuotes = !inBackQuotes
		}
		if last {
			return false, fmt.Errorf("no open brace found")
		}
		escapeInQuote = false
		foundOpenBrace = false
		return true, nil
	})

	if err != nil {
		return nil, err
	} else {
		goCode = in.TrimSliceSuffix(goCode, "◊")
		return token.NewToken(token.TTcodeLocalBlock, goCode), nil
	}
}

func (ct *ContentTokenizer) parseLozenge(in *input.Input) ([]*token.Token, error) {
	wrap := func(tok *token.Token, err error) ([]*token.Token, error) {
		if err != nil {
			return nil, err
		} else {
			return []*token.Token{tok}, nil
		}
	}

	loz := token.NewToken(token.TTcontent, in.SliceOffset(-utf8.RuneLen(ct.loz)))
	singletonLoz := []*token.Token{loz}
	r, found := in.Peek()
	if !found {
		return singletonLoz, nil
	}
	switch r {
	case ' ', '\n':
		return singletonLoz, nil
	case '{':
		return wrap(ct.ParseGoCodeFromTo(in, token.TTcodeLocalBlock, '{', '}', false))
	case '(':
		return wrap(ct.ParseGoCodeFromTo(in, token.TTcodeLocalExpr, '(', ')', true))
	case '.':
		in.Shift(r)
		return ct.parseMacroIdentifier(loz, in)
	case '^':
		in.Shift(r)
		r, found = in.Peek()
		if found && r == '{' {
			return wrap(ct.ParseGoCodeFromTo(in, token.TTcodeGlobalBlock, '{', '}', false))
		} else {
			in.Unshift('^')
			return singletonLoz, nil
		}
	case ct.loz:
		in.Shift(ct.loz)
		return singletonLoz, nil
	default:
		identifier := ct.readIdentifier(in)
		if identifier.Len() == 0 {
			return singletonLoz, nil
		} else {
			return wrap(token.NewToken(token.TTcodeLocalExpr, identifier), nil)
		}
	}
}

func (ct *ContentTokenizer) parseMacroIdentifier(lozSlc *token.Token, in *input.Input) (tokens []*token.Token, err error) {
	identifier := ct.readIdentifier(in)
	if identifier.Len() == 0 {
		return []*token.Token{lozSlc}, nil
	}
	in.UnshiftSlice(identifier)
	m, found := ct.macros.Get(identifier.S)
	if found {
		tokens = []*token.Token{token.NewToken(token.TTmacro, identifier)}
		var nextTokens []*token.Token
		nextTokens, err = m.NextTokens(ct, in)
		if err != nil {
			return nil, err
		}
		for _, nextToken := range nextTokens {
			tokens = append(tokens, nextToken)
		}
	} else {
		return nil, fmt.Errorf("unknown macro %s", identifier)
	}
	return
}

func (ct *ContentTokenizer) ParseGoToClosingBrace(in *input.Input) (*token.Token, error) {
	var tt token.TokenType
	if in.Consume('^') {
		tt = token.TTcodeGlobalBlock
	} else {
		tt = token.TTcodeLocalBlock
	}
	return ct.ParseGoCodeFromTo(in, tt, '{', '}', false)
}

func (ct *ContentTokenizer) ParseGoCodeFromTo(in *input.Input, tt token.TokenType, open, close rune, keep bool) (*token.Token, error) {
	var openCloseCount int
	var foundBalance, inString, inBackQuotes, inChar, escapeInQuote bool
	goCode, err := in.TryReadWhile(func(r rune, last bool) (bool, error) {
		if foundBalance {
			return false, nil
		}
		if r == open && !(inString || inBackQuotes || inChar) {
			openCloseCount += 1
		} else if r == close && !(inString || inBackQuotes || inChar) {
			openCloseCount -= 1
			if openCloseCount == 0 {
				foundBalance = true
			}
		} else if r == '\\' && (inString || inChar) {
			if !escapeInQuote {
				escapeInQuote = true
				return true, nil
			}
		} else if r == '"' && !inBackQuotes {
			if !escapeInQuote {
				inString = !inString
			}
		} else if r == '`' && !inString {
			inBackQuotes = !inBackQuotes
		} else if r == '\'' && !(inString || inBackQuotes) {
			if !escapeInQuote {
				inChar = !inChar
			}
		}
		if escapeInQuote {
			escapeInQuote = false
		}
		if last && openCloseCount > 0 {
			return false, fmt.Errorf("did not find matched '%c'", close)
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	if !keep {
		goCode = in.SliceAt(goCode.Start.Idx+1, goCode.End.Idx-1)
	}
	return token.NewToken(tt, goCode), nil
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func (ct *ContentTokenizer) readIdentifier(in *input.Input) input.Slice {
	if r, found := in.Peek(); !found || !isLetter(r) {
		return input.EmptySlice()
	}
	var afterFirst bool
	ident := in.ReadWhile(func(r rune) bool {
		if !afterFirst {
			afterFirst = true
			return isLetter(r)
		} else {
			return isLetter(r) || unicode.IsNumber(r)
		}
	})
	return ident
}
