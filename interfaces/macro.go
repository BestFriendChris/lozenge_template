package interfaces

type Macro interface {
	Name() string
	Tokenizer
	TokenizerSlc
	Parser
	ParserSlc
}

type Macros struct {
	mm map[string]Macro
}

func NewMacros() *Macros {
	return &Macros{make(map[string]Macro)}
}

func (ms *Macros) Add(m Macro) {
	ms.mm[m.Name()] = m
}

func (ms *Macros) Merge(other *Macros) *Macros {
	newMacros := NewMacros()
	for _, m := range ms.mm {
		newMacros.Add(m)
	}
	if other == nil {
		return newMacros
	}
	for _, m := range other.mm {
		newMacros.Add(m)
	}
	return newMacros
}

func (ms *Macros) Get(name string) (m Macro, found bool) {
	m, found = ms.mm[name]
	return
}

func (ms *Macros) Known() []string {
	keys := make([]string, 0)
	for name := range ms.mm {
		keys = append(keys, name)
	}
	return keys
}
