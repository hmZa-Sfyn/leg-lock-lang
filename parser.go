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
				return nil, nil, NewAsmError(token.Line, token.Column,
					"Duplicate label: "+labelName,
					"Labels must be unique – choose a different name")
			}
			labels[labelName] = len(instructions)
			i++
			continue
		}

		if token.Type != TokenInstruction {
			return nil, nil, NewAsmError(token.Line, token.Column,
				"Expected instruction, got "+token.Value,
				"Valid instructions: mov, add, sub, jmp, syscall")
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
			return nil, nil, NewAsmError(token.Line, token.Column,
				"Unknown instruction: "+token.Value,
				"Supported instructions: mov, add, sub, jmp, syscall")
		}

		i++ // consume instruction token

		var operands []Operand
		// collect operands until end of line / comment / EOF
		for tokens[i].Type != TokenComment &&
			tokens[i].Type != TokenEOF &&
			tokens[i].Line == token.Line {

			if tokens[i].Type == TokenComma {
				i++
				continue
			}

			var op Operand
			switch tokens[i].Type {
			case TokenRegister:
				op.Type = OpRegister
				op.Value = strings.ToLower(tokens[i].Value) // normalize
			case TokenImmediate:
				val, err := strconv.Atoi(tokens[i].Value)
				if err != nil {
					return nil, nil, NewAsmError(tokens[i].Line, tokens[i].Column,
						"Invalid immediate value: "+tokens[i].Value,
						"Immediates must be integers (e.g. 42, -10)")
				}
				op.Type = OpImmediate
				op.Value = val
			case TokenLabel:
				op.Type = OpLabel
				op.Value = tokens[i].Value
			default:
				return nil, nil, NewAsmError(tokens[i].Line, tokens[i].Column,
					"Unexpected token in operand position: "+tokens[i].Value,
					"Expected register (r0–r15), immediate, or label")
			}

			operands = append(operands, op)
			i++
		}

		// Basic validation per instruction type
		switch instType {
		case InstMov, InstAdd, InstSub:
			if len(operands) != 2 {
				return nil, nil, NewAsmError(token.Line, token.Column,
					fmt.Sprintf("%s expects exactly 2 operands, got %d",
						strings.ToUpper(token.Value), len(operands)),
					"Correct syntax: "+strings.ToLower(token.Value)+" r0, 42   or   "+strings.ToLower(token.Value)+" r1, r2")
			}
			if operands[0].Type != OpRegister {
				return nil, nil, NewAsmError(token.Line, token.Column,
					"Destination must be a register",
					"First operand must be a register (e.g. r0, r1, ...)")
			}
			if operands[1].Type == OpLabel {
				return nil, nil, NewAsmError(token.Line, token.Column,
					"Cannot use label as source for "+strings.ToLower(token.Value),
					"Use a register or immediate value instead")
			}

		case InstJmp:
			if len(operands) != 1 {
				return nil, nil, NewAsmError(token.Line, token.Column,
					"jmp expects exactly 1 operand (label), got "+fmt.Sprint(len(operands)),
					"Correct syntax: jmp mylabel")
			}
			if operands[0].Type != OpLabel {
				return nil, nil, NewAsmError(token.Line, token.Column,
					"jmp operand must be a label",
					"Example: jmp start  or  jmp loop1")
			}

		case InstSyscall:
			if len(operands) != 0 {
				return nil, nil, NewAsmError(token.Line, token.Column,
					"syscall takes no operands",
					"Just write: syscall")
			}
		}

		instructions = append(instructions, Instruction{
			Type:     instType,
			Operands: operands,
			Line:     token.Line,
		})
	}

	// Second pass: resolve label references to instruction indices
	for idx := range instructions {
		inst := &instructions[idx]
		for opIdx := range inst.Operands {
			op := &inst.Operands[opIdx]
			if op.Type == OpLabel {
				labelName := op.Value.(string)
				target, exists := labels[labelName]
				if !exists {
					return nil, nil, NewAsmError(inst.Line, 0,
						"Undefined label: "+labelName,
						"Make sure the label is defined with 'label:' somewhere before use")
				}
				op.Type = OpImmediate
				op.Value = target
			}
		}
	}

	return instructions, labels, nil
}
