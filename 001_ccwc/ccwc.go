package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	// "strings"
	"unicode/utf8"
)

func countBytes(filename string) (int, error) {
	// data, err := ioutil.ReadFile(filename)
	data, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	// return len(data), nil
	return int(data.Size()), nil
}

func countLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count, scanner.Err()
}

func countWords(filename string) (int, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	words := bytes.Fields(data) // splits by all whitespace
	return len(words), nil
}

func countChars(filename string) (int, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	return utf8.RuneCount(data), nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ccwc -<option> <filename>")
		return
	}

	option := os.Args[1]
	filename := os.Args[2]

	switch option {

	case "-c":
		bytes, err := countBytes(filename)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("%8d %s\n", bytes, filename)

	case "-l":
		lines, err := countLines(filename)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("%8d %s\n", lines, filename)

	case "-w":
		words, err := countWords(filename)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("%8d %s\n", words, filename)

	case "-m":
		chars, err := countChars(filename)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("%8d %s\n", chars, filename)

	default:
		fmt.Println("Unsupported option:", option)
	}
}
