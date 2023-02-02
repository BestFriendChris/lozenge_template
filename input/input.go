package input

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/BestFriendChris/lozenge_template/internal/logic/errors"
)

type Input struct {
	name    string
	str     string
	idx     int
	lineNo  int
	lineIdx []int
}

func NewInput(name string, s string) *Input {
	return &Input{
		name:    name,
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

func (i *Input) RestSlice() Slice {
	return i.SliceAt(i.idx, len(i.str))
}

func (i *Input) Pos() Pos {
	return Pos{i.idx, i.lineNo, i.col()}
}

func (i *Input) PosAt(idx int) Pos {
	line, col := i.findLineAndCol(idx)
	return Pos{idx, line, col}
}

func (i *Input) SliceOffset(offset int) Slice {
	var from, to int
	if offset < 0 {
		from, to = i.idx+offset, i.idx
	} else {
		from, to = i.idx, i.idx+offset
	}
	return i.SliceAt(from, to)
}

func (i *Input) SliceAt(from, to int) Slice {
	return NewSlice(i.name, i.str[from:to], i.PosAt(from), i.PosAt(to))
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

func (i *Input) ConsumeString(prefix string) (Slice, bool) {
	if i.HasPrefix(prefix) {
		slc := i.SliceOffset(len(prefix))
		i.SeekOffset(len(prefix))
		return slc, true
	} else {
		return EmptySlice(), false
	}
}

func (i *Input) ConsumeRegexp(r *regexp.Regexp) (Slice, bool) {
	found := r.FindIndex([]byte(i.Rest()))
	if found != nil && found[0] == 0 {
		match := i.SliceOffset(found[1])
		i.SeekOffset(found[1])
		return match, true
	} else {
		return EmptySlice(), false
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
	i.ShiftSlice(expected)
}

func (i *Input) ShiftSlice(expected rune) Slice {
	r, found := i.Peek()
	if !found {
		panic("nothing to skip")
	}
	if r != expected {
		panic(fmt.Sprintf("unable to skip '%c' (found '%c')", expected, r))
	}
	s := i.SliceOffset(1)
	i.SeekOffset(utf8.RuneLen(r))
	return s
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

func (i *Input) UnshiftSlice(expected Slice) {
	if !strings.HasSuffix(i.str[:i.idx], expected.S) {
		panic(fmt.Sprintf("unable to unshift %q: (found %q)", expected, i.str[i.idx-(expected.Len()):i.idx]))
	}
	i.SeekOffset(-expected.Len())
}

func (i *Input) ErrorHere(err error) error {
	return errors.NewTokenizerError(i.str, i.idx, err)
}

func (i *Input) ReadWhile(f func(r rune) bool) Slice {
	from := i.idx
	for {
		r, found := i.Peek()
		if found && f(r) {
			i.Shift(r)
		} else {
			return i.SliceAt(from, i.idx)
		}
	}
}

func (i *Input) TryReadWhile(f func(r rune, last bool) (bool, error)) (Slice, error) {
	startIdx := i.idx
	for {
		r, found := i.Peek()
		if !found {
			break
		}
		test, err := f(r, i.isLast())
		if err != nil {
			i.Seek(startIdx)
			return EmptySlice(), i.ErrorHere(err)
		}
		if test {
			i.Shift(r)
		} else {
			break
		}
		if i.Rest() == "" {
			break
		}
	}
	return i.SliceAt(startIdx, i.idx), nil
}

func (i *Input) TrimSliceSuffix(slc Slice, suffix string) Slice {
	if !strings.HasSuffix(slc.S, suffix) {
		return slc
	}
	return i.SliceAt(slc.Start.Idx, slc.End.Idx-len(suffix))
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

func (i *Input) findLineAndCol(idx int) (line, col int) {
	line, col = -1, -1
	for lineNo, endIdx := range i.lineIdx {
		if idx < endIdx || idx == len(i.str) {
			line = lineNo + 1
			break
		}
	}
	if line == -1 {
		panic(fmt.Sprintf("unable to find line for idx %d", idx))
	}

	var leftIdx int
	if line > 1 {
		leftIdx = i.lineIdx[line-2]
	}
	col = (idx - leftIdx) + 1
	return line, col
}
