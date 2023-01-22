package token

import "fmt"

type TokenType int

const (
	TTunknown TokenType = iota
	TTws
	TTnl
	TTcontent
	TTcodeGlobalBlock
	TTcodeLocalBlock
	TTcodeLocalExpr
	TTmacro

	TTcustom = 999
	// Any custom types should be > 999
)

var a = make(map[TokenType]string)
var customIdx = TTcustom

func RegisterCustomTokenType(s string) TokenType {
	customIdx++
	tt := TokenType(customIdx)
	a[tt] = `TT.` + s
	return tt
}

func (t TokenType) String() string {
	switch t {
	case TTunknown:
		return "TT.Unknown"
	case TTws:
		return "TT.WS"
	case TTnl:
		return "TT.NL"
	case TTcontent:
		return "TT.Content"
	case TTcodeGlobalBlock:
		return "TT.CodeGlobalBlock"
	case TTcodeLocalBlock:
		return "TT.CodeLocalBlock"
	case TTcodeLocalExpr:
		return "TT.CodeLocalExpr"
	case TTmacro:
		return "TT.Macro"
	case TTcustom:
		return "TT.Custom"
	default:
		customTT, found := a[t]
		if found {
			return customTT
		}
		panic(fmt.Sprintf("unrecognized TokenType with value %d", t))

	}
}
