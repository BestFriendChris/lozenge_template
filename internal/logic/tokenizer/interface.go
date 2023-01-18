package tokenizer

type Tokenizer interface {
	NextToken(s string) (Token, string)
}
