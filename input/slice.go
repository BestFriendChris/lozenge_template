package input

import "unicode/utf8"

type Slice struct {
	S          string
	Start, End Pos
}

func EmptySlice() Slice {
	return Slice{}
}

func NewSlice(s string, start Pos, end Pos) Slice {
	if start.Idx > end.Idx {
		panic("invalid state")
	}
	return Slice{S: s, Start: start, End: end}
}

func (slc Slice) String() string {
	return slc.S
}

func (slc Slice) Join(other Slice) Slice {
	return Slice{
		S:     slc.S + other.S,
		Start: slc.Start,
		End:   other.End,
	}
}

func (slc Slice) Len() int {
	return utf8.RuneCountInString(slc.S)
}
