package lexer

import (
	"regexp"
	"strings"
)

const NotMatched = -1

type Rule[TokenKind comparable] struct {
	keep      bool
	TokenKind TokenKind
	match     func(string) int // 匹配返回 EndRuneCount , 失败返回 NotMatched
}

func (r *Rule[TokenKind]) Skip() *Rule[TokenKind] { r.keep = false; return r }

// Lexicon Lexical grammar
type Lexicon[TokenKind comparable] struct {
	rules []*Rule[TokenKind]
}

func NewLexicon[Kind comparable]() Lexicon[Kind] {
	return Lexicon[Kind]{}
}

func (l *Lexicon[TokenKind]) Rule(r Rule[TokenKind]) *Rule[TokenKind] {
	l.rules = append(l.rules, &r)
	return &r
}
func (l *Lexicon[TokenKind]) Str(k TokenKind, s string) *Rule[TokenKind] { return l.Rule(str(k, s)) }
func (l *Lexicon[TokenKind]) Keyword(k TokenKind, s string) *Rule[TokenKind] {
	return l.Rule(keyword(k, s))
}
func (l *Lexicon[TokenKind]) Regex(k TokenKind, pattern string) *Rule[TokenKind] {
	return l.Rule(regex(k, pattern))
}
func (l *Lexicon[TokenKind]) PrimOper(k TokenKind, oper string) *Rule[TokenKind] {
	return l.Rule(primOper(k, oper))
}
func (l *Lexicon[TokenKind]) Oper(k TokenKind, oper string) *Rule[TokenKind] {
	if IsIdentOp(oper) {
		return l.Keyword(k, oper)
	} else {
		return l.Str(k, oper)
	}
}

func str[TokenKind comparable](k TokenKind, str string) Rule[TokenKind] {
	return Rule[TokenKind]{true, k, func(s string) int {
		if strings.HasPrefix(s, str) {
			return runeCount(str)
		} else {
			return NotMatched
		}
	}}
}

var keywordPostfix = regexp.MustCompile(`^[a-zA-Z\d\p{L}_]+`)

func keyword[TokenKind comparable](k TokenKind, kw string) Rule[TokenKind] {
	return Rule[TokenKind]{true, k, func(s string) int {
		// golang regexp 不支持 lookahead
		completedWord := strings.HasPrefix(s, kw) &&
			!keywordPostfix.MatchString(s[len(kw):])
		if completedWord {
			return runeCount(kw)
		} else {
			return NotMatched
		}
	}}
}

func regex[TokenKind comparable](k TokenKind, pattern string) Rule[TokenKind] {
	startWith := regexp.MustCompile("^(?:" + pattern + ")")
	return Rule[TokenKind]{true, k, func(s string) int {
		found := startWith.FindString(s)
		if found == "" {
			return NotMatched
		} else {
			return runeCount(found)
		}
	}}
}

// primOper . ? 内置操作符的优先级高于自定义操作符, 且不是匹配最长, 需要特殊处理
// e.g 比如自定义操作符 .^. 不能匹配成 [`.`, `^.`]
func primOper[TokenKind comparable](k TokenKind, oper string) Rule[TokenKind] {
	return Rule[TokenKind]{true, k, func(s string) int {
		if !strings.HasPrefix(s, oper) {
			return NotMatched
		}
		completedOper := len(s) == len(oper) || !HasOperPrefix(s[len(oper):])
		if completedOper {
			return runeCount(oper)
		} else {
			return NotMatched
		}
	}}
}
