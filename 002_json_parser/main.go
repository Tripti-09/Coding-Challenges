package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("JSON Parser starting...")

	if len(os.Args) < 2 {
		fmt.Println("Usage: jsonparser <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	lexer := NewLexer(string(data))
	parser := NewParser(lexer)

	if parser.currentToken.Type == TokenLBrace {
		err = parser.ParseObject()
	} else if parser.currentToken.Type == TokenLBracket {
		err = parser.ParseArray()
	} else {
		fmt.Println("JSON must start with an object or array")
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("Invalid JSON:", err)
		os.Exit(1)
	} else {
		fmt.Println("Valid JSON")
	}
}
