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

// peekChar read a char from buffer into l.ch
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
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
		if l.peekChar() == '=' {
			// move ahead
			l.readChar()
			tok = token.CreateForStr(token.EQ, "==")
		} else {
			tok = token.CreateForByte(token.ASSIGN, l.ch)
		}
	case '+':
		tok = token.CreateForByte(token.PLUS, l.ch)
	case ';':
		tok = token.CreateForByte(token.SEMICOLON, l.ch)
	case '(':
		tok = token.CreateForByte(token.LPAREN, l.ch)
	case ')':
		tok = token.CreateForByte(token.RPAREN, l.ch)
	case '{':
		tok = token.CreateForByte(token.LBRACE, l.ch)
	case '}':
		tok = token.CreateForByte(token.RBRACE, l.ch)
	case ',':
		tok = token.CreateForByte(token.COMMA, l.ch)
	case '-':
		tok = token.CreateForByte(token.MINUS, l.ch)
	case '/':
		tok = token.CreateForByte(token.SLASH, l.ch)
	case '!':

		if l.peekChar() == '=' {
			// move ahead
			l.readChar()
			tok = token.CreateForStr(token.NOT_EQ, "!=")
		} else {
			tok = token.CreateForByte(token.BANG, l.ch)
		}

	case '>':
		tok = token.CreateForByte(token.GT, l.ch)

	case '<':
		tok = token.CreateForByte(token.LT, l.ch)
	case '*':
		tok = token.CreateForByte(token.ASTERISK, l.ch)

	case '"':
		// start of string literal
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok = token.Eof()

	case '[':
		tok = token.CreateForByte(token.LBRACKET, l.ch)
	case ']':
		tok = token.CreateForByte(token.RBRACKET, l.ch)

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

func (l *Lexer) readString() string {
	position := l.position + 1 // skip "
	prev := l.ch
	for {
		l.readChar()

		if (l.ch == '"' && prev != '\\') || l.ch == 0 {
			break
		}
		prev = l.ch
	}
	return l.input[position:l.position]
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
