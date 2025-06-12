package main

import (
	"unicode"
)

type TokenType string

const (
	TokenLBrace   TokenType = "LBRACE"
	TokenRBrace   TokenType = "RBRACE"
	TokenColon    TokenType = "COLON"
	TokenComma    TokenType = "COMMA"
	TokenString   TokenType = "STRING"
	TokenNumber   TokenType = "NUMBER"
	TokenTrue     TokenType = "TRUE"
	TokenFalse    TokenType = "FALSE"
	TokenNull     TokenType = "NULL"
	TokenLBracket TokenType = "LBRACKET"
	TokenRBracket TokenType = "RBRACKET"
	TokenEOF      TokenType = "EOF"
	TokenIllegal  TokenType = "ILLEGAL"
)

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	var result []rune
	l.readChar() // skip initial "

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case '"':
				result = append(result, '"')
			case '\\':
				result = append(result, '\\')
			case '/':
				result = append(result, '/')
			case 'b':
				result = append(result, '\b')
			case 'f':
				result = append(result, '\f')
			case 'n':
				result = append(result, '\n')
			case 'r':
				result = append(result, '\r')
			case 't':
				result = append(result, '\t')
			default:
				result = append(result, rune(l.ch))
			}
		} else {
			result = append(result, rune(l.ch))
		}
		l.readChar()
	}

	l.readChar() // skip closing "
	return string(result)
}

func (l *Lexer) readNumber() string {
	start := l.position

	if l.ch == '-' {
		l.readChar()
	}

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	if l.ch == 'e' || l.ch == 'E' {
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[start:l.position]
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	switch l.ch {
	case '{':
		l.readChar()
		return Token{Type: TokenLBrace, Literal: "{"}
	case '}':
		l.readChar()
		return Token{Type: TokenRBrace, Literal: "}"}
	case '[':
		l.readChar()
		return Token{Type: TokenLBracket, Literal: "["}
	case ']':
		l.readChar()
		return Token{Type: TokenRBracket, Literal: "]"}
	case ':':
		l.readChar()
		return Token{Type: TokenColon, Literal: ":"}
	case ',':
		l.readChar()
		return Token{Type: TokenComma, Literal: ","}
	case '"':
		return Token{Type: TokenString, Literal: l.readString()}
	case 't':
		if l.peekKeyword("true") {
			l.readN(4)
			return Token{Type: TokenTrue, Literal: "true"}
		}
	case 'f':
		if l.peekKeyword("false") {
			l.readN(5)
			return Token{Type: TokenFalse, Literal: "false"}
		}
	case 'n':
		if l.peekKeyword("null") {
			l.readN(4)
			return Token{Type: TokenNull, Literal: "null"}
		}
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		num := l.readNumber()
		return Token{Type: TokenNumber, Literal: num}
	case 0:
		return Token{Type: TokenEOF, Literal: ""}
	}

	return Token{Type: TokenIllegal, Literal: string(l.ch)}
}

func (l *Lexer) peekKeyword(keyword string) bool {
	return len(l.input[l.position:]) >= len(keyword) &&
		l.input[l.position:l.position+len(keyword)] == keyword
}

func (l *Lexer) readN(n int) {
	for i := 0; i < n; i++ {
		l.readChar()
	}
}
