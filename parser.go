// parser.go
package main

import (
	"fmt"
	"strconv"
	"strings"
)

type InstructionType int

const (
	InstMov InstructionType = iota
	InstAdd
	InstSub
	InstJmp
	InstSyscall
)

type OperandType int

const (
	OpRegister OperandType = iota
	OpImmediate
	OpLabel
)

type Operand struct {
	Type  OperandType
	Value interface{} // string for reg/label, int for immediate
}

type Instruction struct {
	Type     InstructionType
	Operands []Operand
	Line     int
}

func Parse(tokens []Token) ([]Instruction, map[string]int, *AsmError) {
	var instructions []Instruction
	labels := make(map[string]int)
	i := 0

	for tokens[i].Type != TokenEOF {
		token := tokens[i]

		if token.Type == TokenComment {
			i++
			continue
		}

		if token.Type == TokenLabelDef {
			labelName := strings.TrimSuffix(token.Value, ":")
			if _, exists := labels[labelName]; exists {
				return nil, NewAsmError(token.Line, token.Column, "Duplicate label: "+labelName, "Labels must be unique")
			}
			labels[labelName] = len(instructions)
			i++
			continue
		}

		if token.Type != TokenInstruction {
			return nil, NewAsmError(token.Line, token.Column, "Expected instruction, got "+token.Value, "Instructions like mov, add, etc.")
		}

		var instType InstructionType
		switch strings.ToLower(token.Value) {
		case "mov":
			instType = InstMov
		case "add":
			instType = InstAdd
		case "sub":
			instType = InstSub
		case "jmp":
			instType = InstJmp
		case "syscall":
			instType = InstSyscall
		default:
			return nil, NewAsmError(token.Line, token.Column, "Unknown instruction: "+token.Value, "Supported: mov, add, sub, jmp, syscall")
		}

		i++ // Move past instruction

		var operands []Operand
		for tokens[i].Type != TokenComment && tokens[i].Type != TokenEOF && tokens[i].Line == token.Line {
			if tokens[i].Type == TokenComma {
				i++
				continue
			}

			var op Operand
			switch tokens[i].Type {
			case TokenRegister:
				op.Type = OpRegister
				op.Value = tokens[i].Value
			case TokenImmediate:
				val, err := strconv.Atoi(tokens[i].Value)
				if err != nil {
					return nil, NewAsmError(tokens[i].Line, tokens[i].Column, "Invalid immediate: "+tokens[i].Value, "Immediates should be integers")
				}
				op.Type = OpImmediate
				op.Value = val
			case TokenLabel:
				op.Type = OpLabel
				op.Value = tokens[i].Value
			default:
				return nil, NewAsmError(tokens[i].Line, tokens[i].Column, "Unexpected token: "+tokens[i].Value, "Expected register, immediate, or label")
			}

			operands = append(operands, op)
			i++
		}

		// Validate operand count
		switch instType {
		case InstMov, InstAdd, InstSub:
			if len(operands) != 2 {
				return nil, NewAsmError(token.Line, token.Column, fmt.Sprintf("Expected 2 operands, got %d", len(operands)), "E.g., mov r0, 5")
			}
			if operands[0].Type != OpRegister {
				return nil, NewAsmError(token.Line, token.Column, "Destination must be register", "First operand should be like r0")
			}
			if operands[1].Type == OpLabel {
				return nil, NewAsmError(token.Line, token.Column, "Source cannot be label for this instruction", "Use register or immediate")
			}
		case InstJmp:
			if len(operands) != 1 {
				return nil, NewAsmError(token.Line, token.Column, fmt.Sprintf("Expected 1 operand, got %d", len(operands)), "E.g., jmp label")
			}
			if operands[0].Type != OpLabel {
				return nil, NewAsmError(token.Line, token.Column, "Operand must be label", "E.g., jmp label")
			}
		case InstSyscall:
			if len(operands) != 0 {
				return nil, NewAsmError(token.Line, token.Column, "Syscall takes no operands", "Just 'syscall'")
			}
		}

		instructions = append(instructions, Instruction{instType, operands, token.Line})
	}

	// Resolve labels in instructions
	for idx, inst := range instructions {
		for opIdx, op := range inst.Operands {
			if op.Type == OpLabel {
				labelName := op.Value.(string)
				target, exists := labels[labelName]
				if !exists {
					return nil, NewAsmError(inst.Line, 0, "Undefined label: "+labelName, "Define the label with 'label:'")
				}
				instructions[idx].Operands[opIdx].Type = OpImmediate
				instructions[idx].Operands[opIdx].Value = target
			}
		}
	}

	return instructions, labels, nil
}
