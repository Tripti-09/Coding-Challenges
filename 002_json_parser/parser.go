package main

import (
	"fmt"
	"strconv"
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

func (p *Parser) ParseObject() error {
	if p.currentToken.Type != TokenLBrace {
		return fmt.Errorf("expected { but got %s", p.currentToken.Literal)
	}
	p.nextToken()

	if p.currentToken.Type == TokenRBrace {
		p.nextToken()
		return nil
	}

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

		if p.currentToken.Type == TokenComma {
			p.nextToken()
			continue
		} else if p.currentToken.Type == TokenRBrace {
			p.nextToken()
			break
		} else {
			return fmt.Errorf("expected , or } but got %s", p.currentToken.Literal)
		}
	}

	return nil
}

func (p *Parser) ParseArray() error {
	if p.currentToken.Type != TokenLBracket {
		return fmt.Errorf("expected [ but got %s", p.currentToken.Literal)
	}
	p.nextToken()

	if p.currentToken.Type == TokenRBracket {
		p.nextToken()
		return nil
	}

	for {
		if err := p.parseValue(); err != nil {
			return err
		}

		if p.currentToken.Type == TokenComma {
			p.nextToken()
			continue
		} else if p.currentToken.Type == TokenRBracket {
			p.nextToken()
			break
		} else {
			return fmt.Errorf("expected , or ] but got %s", p.currentToken.Literal)
		}
	}

	return nil
}

func (p *Parser) parseValue() error {
	switch p.currentToken.Type {
	case TokenString:
		p.nextToken()
		return nil
	case TokenNumber:
		if _, err := strconv.ParseFloat(p.currentToken.Literal, 64); err != nil {
			return fmt.Errorf("invalid number: %s", p.currentToken.Literal)
		}
		p.nextToken()
		return nil
	case TokenTrue, TokenFalse, TokenNull:
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
