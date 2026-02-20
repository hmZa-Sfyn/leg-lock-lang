// main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: asm-interpreter <file.asm>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	code, err := readFile(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tokens, lexErr := Lex(code)
	if lexErr != nil {
		lexErr.Print(code)
		os.Exit(1)
	}

	instructions, labels, parseErr := Parse(tokens)
	if parseErr != nil {
		parseErr.Print(code)
		os.Exit(1)
	}

	err = Interpret(instructions, labels)
	if err != nil {
		fmt.Println("Runtime error:", err)
		os.Exit(1)
	}
}

func readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var sb strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sb.WriteString(scanner.Text() + "\n")
	}
	return sb.String(), scanner.Err()
}
