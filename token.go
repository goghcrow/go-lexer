package lexer

import "fmt"

type Token[Kind comparable] struct {
	TokenKind Kind
	Pos
	Lexeme string
}

func (t *Token[Kind]) String() string {
	// return t.Lexeme
	return fmt.Sprintf("<'%s', %v>", t.Lexeme, t.TokenKind)
}
