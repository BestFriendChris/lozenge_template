package input

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/BestFriendChris/lozenge_template/internal/logic/errors"
)

type Input struct {
	str string
	idx int
}

func NewInput(s string) *Input {
	return &Input{str: s}
}

func (i *Input) Rest() string {
	return i.str[i.idx:]
}

func (i *Input) Consumed() bool {
	return i.idx >= len(i.str)
}

func (i *Input) Consume(r rune) bool {
	got, found := i.Peek()
	if found && got == r {
		i.Shift(got)
		return true
	} else {
		return false
	}
}

func (i *Input) ConsumeString(prefix string) bool {
	if i.HasPrefix(prefix) {
		i.idx += len(prefix)
		return true
	} else {
		return false
	}
}

func (i *Input) ConsumeRegexp(r *regexp.Regexp) (string, bool) {
	found := r.FindIndex([]byte(i.Rest()))
	if found != nil && found[0] == 0 {
		match := i.Rest()[:found[1]]
		i.idx += found[1]
		return match, true
	} else {
		return "", false
	}
}

func (i *Input) HasPrefix(prefix string) bool {
	return strings.HasPrefix(i.Rest(), prefix)
}

func (i *Input) HasPrefixRegexp(r *regexp.Regexp) bool {
	found := r.FindIndex([]byte(i.Rest()))
	return found != nil && found[0] == 0
}

func (i *Input) Peek() (r rune, found bool) {
	var size int
	r, size = utf8.DecodeRuneInString(i.Rest())
	if size == 0 {
		return utf8.RuneError, false
	}
	return r, true
}

func (i *Input) Shift(expected rune) {
	r, found := i.Peek()
	if !found {
		panic("nothing to skip")
	}
	if r != expected {
		panic(fmt.Sprintf("unable to skip '%c' (found '%c')", expected, r))
	}
	i.idx += utf8.RuneLen(r)
}

func (i *Input) Unshift(expected rune) {
	i.UnshiftString(string(expected))
}

func (i *Input) UnshiftString(expected string) {
	if !strings.HasSuffix(i.str[:i.idx], expected) {
		panic(fmt.Sprintf("unable to unshift %q: (found %q)", expected, i.str[i.idx-(len(expected)):i.idx]))
	}
	i.idx -= len(expected)
}

func (i *Input) ErrorHere(err error) error {
	return errors.NewTokenizerError(i.str, i.idx, err)
}

func (i *Input) ReadWhile(f func(r rune) bool) string {
	var sb bytes.Buffer
	for {
		r, found := i.Peek()
		if found && f(r) {
			sb.WriteRune(r)
			i.Shift(r)
		} else {
			return sb.String()
		}
	}
}

func (i *Input) TryReadWhile(f func(r rune, last bool) (bool, error)) (string, error) {
	var sb bytes.Buffer
	startIdx := i.idx
	for {
		r, found := i.Peek()
		if !found {
			break
		}
		test, err := f(r, i.isLast())
		if err != nil {
			i.idx = startIdx
			return "", i.ErrorHere(err)
		}
		if test {
			sb.WriteRune(r)
			i.Shift(r)
		} else {
			break
		}
		if i.Rest() == "" {
			break
		}
	}
	return sb.String(), nil
}

func (i *Input) isLast() bool {
	_, size := utf8.DecodeRuneInString(i.Rest())
	return len(i.str)-size == i.idx
}
