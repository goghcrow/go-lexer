package lexer

import (
	"fmt"
	"unicode/utf8"
)

// BuildLexer
// 这里没有取最长匹配, 而是首次匹配, so, 注意规则顺序
// 具体可以参考 example/lexicon.go
func BuildLexer[TokenKind comparable](f func(lexicon *Lexicon[TokenKind])) *Lexer[TokenKind] {
	lex := NewLexicon[TokenKind]()
	f(&lex)
	return NewLexer(lex)
}

func NewLexer[TokenKind comparable](lex Lexicon[TokenKind]) *Lexer[TokenKind] {
	return &Lexer[TokenKind]{Lexicon: lex}
}

func (l *Lexer[TokenKind]) MustLex(input string) []*Token[TokenKind] {
	toks, err := l.Lex(input)
	if err != nil {
		panic(err)
	}
	return toks
}

func (l *Lexer[TokenKind]) Lex(input string) ([]*Token[TokenKind], error) {
	l.input = []rune(input)
	l.Pos = Pos{}
	var toks []*Token[TokenKind]
	for {
		t, keep, err := l.next()
		if err != nil {
			return toks, err
		}
		if t == nil {
			break
		}
		if keep {
			toks = append(toks, t)
		}
	}
	return toks, nil
}

type Lexer[TokenKind comparable] struct {
	Lexicon[TokenKind]
	Pos
	input []rune
}

func (l *Lexer[TokenKind]) next() (tok *Token[TokenKind], keep bool, err error) {
	if l.Idx >= len(l.input) {
		return nil, true, nil
	}

	pos := l.Pos
	sub := string(l.input[l.Idx:])
	for _, rl := range l.Lexicon.rules {
		offset := rl.match(sub)
		if offset >= 0 {
			matched := l.input[l.Idx : l.Idx+offset]
			for _, r := range matched {
				l.Move(r)
			}
			pos.IdxEnd = l.Pos.Idx
			return &Token[TokenKind]{TokenKind: rl.TokenKind, Lexeme: string(matched), Pos: pos}, rl.keep, nil
		}
	}
	return nil, false, fmt.Errorf("syntax error in %s: nothing token matched", l.Pos)
}

func runeCount(s string) int { return utf8.RuneCountInString(s) }
