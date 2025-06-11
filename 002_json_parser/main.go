package main

import (
	"fmt" // for printing output
	"os" // for file operations and exiting with codes
	// "strings" // for trimming spaces from file content
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
		fmt.Println("Error reading file: ", err)
		os.Exit(1)
	}

	// content := strings.TrimSpace(string(data))

	// if content == "{}" {
	// 	fmt.Println("Valid JSON: empty object")
	// 	os.Exit(0)
	// } else {
	// 	fmt.Println("Invalid JSON")
	// 	os.Exit(1)
	// }

	lexer := NewLexer(string(data))
	parser := NewParser(lexer)

	err = parser.ParseObject()
	if err != nil {
		fmt.Println("Invalid JSON, error while parsing:", err)
		os.Exit(1)
	} else {
		fmt.Println("Valid JSON")
		os.Exit(0)
	}

}