// errors.go
package main

import (
	"fmt"
	"strings"
)

type AsmError struct {
	Line   int
	Column int
	Msg    string
	Hint   string
}

func NewAsmError(line, col int, msg, hint string) *AsmError {
	return &AsmError{Line: line, Column: col, Msg: msg, Hint: hint}
}

func (e *AsmError) Print(code string) {
	lines := strings.Split(code, "\n")
	if e.Line-1 >= len(lines) {
		fmt.Printf("Error: %s\nHint: %s\n", e.Msg, e.Hint)
		return
	}
	line := lines[e.Line-1]

	fmt.Printf("Error at line %d, column %d: %s\n", e.Line, e.Column, e.Msg)
	fmt.Println(line)
	if e.Column > 0 {
		fmt.Println(strings.Repeat(" ", e.Column-1) + "^")
	}
	fmt.Printf("Hint: %s\n", e.Hint)
}
