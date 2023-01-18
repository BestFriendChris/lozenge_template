package parser

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/BestFriendChris/lozenge/internal/infra/go_format"
)

type Macro func(Parser, string) (string, error)

type Parser interface {
	ParseSubstring(rest string, s string) (string, error)
	Handler() Handler
}

type DefaultParser struct {
	H      Handler
	macros map[string]Macro
	config ParserConfig
}

func (p DefaultParser) Handler() Handler {
	return p.H
}

func New(h Handler, macros map[string]Macro, config ParserConfig) DefaultParser {
	return DefaultParser{
		H:      h,
		macros: defaultMacros(h.DefaultMacros(), macros),
		config: config,
	}
}

func defaultMacros(handlerMacros, overrideMacros map[string]Macro) map[string]Macro {
	ms := map[string]Macro{
		"if":  Macro(macroIf),
		"for": Macro(macroFor),
	}
	if handlerMacros != nil {
		for k, v := range handlerMacros {
			ms[k] = v
		}
	}
	if overrideMacros != nil {
		for k, v := range overrideMacros {
			ms[k] = v
		}
	}
	return ms
}

func (p DefaultParser) Parse(s string) (string, error) {
	_, err := p.ParseSubstring(s, "")
	if err != nil {
		return "", err
	}
	fullOutput, err := p.H.Done()
	if err != nil {
		return "", err
	}
	return go_format.Format(fullOutput)
}

func debug(s string, vals ...any) {
	if DEBUGLOG {
		fmt.Printf(s, vals...)
	}
}

func (p DefaultParser) ParseSubstring(s, stopAt string) (rest string, err error) {
	var stopIdx, nextIdx int

	stopIdx = -1
	rest = s
	if p.config.TrimSpaces {
		debug("BEFORE TrimLeft:\n%q\n", rest)
		rest = strings.TrimLeft(rest, " \t")
		debug("AFTER:\n%q\n\n", rest)
	}
	for {
		debug("LOOP:\n%q\n\n", rest)
		if stopAt != "" {
			stopIdx = strings.Index(rest, stopAt)
			// debug("stopIdx: %d\n", stopIdx)
		}
		nextIdx = strings.IndexAny(rest, "◊\n")
		// debug("nextIdx: %d\n", nextIdx)

		isNewline := nextIdx >= 0 && rest[nextIdx] == '\n'

		if stopIdx >= 0 && nextIdx >= 0 && stopIdx <= nextIdx {
			// debug("stopIdx <= nextIdx: %d <= %d\n", nextIdx, stopIdx)
			before, after := rest[:stopIdx], strings.TrimPrefix(rest[stopIdx:], stopAt)
			// debug("before:\n%s\nafter:\n%s\n", before, after)
			if strings.TrimSpace(before) != "" {
				if p.config.TrimSpaces {
					debug("BEFORE TrimRight:\n%q\n", before)
					before = strings.TrimRight(before, " \t")
					debug("AFTER:\n%q\n\n", before)
				}
				if isNewline {
					p.H.WriteContent(before + "\n")
				} else {
					p.H.WriteContent(before)
				}
			}
			rest = after
			if isNewline {
				continue
			}

			if p.config.TrimSpaces {
				debug("BEFORE TrimLeft:\n%q\n", rest)
				rest = strings.TrimLeft(rest, " \t")
				debug("AFTER:\n%q\n\n", rest)
			}
			break
		}

		if isNewline {
			before, after := rest[:nextIdx+1], rest[nextIdx+1:]

			if p.config.TrimSpaces && strings.TrimSpace(before) != "" {
				debug("BEFORE TrimRight:\n%q\n", before)
				before = strings.TrimRight(before, " \t\n")
				debug("AFTER:\n%q\n\n", before)
				p.H.WriteContent(before)
			} else if before != "" {
				p.H.WriteContent(before)
			}

			rest = after
			if p.config.TrimSpaces {
				debug("BEFORE TrimLeft:\n%q\n", rest)
				rest = strings.TrimLeft(rest, " \t")
				debug("AFTER:\n%q\n\n", rest)
			}
			continue
		} else if nextIdx >= 0 {
			before, after := rest[:nextIdx], strings.TrimPrefix(rest[nextIdx:], "◊")
			// debug("AHH before:\n%q\nAHH after:\n%q\n", before, after)
			if p.config.TrimSpaces && strings.TrimSpace(before) != "" {
				debug("BEFORE TrimRight:\n%q\n", before)
				before = strings.TrimRight(before, " \t\n")
				debug("AFTER:\n%q\n\n", before)
				p.H.WriteContent(before)
			} else if before != "" {
				p.H.WriteContent(before)
			}
			// debug("BAHH before:\n%q\nBAHH after:\n%q\n", before, after)
			after, err = p.ParseCode(after)
			// debug("CAHH after:\n%q\n", after)
			if err != nil {
				return
			}
			rest = after
			if p.config.TrimSpaces {
				debug("BEFORE TrimLeft:\n%q\n", rest)
				rest = strings.TrimLeft(rest, " \t")
				debug("AFTER:\n%q\n\n", rest)
			}
		} else {
			if strings.TrimSpace(rest) != "" {
				if p.config.TrimSpaces {
					debug("BEFORE TrimRight:\n%q\n", rest)
					rest = strings.TrimRight(rest, " \t\n")
					debug("AFTER:\n%q\n\n", rest)
				}
				p.H.WriteContent(rest)
			}
			break
		}

		// before, after, found := strings.Cut(s, "◊")
		// if before != "" {
		// 	p.H.WriteContent(before)
		// }
		// if !found {
		// 	break
		// }
		// after, err = p.ParseCode(after)
		// if err != nil {
		// 	return
		// }
		// s = after
	}
	return
}

