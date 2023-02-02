package input

import (
	"fmt"
	"unicode/utf8"
)

type Slice struct {
	Name, S    string
	Start, End Pos
}

func EmptySlice() Slice {
	return Slice{}
}

func NewSlice(name string, s string, start Pos, end Pos) Slice {
	if start.Idx > end.Idx {
		panic("invalid state")
	}
	return Slice{Name: name, S: s, Start: start, End: end}
}

func (slc Slice) String() string {
	if slc.Start.Row == 0 {
		return "<empty slice>"
	}
	return fmt.Sprintf("%s:%d - %q", slc.Name, slc.Start.Row, slc.S)
}

func (slc Slice) CanJoin(other Slice) bool {
	return slc.Name == other.Name && slc.End.Idx == other.Start.Idx
}

func (slc Slice) Join(other Slice) Slice {
	if !slc.CanJoin(other) {
		panic("unable to join lines")
	}
	return NewSlice(slc.Name, slc.S+other.S, slc.Start, other.End)
}

func (slc Slice) Len() int {
	return utf8.RuneCountInString(slc.S)
}
