package lexer

import "github.com/NishanthSpShetty/monkey/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()

	return l
}

// readChar read a char from buffer into l.ch
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpaces()
	switch l.ch {
	case '=':
		tok = token.Create(token.ASSIGN, l.ch)
	case '+':
		tok = token.Create(token.PLUS, l.ch)
	case ';':
		tok = token.Create(token.SEMICOLON, l.ch)
	case '(':
		tok = token.Create(token.LPAREN, l.ch)
	case ')':
		tok = token.Create(token.RPAREN, l.ch)
	case '{':
		tok = token.Create(token.LBRACE, l.ch)
	case '}':
		tok = token.Create(token.RBRACE, l.ch)
	case ',':
		tok = token.Create(token.COMMA, l.ch)
	case 0:
		tok = token.Eof()

	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = token.Ill()
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) readNumber() string {
	position := l.position
	// read until we see letter
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	// read until we see letter
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhiteSpaces() {
	// loop till all whitspaces are consumed
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ('0' <= ch && ch <= '9')
}
