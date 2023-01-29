package errors

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type TokenizerError struct {
	Msg string
	Err error
}

func NewTokenizerError(input string, idx int, err error) *TokenizerError {
	var sb strings.Builder
	lineNo, line, newIdx := findLine(input, idx)
	linePrefix := fmt.Sprintf("line %d: ", lineNo)
	sb.WriteString(linePrefix)
	sb.WriteString(line + "\n")
	spaces := strings.Repeat(" ", len(linePrefix)+newIdx)
	sb.WriteString(spaces + "▲\n")
	sb.WriteString(spaces + "└── %s")

	s := fmt.Sprintf(sb.String(), err.Error())
	return &TokenizerError{Msg: s, Err: err}
}

func (e *TokenizerError) Error() string {
	return e.Msg
}

func (e *TokenizerError) Unwrap() error {
	return e.Err
}

func findLine(input string, idx int) (lineNo int, line string, newIdx int) {
	lineStartIdx := strings.LastIndex(input[:idx], "\n")
	if lineStartIdx == -1 {
		lineStartIdx = 0
	} else {
		lineStartIdx += 1
	}
	lineEndIdx := strings.Index(input[idx:], "\n")
	if lineEndIdx == -1 {
		lineEndIdx = len(input)
	} else {
		lineEndIdx = idx + lineEndIdx
	}
	lineNo = strings.Count(input[:lineStartIdx], "\n") + 1
	line = input[lineStartIdx:lineEndIdx]
	newIdx = idx - lineStartIdx
	for _, r := range []rune(line[:newIdx]) {
		runeLen := utf8.RuneLen(r)
		if runeLen > 1 {
			newIdx -= runeLen - 1
		}
	}
	return lineNo, line, newIdx
}
