package lozenge_template

type ParserConfig struct {
	Loz        rune
	TrimSpaces bool
}

func NewParserConfig() ParserConfig {
	return ParserConfig{
		Loz: '◊',
	}
}

func (pc ParserConfig) WithTrimSpaces() ParserConfig {
	pc.TrimSpaces = true
	return pc
}

func (pc ParserConfig) WithMarker(c rune) ParserConfig {
	pc.Loz = c
	return pc
}
