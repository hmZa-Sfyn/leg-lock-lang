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
	InstCmp
	InstJz
	InstJnz
	InstJe
	InstJne
	InstJg
	InstJl
)

type OperandType int

const (
	OpRegister OperandType = iota
	OpImmediate
	OpLabel
	OpMemory // [reg] or [imm]
)

type Operand struct {
	Type  OperandType
	Value interface{} // string for reg/label, int for immediate, struct for memory {Base: string or int}
}

type MemoryOperand struct {
	Base interface{} // string (reg) or int (imm)
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
				"Valid instructions: mov, add, sub, jmp, syscall, cmp, jz, jnz, je, jne, jg, jl")
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
		case "cmp":
			instType = InstCmp
		case "jz":
			instType = InstJz
		case "jnz":
			instType = InstJnz
		case "je":
			instType = InstJe
		case "jne":
			instType = InstJne
		case "jg":
			instType = InstJg
		case "jl":
			instType = InstJl
		default:
			return nil, nil, NewAsmError(token.Line, token.Column,
				"Unknown instruction: "+token.Value,
				"Supported instructions: mov, add, sub, jmp, syscall, cmp, jz, jnz, je, jne, jg, jl")
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

			// Check for memory operand [ ... ]
			if tokens[i].Type == TokenLBracket {
				i++ // consume [
				if tokens[i].Type == TokenRegister {
					op.Type = OpMemory
					op.Value = MemoryOperand{Base: strings.ToLower(tokens[i].Value)}
					i++
				} else if tokens[i].Type == TokenImmediate {
					val, err := strconv.Atoi(tokens[i].Value)
					if err != nil {
						return nil, nil, NewAsmError(tokens[i].Line, tokens[i].Column,
							"Invalid memory address: "+tokens[i].Value,
							"Memory addresses must be integers or registers")
					}
					op.Type = OpMemory
					op.Value = MemoryOperand{Base: val}
					i++
				} else {
					return nil, nil, NewAsmError(tokens[i].Line, tokens[i].Column,
						"Invalid memory operand inside []",
						"Use [r0] or [123]")
				}
				if tokens[i].Type != TokenRBracket {
					return nil, nil, NewAsmError(tokens[i].Line, tokens[i].Column,
						"Missing closing ] for memory operand",
						"Memory access: [reg] or [imm]")
				}
				i++ // consume ]
			} else {
				switch tokens[i].Type {
				case TokenRegister:
					op.Type = OpRegister
					op.Value = strings.ToLower(tokens[i].Value) // normalize
					i++
				case TokenImmediate:
					val, err := strconv.Atoi(tokens[i].Value)
					if err != nil {
						return nil, nil, NewAsmError(tokens[i].Line, tokens[i].Column,
							"Invalid immediate value: "+tokens[i].Value,
							"Immediates must be integers (e.g. 42, -10)")
					}
					op.Type = OpImmediate
					op.Value = val
					i++
				case TokenLabel:
					op.Type = OpLabel
					op.Value = tokens[i].Value
					i++
				default:
					return nil, nil, NewAsmError(tokens[i].Line, tokens[i].Column,
						"Unexpected token in operand position: "+tokens[i].Value,
						"Expected register (r0–r15), immediate, label, or [memory]")
				}
			}

			operands = append(operands, op)
		}

		// Basic validation per instruction type
		switch instType {
		case InstMov, InstAdd, InstSub, InstCmp:
			if len(operands) != 2 {
				return nil, nil, NewAsmError(token.Line, token.Column,
					fmt.Sprintf("%s expects exactly 2 operands, got %d",
						strings.ToUpper(token.Value), len(operands)),
					"Correct syntax: "+strings.ToLower(token.Value)+" r0, 42   or   "+strings.ToLower(token.Value)+" r1, r2   or   mov r0, [123]")
			}
			// For mov: allow reg <-> mem, but not mem <-> mem
			if instType == InstMov {
				if operands[0].Type == OpMemory && operands[1].Type == OpMemory {
					return nil, nil, NewAsmError(token.Line, token.Column,
						"Cannot mov memory to memory directly",
						"Use a register as temporary: mov r0, [src]; mov [dst], r0")
				}
			}
			// For add/sub/cmp: dest must be reg, src can be reg/imm/mem
			if instType != InstMov && operands[0].Type != OpRegister {
				return nil, nil, NewAsmError(token.Line, token.Column,
					"Destination for "+strings.ToUpper(token.Value)+" must be a register",
					"First operand must be a register (e.g. r0, r1, ...)")
			}

		case InstJmp, InstJz, InstJnz, InstJe, InstJne, InstJg, InstJl:
			if len(operands) != 1 {
				return nil, nil, NewAsmError(token.Line, token.Column,
					strings.ToUpper(token.Value)+" expects exactly 1 operand (label), got "+fmt.Sprint(len(operands)),
					"Correct syntax: "+strings.ToLower(token.Value)+" mylabel")
			}
			if operands[0].Type != OpLabel {
				return nil, nil, NewAsmError(token.Line, token.Column,
					strings.ToUpper(token.Value)+" operand must be a label",
					"Example: "+strings.ToLower(token.Value)+" start  or  "+strings.ToLower(token.Value)+" loop1")
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
						"Make sure the label is defined with 'label:' somewhere")
				}
				op.Type = OpImmediate
				op.Value = target
			}
		}
	}

	return instructions, labels, nil
}
