package lexer

import "fmt"

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) Print() {
	fmt.Printf("[type: %s, literal: %s]\n", t.Type, t.Literal)
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func newTokenFromLiteral(tokenType TokenType, literal string) Token {
	return Token{tokenType, literal}
}

const (
	Illegal = "Illegal"
	Eof     = "Eof"

	Identifier = "Identifier"

	Assign     = "="
	ParamStart = "?"

	LeftParenthesis  = "("
	RightParenthesis = ")"
	Comma            = ","
	Slash            = "/"

	Equals         = "equals"
	LessThan       = "lessThan"
	LessOrEqual    = "lessOrEqual"
	GreaterThan    = "greaterThan"
	GreaterOrEqual = "greaterOrEqual"
	Contains       = "contains"
	StartWith      = "startsWith"
	EndsWith       = "endsWith"
	Any            = "any"
	Has            = "has"
	Not            = "not"
	Or             = "or"
	And            = "and"
	Filter         = "filter"
	Null           = "null"
)

var keywords = map[string]TokenType{
	"filter":         Filter,
	"equals":         Equals,
	"lessThan":       LessThan,
	"lessOrEqual":    LessOrEqual,
	"greaterThan":    GreaterThan,
	"greaterOrEqual": GreaterOrEqual,
	"contains":       Contains,
	"startsWith":     StartWith,
	"endsWith":       EndsWith,
	"any":            Any,
	"has":            Has,
	"not":            Not,
	"or":             Or,
	"and":            And,
	"null":           Null,
}

func lookUpIdentifier(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Identifier
}

type Lexer struct {
	input        string // input for lexer
	position     int    //current position in the input
	readPosition int    // current reading position in the input
	ch           byte   //current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() Token {
	l.skipWhiteSpace()

	var tok Token

	switch l.ch {
	case '?':
		tok = newToken(ParamStart, l.ch)
	case '=':
		tok = newToken(Equals, l.ch)
	case '(':
		tok = newToken(LeftParenthesis, l.ch)
	case ')':
		tok = newToken(RightParenthesis, l.ch)
	case '\000': // eof
		tok = newToken(Eof, ' ')
	case '\'':
		quoted := l.readQuoted()
		tok = newTokenFromLiteral(Identifier, quoted)
	case ',':
		tok = newToken(Comma, l.ch)
	default:
		if isLetter(l.ch) { // either identifier or keyword - /users?filter=equals(displayName,'Brian O''Connor')
			tok.Literal = l.readIdentifier()

			if lookUpIdentifier(tok.Literal) == Identifier {
				tok.Type = Identifier
			} else {
				tok.Type = "Keyword"
			}

			return tok
		} else {
			tok = newToken(Illegal, l.ch)
			return tok
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readQuoted() string {
	position := l.position
	for l.peekChar() != '\'' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) peekChar2() byte {
	if l.readPosition+1 >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition+1]
}

func (l *Lexer) skipWhiteSpace() {

	for l.ch == '%' && l.peekChar() == '2' && l.peekChar2() == '0' {
		l.readChar()
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}
