// lexer.go
package main

import (
	"strings"
	"unicode"
)

type TokenType int

const (
	TokenInstruction TokenType = iota
	TokenRegister
	TokenImmediate
	TokenLabel
	TokenLabelDef
	TokenComma
	TokenComment
	TokenLBracket // [
	TokenRBracket // ]
	TokenEOF
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

func Lex(code string) ([]Token, *AsmError) {
	var tokens []Token
	lines := strings.Split(code, "\n")
	lineNum := 1

	for _, line := range lines {
		col := 1
		i := 0
		line = strings.TrimSpace(line) // Trim leading/trailing whitespace
		if len(line) == 0 {
			lineNum++
			continue
		}

		for i < len(line) {
			c := rune(line[i])

			if unicode.IsSpace(c) {
				i++
				col++
				continue
			}

			if c == ';' {
				// Comment
				tokens = append(tokens, Token{TokenComment, line[i:], lineNum, col})
				break // Rest of line is comment
			}

			if c == ':' {
				// Label definition, but check if preceded by identifier
				if len(tokens) > 0 && tokens[len(tokens)-1].Type == TokenLabel {
					prev := tokens[len(tokens)-1]
					tokens[len(tokens)-1] = Token{TokenLabelDef, prev.Value + ":", lineNum, prev.Column}
					i++
					col++
					continue
				} else {
					return nil, NewAsmError(lineNum, col, "Unexpected ':' without label name", "Label definitions should be like 'label:'")
				}
			}

			if c == ',' {
				tokens = append(tokens, Token{TokenComma, ",", lineNum, col})
				i++
				col++
				continue
			}

			if c == '[' {
				tokens = append(tokens, Token{TokenLBracket, "[", lineNum, col})
				i++
				col++
				continue
			}

			if c == ']' {
				tokens = append(tokens, Token{TokenRBracket, "]", lineNum, col})
				i++
				col++
				continue
			}

			if unicode.IsLetter(c) {
				// Instruction, register, or label
				start := i
				for i < len(line) && (unicode.IsLetter(rune(line[i])) || unicode.IsDigit(rune(line[i]))) {
					i++
				}
				value := line[start:i]
				col += i - start

				lowerValue := strings.ToLower(value)
				if lowerValue == "mov" || lowerValue == "add" || lowerValue == "sub" || lowerValue == "jmp" ||
					lowerValue == "syscall" || lowerValue == "cmp" || lowerValue == "jz" || lowerValue == "jnz" ||
					lowerValue == "je" || lowerValue == "jne" || lowerValue == "jg" || lowerValue == "jl" {
					tokens = append(tokens, Token{TokenInstruction, value, lineNum, start + 1})
				} else if strings.HasPrefix(lowerValue, "r") && len(value) > 1 && unicode.IsDigit(rune(value[1])) {
					tokens = append(tokens, Token{TokenRegister, value, lineNum, start + 1})
				} else {
					tokens = append(tokens, Token{TokenLabel, value, lineNum, start + 1})
				}
				continue
			}

			if unicode.IsDigit(c) || c == '-' {
				// Immediate
				start := i
				if c == '-' {
					i++
				}
				for i < len(line) && unicode.IsDigit(rune(line[i])) {
					i++
				}
				value := line[start:i]
				col += i - start
				tokens = append(tokens, Token{TokenImmediate, value, lineNum, start + 1})
				continue
			}

			return nil, NewAsmError(lineNum, col, "Unexpected character: "+string(c), "Check for typos in instructions or operands")
		}

		lineNum++
	}

	tokens = append(tokens, Token{TokenEOF, "", lineNum, 0})
	return tokens, nil
}
