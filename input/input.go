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
	str     string
	idx     int
	lineNo  int
	lineIdx []int
}

func NewInput(s string) *Input {
	return &Input{
		str:     s,
		lineNo:  1,
		lineIdx: makeLineIdx(s),
	}
}

func (i *Input) Seek(idx int) {
	if idx < 0 {
		panic("unable to seek to negative idx")
	}
	strLen := len(i.str)
	if idx > strLen {
		panic(fmt.Sprintf("unable to seek past end of Input: %d vs %d", idx, strLen))
	}
	if i.idx+1 == idx {
		// Fast pass for moving to next character
		if i.str[i.idx] == '\n' {
			i.lineNo++
		}
		i.idx = idx
		return
	} else {
		for lineNo, endIdx := range i.lineIdx {
			if idx <= endIdx {
				i.lineNo = lineNo + 1
				i.idx = idx
				return
			}
		}
	}
	panic(fmt.Sprintf("unable to find line for idx %d", idx))
}

func (i *Input) SeekOffset(offset int) {
	i.Seek(i.idx + offset)
}

func (i *Input) Rest() string {
	return i.str[i.idx:]
}

func (i *Input) Pos() Pos {
	return Pos{i.idx, i.lineNo, i.col()}
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
		i.SeekOffset(len(prefix))
		return true
	} else {
		return false
	}
}

func (i *Input) ConsumeRegexp(r *regexp.Regexp) (string, bool) {
	found := r.FindIndex([]byte(i.Rest()))
	if found != nil && found[0] == 0 {
		match := i.Rest()[:found[1]]
		i.SeekOffset(found[1])
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
	i.SeekOffset(utf8.RuneLen(r))
}

func (i *Input) Unshift(expected rune) {
	i.UnshiftString(string(expected))
}

func (i *Input) UnshiftString(expected string) {
	if !strings.HasSuffix(i.str[:i.idx], expected) {
		panic(fmt.Sprintf("unable to unshift %q: (found %q)", expected, i.str[i.idx-(len(expected)):i.idx]))
	}
	i.SeekOffset(-len(expected))
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
			i.Seek(startIdx)
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

func makeLineIdx(s string) []int {
	lines := make([]int, 0)
	var idxSoFar int
	for {
		idx := strings.IndexRune(s[idxSoFar:], '\n')
		if idx == -1 {
			lines = append(lines, len(s))
			break
		}
		idxSoFar += idx + 1
		lines = append(lines, idxSoFar)
	}
	return lines
}

func (i *Input) col() int {
	var leftIdx int
	if i.lineNo > 1 {
		leftIdx = i.lineIdx[i.lineNo-2]
	}
	return (i.idx - leftIdx) + 1
}
