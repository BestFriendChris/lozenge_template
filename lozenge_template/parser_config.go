package lozenge_template

type ParserConfig struct {
	TrimSpaces bool
}

func NewParserConfig() ParserConfig {
	return ParserConfig{}
}

func (pc ParserConfig) WithTrimSpaces() ParserConfig {
	pc.TrimSpaces = true
	return pc
}
