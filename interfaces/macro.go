package interfaces

type Macro interface {
	Name() string
	Tokenizer
	Parser
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
	for name, m := range ms.mm {
		newMacros.mm[name] = m
	}
	if other == nil {
		return newMacros
	}
	for name, m := range other.mm {
		newMacros.mm[name] = m
	}
	return newMacros
}

func (ms *Macros) Get(name string) (m Macro, found bool) {
	m, found = ms.mm[name]
	return
}

func (ms *Macros) Known() []string {
	var keys []string
	for name := range ms.mm {
		keys = append(keys, name)
	}
	return keys
}
