package interfaces

type Macro interface {
	Tokenizer
}

type MapMacros = map[string]Macro

func NewMapMacros() MapMacros {
	return make(MapMacros)
}
