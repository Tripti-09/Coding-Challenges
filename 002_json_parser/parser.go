package main

import (
	"fmt"
	// "strconv"
)

type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// Entry point for JSON parsing (object only for now)
func (p *Parser) ParseObject() error {
	if p.currentToken.Type != TokenLBrace {
		return fmt.Errorf("expected { but got %s", p.currentToken.Literal)
	}
	p.nextToken()

	// Handle empty object
	if p.currentToken.Type == TokenRBrace {
		p.nextToken()
		return nil
	}

parseLoop:
	for {
		if p.currentToken.Type != TokenString {
			return fmt.Errorf("expected string key but got %s", p.currentToken.Literal)
		}
		p.nextToken()

		if p.currentToken.Type != TokenColon {
			return fmt.Errorf("expected : but got %s", p.currentToken.Literal)
		}
		p.nextToken()

		if err := p.parseValue(); err != nil {
			return err
		}

		switch p.currentToken.Type {
		case TokenComma:
			p.nextToken()
			continue parseLoop
		case TokenRBrace:
			p.nextToken()
			if p.currentToken.Type != TokenEOF {
				continue
			}
			break parseLoop
		default:
			return fmt.Errorf("expected , or } but got %s", p.currentToken.Literal)
		}
	}

	return nil
}

// Parses any JSON value
func (p *Parser) parseValue() error {
	// fmt.Println("token:", p.currentToken.Type)
	
	switch p.currentToken.Type {
	case TokenString, TokenNumber, TokenTrue, TokenFalse, TokenNull:
		p.nextToken()
		return nil
	case TokenLBrace:
		return p.ParseObject()
	case TokenLBracket:
		return p.ParseArray()
	default:
		return fmt.Errorf("unexpected value: %s", p.currentToken.Literal)
	}
}

// Parses JSON arrays (e.g., [1, "two", false])
func (p *Parser) ParseArray() error {
	if p.currentToken.Type != TokenLBracket {
		return fmt.Errorf("expected [ but got %s", p.currentToken.Literal)
	}
	p.nextToken()

	// Handle empty array
	if p.currentToken.Type == TokenRBracket {
		p.nextToken()
		return nil
	}

parseArray:
	for {
		if err := p.parseValue(); err != nil {
			return err
		}

		switch p.currentToken.Type {
		case TokenComma:
			p.nextToken()
			continue parseArray
		case TokenRBracket:
			p.nextToken()
			break parseArray
		default:
			return fmt.Errorf("expected , or ] but got %s", p.currentToken.Literal)
		}
	}
	return nil
}