func (p DefaultParser) ParseCode(s string) (string, error) {
	switch {
	case s[0] == ' ':
		// "◊ " => "◊ "
		p.H.WriteContent("◊ ")
		s = s[1:]

	case s[0] == '\n':
		// "◊\n" => "◊\n"
		p.H.WriteContent("◊\n")
		s = s[1:]

	case strings.HasPrefix(s, "◊"):
		// "◊◊" => "◊"
		p.H.WriteContent("◊")
		s = s[3:]

	case s[0] == '.':
		// Macro
		after, err := p.parsecodeMacro(s)
		if err != nil {
			return "", err
		}
		s = after

	case s[0] == '{':
		// Inline Block
		after, err := p.parseCodeInlineBlock(s)
		if err != nil {
			return "", err
		}
		s = after

	case s[0] == '(':
		// Inline Expr
		after, err := p.parseCodeInlineExpr(s)
		if err != nil {
			return "", err
		}
		s = after

	case len(s) > 1 && s[0:2] == "^{":
		// Global Block
		after, err := p.parsecodeGlobalblock(s)
		if err != nil {
			return "", err
		}
		s = after

	default:
		before, after := p.parsecodeIdentifier(s)
		p.H.WriteCodeExpression(before)
		s = after
	}
	return s, nil
}

func (p DefaultParser) parseCodeInlineBlock(s string) (string, error) {
	numOfOpenBraces := 0
	closeBraceIdx := -1
loop:
	for i, c := range s {
		switch c {
		case '{':
			numOfOpenBraces += 1
		case '}':
			numOfOpenBraces -= 1
			if numOfOpenBraces == 0 {
				closeBraceIdx = i
				break loop
			}
		}
	}
	if closeBraceIdx < 0 {
		return "", fmt.Errorf("unable to find closing brace in %q", s)
	}
	p.H.WriteCodeBlock(s[1:closeBraceIdx])
	return s[closeBraceIdx+1:], nil
}

func (p DefaultParser) parseCodeInlineExpr(s string) (string, error) {
	numOfOpenParens := 0
	closeParenIdx := -1
loop:
	for i, c := range s {
		switch c {
		case '(':
			numOfOpenParens += 1
		case ')':
			numOfOpenParens -= 1
			if numOfOpenParens == 0 {
				closeParenIdx = i + 1
				break loop
			}
		}
	}
	if closeParenIdx < 0 {
		return "", fmt.Errorf("unable to find closing paren in %q", s)
	}
	p.H.WriteCodeExpression(s[:closeParenIdx])
	return s[closeParenIdx:], nil
}

func (p DefaultParser) parsecodeGlobalblock(s string) (string, error) {
	numOfOpenBraces := 0
	closeBraceIdx := -1
loop:
	for i, c := range s {
		switch c {
		case '{':
			numOfOpenBraces += 1
		case '}':
			numOfOpenBraces -= 1
			if numOfOpenBraces == 0 {
				closeBraceIdx = i
				break loop
			}
		}
	}
	if closeBraceIdx < 0 {
		return "", fmt.Errorf("unable to find closing brace in %q", s)
	}
	p.H.WriteCodeGlobalBlock(s[2:closeBraceIdx])
	return s[closeBraceIdx+1:], nil
}

func (p DefaultParser) parsecodeMacro(s string) (string, error) {
	macroName, rest := p.parsecodeIdentifier(s[1:])
	macro, found := p.macros[macroName]
	if !found {
		return "", fmt.Errorf("unknown macro %q", macroName)
	}
	return macro(p, rest)
}

func (p DefaultParser) parsecodeIdentifier(s string) (identifier, rest string) {
	foundIdx := len(s)
	for i, c := range s {
		if c != '_' && !unicode.IsDigit(c) && !unicode.IsLetter(c) {
			foundIdx = i
			break
		}
	}
	identifier, rest = s[:foundIdx], s[foundIdx:]
	return
}

var DEBUGLOG bool

func macroIf(p Parser, s string) (string, error) {
	var (
		ifStmt, rest string
		found        bool
		err          error
	)
	ifStmt, rest, found = strings.Cut(s, "{")
	if !found {
		return "", fmt.Errorf("macro[if]: '{' not found")
	}
	p.Handler().WriteCodeBlock(fmt.Sprintf("if %s {", ifStmt))
	rest, err = p.ParseSubstring(rest, "◊}")
	if err != nil {
		return "", err
	}
	hasElseIf := regexp.MustCompile("^ *else +if.+{")
	hasElse := regexp.MustCompile("^ *else *{")
	for {
		if hasElseIf.MatchString(rest) || hasElse.MatchString(rest) {
			ifStmt, rest, _ = strings.Cut(rest, "{")
			p.Handler().WriteCodeBlock(fmt.Sprintf("}%s{", ifStmt))
			rest, err = p.ParseSubstring(rest, "◊}")
			if err != nil {
				return "", err
			}
		} else {
			break
		}
	}
	p.Handler().WriteCodeBlock("}")
	return rest, nil
}

func macroFor(p Parser, s string) (string, error) {
	var (
		forStmt, rest string
		found         bool
		err           error
	)
	forStmt, rest, found = strings.Cut(s, "{")
	if !found {
		return "", fmt.Errorf("macro[for]: '{' not found")
	}
	p.Handler().WriteCodeBlock(fmt.Sprintf("for %s {", forStmt))
	rest, err = p.ParseSubstring(rest, "◊}")
	if err != nil {
		return "", err
	}
	p.Handler().WriteCodeBlock("}")
	return rest, nil
}
