package parser

type ParserConfig struct {
	TrimSpaces bool
}

func NewParserConfig() ParserConfig {
	return ParserConfig{}
}

func (pc ParserConfig) WithTrimSpaces(trimSpaces bool) ParserConfig {
	newPc := pc
	newPc.TrimSpaces = trimSpaces
	return newPc
}
