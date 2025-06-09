package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func countBytes(filename string) (int, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ccwc -c <filename>")
		return
	}

	option := os.Args[1]
	filename := os.Args[2]

	if option == "-c" {
		bytes, err := countBytes(filename)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Printf("%8d %s\n", bytes, filename)
	}
}