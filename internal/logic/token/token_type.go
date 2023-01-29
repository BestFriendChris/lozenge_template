package token

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type TokenType int

func (t TokenType) String() string {
	customTT, found := globalRegistry.Lookup(t)
	if !found {
		panic(fmt.Sprintf("unrecognized TokenType with value %d", t))
	}
	return customTT
}

func (t TokenType) IsCustom() bool {
	return t >= TTcustom
}

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

func NewRegistry() *Registry {
	tr := Registry{idx: TTcustom}
	baseRegistry := map[TokenType]string{
		TTunknown:         "TT.Unknown",
		TTws:              "TT.WS",
		TTnl:              "TT.NL",
		TTcontent:         "TT.Content",
		TTcodeGlobalBlock: "TT.CodeGlobalBlock",
		TTcodeLocalBlock:  "TT.CodeLocalBlock",
		TTcodeLocalExpr:   "TT.CodeLocalExpr",
		TTmacro:           "TT.Macro",
		TTcustom:          "TT.Custom",
	}
	tr.registry.Store(baseRegistry)
	return &tr
}

type Registry struct {
	mu       sync.Mutex
	idx      int // only change with mutex
	registry atomic.Value
}

func (tr *Registry) Lookup(tt TokenType) (string, bool) {
	m := tr.registry.Load().(map[TokenType]string)
	s, found := m[tt]
	return s, found
}

func (tr *Registry) Register(name string) TokenType {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	tr.idx++
	tt := TokenType(tr.idx)
	m1 := tr.registry.Load().(map[TokenType]string)
	m2 := make(map[TokenType]string)
	for k, v := range m1 {
		m2[k] = v
	}
	m2[tt] = `TT.` + name
	tr.registry.Store(m2)
	return tt
}

var globalRegistry = NewRegistry()

func RegisterCustomTokenType(name string) TokenType {
	return globalRegistry.Register(name)
}
