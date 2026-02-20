// interpreter.go
package main

import (
	"fmt"
	"strconv"
)

const (
	NumRegisters = 16
	MemorySize   = 1024 // Simple flat memory model (bytes, but we treat as ints for simplicity)
)

type Flags struct {
	Zero     bool // ZF
	Sign     bool // SF (negative)
	Carry    bool // CF
	Overflow bool // OF
}

type VM struct {
	Registers [NumRegisters]int
	Memory    [MemorySize]int // Treat as int array for simplicity (4-byte words)
	PC        int
	Flags     Flags
}

func Interpret(instructions []Instruction, _ map[string]int) error {
	vm := VM{}

	for vm.PC < len(instructions) {
		inst := instructions[vm.PC]
		switch inst.Type {
		case InstMov:
			dest := inst.Operands[0]
			src := inst.Operands[1]

			var srcVal int
			switch src.Type {
			case OpRegister:
				srcVal = vm.Registers[regIndex(src.Value.(string))]
			case OpImmediate:
				srcVal = src.Value.(int)
			case OpMemory:
				addr := getAddress(&vm, src.Value.(MemoryOperand))
				srcVal = vm.Memory[addr]
			default:
				return fmt.Errorf("invalid source operand for mov at line %d", inst.Line)
			}

			switch dest.Type {
			case OpRegister:
				vm.Registers[regIndex(dest.Value.(string))] = srcVal
			case OpMemory:
				addr := getAddress(&vm, dest.Value.(MemoryOperand))
				vm.Memory[addr] = srcVal
			default:
				return fmt.Errorf("invalid dest operand for mov at line %d", inst.Line)
			}

		case InstAdd:
			destIdx := regIndex(inst.Operands[0].Value.(string))
			src := inst.Operands[1]

			var srcVal int
			switch src.Type {
			case OpRegister:
				srcVal = vm.Registers[regIndex(src.Value.(string))]
			case OpImmediate:
				srcVal = src.Value.(int)
			case OpMemory:
				addr := getAddress(&vm, src.Value.(MemoryOperand))
				srcVal = vm.Memory[addr]
			}

			oldVal := vm.Registers[destIdx]
			newVal := oldVal + srcVal
			vm.Registers[destIdx] = newVal
			updateFlags(&vm.Flags, newVal, oldVal > 0 && srcVal > 0 && newVal < 0) // simple overflow check

		case InstSub:
			destIdx := regIndex(inst.Operands[0].Value.(string))
			src := inst.Operands[1]

			var srcVal int
			switch src.Type {
			case OpRegister:
				srcVal = vm.Registers[regIndex(src.Value.(string))]
			case OpImmediate:
				srcVal = src.Value.(int)
			case OpMemory:
				addr := getAddress(&vm, src.Value.(MemoryOperand))
				srcVal = vm.Memory[addr]
			}

			oldVal := vm.Registers[destIdx]
			newVal := oldVal - srcVal
			vm.Registers[destIdx] = newVal
			updateFlags(&vm.Flags, newVal, oldVal < 0 && srcVal > 0 && newVal > 0) // simple overflow

		case InstCmp:
			dest := inst.Operands[0]
			src := inst.Operands[1]

			var destVal, srcVal int
			switch dest.Type {
			case OpRegister:
				destVal = vm.Registers[regIndex(dest.Value.(string))]
			case OpMemory:
				addr := getAddress(&vm, dest.Value.(MemoryOperand))
				destVal = vm.Memory[addr]
			}
			switch src.Type {
			case OpRegister:
				srcVal = vm.Registers[regIndex(src.Value.(string))]
			case OpImmediate:
				srcVal = src.Value.(int)
			case OpMemory:
				addr := getAddress(&vm, src.Value.(MemoryOperand))
				srcVal = vm.Memory[addr]
			}

			result := destVal - srcVal
			updateFlags(&vm.Flags, result, false) // no overflow for cmp

		case InstJmp:
			vm.PC = inst.Operands[0].Value.(int) - 1 // adjust for PC++
			continue

		case InstJz, InstJe:
			if vm.Flags.Zero {
				vm.PC = inst.Operands[0].Value.(int) - 1
				continue
			}
		case InstJnz, InstJne:
			if !vm.Flags.Zero {
				vm.PC = inst.Operands[0].Value.(int) - 1
				continue
			}
		case InstJg:
			if !vm.Flags.Zero && !vm.Flags.Sign {
				vm.PC = inst.Operands[0].Value.(int) - 1
				continue
			}
		case InstJl:
			if vm.Flags.Sign {
				vm.PC = inst.Operands[0].Value.(int) - 1
				continue
			}

		case InstSyscall:
			err := handleSyscall(&vm)
			if err != nil {
				return err
			}
		}
		vm.PC++
	}
	return nil
}

func regIndex(reg string) int {
	idx, _ := strconv.Atoi(reg[1:])
	return idx
}

func getAddress(vm *VM, mem MemoryOperand) int {
	switch base := mem.Base.(type) {
	case string: // reg
		return vm.Registers[regIndex(base)] % MemorySize // simple modulo to bound
	case int: // imm
		return base % MemorySize
	default:
		panic("invalid memory base")
	}
}

func updateFlags(flags *Flags, result int, overflow bool) {
	flags.Zero = result == 0
	flags.Sign = result < 0
	flags.Overflow = overflow
	// Carry: for simplicity, set if result < 0 for sub, but real impl more complex
	flags.Carry = result < 0 // placeholder
}
