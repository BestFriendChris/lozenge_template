package input

import "fmt"

type Pos struct {
	Idx, Row, Col int
}

func (p Pos) String() string {
	return fmt.Sprintf("Pos[line=%d;col=%d]", p.Row, p.Col)
}
