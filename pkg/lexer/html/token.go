package html

import "github.com/dddplayer/hugoverse/pkg/lexer"

type Token struct {
	lexer.BaseToken
	Start lexer.Delim
	End   lexer.Delim
}
