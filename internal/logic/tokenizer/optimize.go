package tokenizer

import (
	"fmt"

	"github.com/BestFriendChris/lozenge_template/internal/logic/token"
)

func Optimize(toks []*token.Token, trimSpaces bool) []*token.Token {
	newToks := make([]*token.Token, 0)
	curTok := func() *token.Token {
		if len(newToks) == 0 {
			return nil
		}
		return newToks[len(newToks)-1]
	}
	appendNewToks := func(tok *token.Token, joinNewline bool) {
		if cur := curTok(); cur != nil && cur.TT == tok.TT {
			if joinNewline && len(cur.S) > 0 && cur.S[len(cur.S)-1] != '\n' {
				cur.S += "\n"
			}
			if len(cur.S)+len(tok.S) > 60 {
				newToks = append(newToks, tok)
			} else {
				cur.S += tok.S
			}
		} else {
			newToks = append(newToks, tok)
		}
	}
	var alreadyTrimmedNL bool
	for i := 0; i < len(toks); i++ {
		tok := toks[i]
		switch tok.TT {
		case token.TTcontent:
			appendNewToks(tok, false)
		case token.TTws:
			if trimSpaces && isNextRealTokenCodeBlock(toks[i:]) {
				continue
			}
			tok = token.NewToken(token.TTcontent, tok.S)
			appendNewToks(tok, false)
		case token.TTnl:
			if trimSpaces && isPrevRealTokenCodeBlock(newToks) && !alreadyTrimmedNL {
				alreadyTrimmedNL = true
				continue
			}
			tok = token.NewToken(token.TTcontent, tok.S)
			appendNewToks(tok, false)
		case token.TTcodeGlobalBlock, token.TTcodeLocalBlock:
			appendNewToks(tok, true)
		default:
			newToks = append(newToks, tok)
		}
		alreadyTrimmedNL = false
	}
	return newToks
}

func isNextRealTokenCodeBlock(toks []*token.Token) bool {
	for _, tok := range toks {
		if tok.TT.IsCustom() || tok.TT == token.TTmacro || tok.TT == token.TTws {
			continue
		}
		return tok.TT == token.TTcodeGlobalBlock || tok.TT == token.TTcodeLocalBlock
	}
	return false
}

func isPrevRealTokenCodeBlock(toks []*token.Token) bool {
	if len(toks) == 0 {
		return false
	}
	for i := len(toks) - 1; i >= 0; i-- {
		tok := toks[i]
		if tok.TT.IsCustom() || tok.TT == token.TTmacro {
			fmt.Println("skipping custom or macro")
			continue
		}
		return tok.TT == token.TTcodeGlobalBlock || tok.TT == token.TTcodeLocalBlock
	}
	return false
}

func OptimizeSlc(toks []*token.TokenSlice, trimSpaces bool) []*token.TokenSlice {
	newToks := make([]*token.TokenSlice, 0)
	curTok := func() *token.TokenSlice {
		if len(newToks) == 0 {
			return nil
		}
		return newToks[len(newToks)-1]
	}
	appendNewToks := func(tok *token.TokenSlice, joinNewline bool) {
		if cur := curTok(); cur != nil && cur.TT == tok.TT {
			curSlcLen := cur.Slc.Len()
			if joinNewline && curSlcLen > 0 && cur.Slc.S[curSlcLen-1] != '\n' {
				cur.Slc.S += "\n"
			}
			if curSlcLen+tok.Slc.Len() > 60 {
				newToks = append(newToks, tok)
			} else {
				cur.Slc = cur.Slc.Join(tok.Slc)
			}
		} else {
			newToks = append(newToks, tok)
		}
	}
	var alreadyTrimmedNL bool
	for i := 0; i < len(toks); i++ {
		tok := toks[i]
		switch tok.TT {
		case token.TTcontent:
			appendNewToks(tok, false)
		case token.TTws:
			if trimSpaces && isNextRealTokenCodeBlockSlc(toks[i:]) {
				continue
			}
			tok = token.NewTokenSlice(token.TTcontent, tok.Slc)
			appendNewToks(tok, false)
		case token.TTnl:
			if trimSpaces && isPrevRealTokenCodeBlockSlc(newToks) && !alreadyTrimmedNL {
				alreadyTrimmedNL = true
				continue
			}
			tok = token.NewTokenSlice(token.TTcontent, tok.Slc)
			appendNewToks(tok, false)
		case token.TTcodeGlobalBlock, token.TTcodeLocalBlock:
			appendNewToks(tok, true)
		default:
			newToks = append(newToks, tok)
		}
		alreadyTrimmedNL = false
	}
	return newToks
}

func isNextRealTokenCodeBlockSlc(toks []*token.TokenSlice) bool {
	for _, tok := range toks {
		if tok.TT.IsCustom() || tok.TT == token.TTmacro || tok.TT == token.TTws {
			continue
		}
		return tok.TT == token.TTcodeGlobalBlock || tok.TT == token.TTcodeLocalBlock
	}
	return false
}

func isPrevRealTokenCodeBlockSlc(toks []*token.TokenSlice) bool {
	if len(toks) == 0 {
		return false
	}
	for i := len(toks) - 1; i >= 0; i-- {
		tok := toks[i]
		if tok.TT.IsCustom() || tok.TT == token.TTmacro {
			fmt.Println("skipping custom or macro")
			continue
		}
		return tok.TT == token.TTcodeGlobalBlock || tok.TT == token.TTcodeLocalBlock
	}
	return false
}
