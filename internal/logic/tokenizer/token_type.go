package tokenizer

type TokenType int

const (
	Unknown TokenType = iota
	WS
	NL
	Content
	GoCodeGlobalBlock
	GoCodeLocalBlock
	GoCodeExpr
	Macro

	Custom = 999
)

func (t TokenType) String() string {
	switch t {
	case Unknown:
		return "TT.Unknown"
	case WS:
		return "TT.WS"
	case NL:
		return "TT.NL"
	case Content:
		return "TT.Content"
	case GoCodeGlobalBlock:
		return "TT.GoCodeGlobalBlock"
	case GoCodeLocalBlock:
		return "TT.GoCodeLocalBlock"
	case GoCodeExpr:
		return "TT.GoCodeExpr"
	case Macro:
		return "TT.Macro"
	case Custom:
		return "TT.Custom"
	default:
		panic("unrecognized TokenType")
	}
}
