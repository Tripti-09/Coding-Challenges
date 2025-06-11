package main

import (
	// "strings"
	"unicode"
)

type TokenType string

const (
	TokenLBrace TokenType = "LBRACE"
	TokenRBrace TokenType = "RBRACE"
	TokenColon TokenType = "COLON"
	TokenComma   TokenType = "COMMA"
	TokenString  TokenType = "STRING"
	TokenEOF     TokenType = "EOF"
	TokenIllegal TokenType = "ILLEGAL"
	TokenTrue TokenType = "BTRUE"
	TokenFalse TokenType = "BFALSE"
	TokenNull TokenType = "NULL"
	TokenNumber TokenType = "NUMBER"
	TokenLBracket TokenType = "LBRACKET"
	TokenRBracket TokenType = "RBRACKET"
)

type Token struct {
	Type TokenType
	Literal string
}

type Lexer struct {
	input string
	position int
	readPosition int
	ch byte
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
	pos := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
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
	case ':':
		l.readChar()
		return Token{Type: TokenColon, Literal: ":"}
	case ',':
		l.readChar()
		return Token{Type: TokenComma, Literal: ","}
	case '"':
		str := l.readString()
		l.readChar()
		return Token{Type: TokenString, Literal: str}
	case 0:
		return Token{Type: TokenEOF, Literal: ""}
	case '[':
		l.readChar()
		return Token{Type: TokenLBracket, Literal: "["}
	case ']':
		l.readChar()
		return Token{Type: TokenRBracket, Literal: "]"}
	default:
	if isLetter(l.ch) {
		lit := l.readIdentifier()
		switch lit {
		case "true":
			return Token{Type: TokenTrue, Literal: "true"}
		case "false":
			return Token{Type: TokenFalse, Literal: "false"}
		case "null":
			return Token{Type: TokenNull, Literal: "null"}
		default:
			return Token{Type: TokenIllegal, Literal: lit}
		}
	} else if isDigit(l.ch) {
		num := l.readNumber()
		return Token{Type: TokenNumber, Literal: num}
	} else {
		return Token{Type: TokenIllegal, Literal: string(l.ch)}
	}

	}
}