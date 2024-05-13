package lang

import (
	"net"
)

type Lexer struct {
	input            string
	fallbackPosition int
	position         int
	nextPosition     int
	ch               byte
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	}
	return l.input[l.nextPosition]
}

func (l *Lexer) readChar() {
	if l.fallbackPosition != -1 {
		l.fallbackPosition = -1
	}
	l.ch = l.peekChar()
	l.position = l.nextPosition
	l.nextPosition += 1
}

func (l *Lexer) fallback() {
	if l.fallbackPosition == -1 {
		// silently ignore
		return
	}
	l.ch = l.input[l.fallbackPosition]
	l.position = l.fallbackPosition
	l.nextPosition = l.position + 1
}

func (l *Lexer) unsafeReadChar() {
	if l.fallbackPosition == -1 {
		l.fallbackPosition = l.position
	}
	l.ch = l.peekChar()
	l.position = l.nextPosition
	l.nextPosition += 1
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, fallbackPosition: -1}
	l.readChar()
	return l
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isHex(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func (l *Lexer) readIdentifierAlphaNumeric() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readCIDR() (string, error) {
	position := l.position
	for isHex(l.ch) || l.ch == ':' || l.ch == '.' || l.ch == '/' {
		l.unsafeReadChar()
	}
	_, _, err := net.ParseCIDR(l.input[position:l.position])
	if err != nil {
		l.fallback()
		return "", err
	}
	return l.input[position:l.position], nil
}

func (l *Lexer) readIPv4Address() string {
	position := l.position
	for isDigit(l.ch) || l.ch == '.' {
		l.unsafeReadChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readIPv6Address() string {
	position := l.position
	for isHex(l.ch) || l.ch == ':' {
		l.unsafeReadChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()
	switch l.ch {
	case '(':
		tok = Token{Type: LeftParen, Literal: string(l.ch)}
	case ')':
		tok = Token{Type: RightParen, Literal: string(l.ch)}
	case ',':
		tok = Token{Type: Comma, Literal: string(l.ch)}
	case '"':
		tok.Type = String
		l.readChar()
		if l.ch == '"' {
			tok.Literal = ""
		} else {
			tok.Literal = l.readIdentifierAlphaNumeric()
		}
		if l.ch == '"' {
      l.readChar()
			return tok
		} else {
			tok.Type = Illegal
			tok.Literal = string(l.ch)
		}
	case ';':
		tok = Token{Type: Semicolon, Literal: string(l.ch)}
	case 0:
		tok = Token{Type: EOF, Literal: ""}
	default:
		if l.ch == ':' {
			tok.Type = CIDR
			if literal, err := l.readCIDR(); err == nil {
				tok.Literal = literal
				return tok
			}
		} else if isLetter(l.ch) {
			tok.Type = CIDR
			if literal, err := l.readCIDR(); err == nil {
				tok.Literal = literal
				return tok
			}
			tok.Literal = l.readIdentifierAlphaNumeric()
			tok.Type = LoopupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = CIDR
			if literal, err := l.readCIDR(); err == nil {
				tok.Literal = literal
				return tok
			}
		}
		tok = Token{Type: Illegal, Literal: string(l.ch)}
	}

	l.readChar()
	return tok
}
