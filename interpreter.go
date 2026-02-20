// interpreter.go (added this for execution logic)
package main

import (
	"strconv"
)

const NumRegisters = 16

type VM struct {
	Registers [NumRegisters]int
	PC        int
}

func Interpret(instructions []Instruction, _ map[string]int) error {
	vm := VM{}

	for vm.PC < len(instructions) {
		inst := instructions[vm.PC]
		switch inst.Type {
		case InstMov:
			dest := regIndex(inst.Operands[0].Value.(string))
			src := inst.Operands[1]
			if src.Type == OpRegister {
				vm.Registers[dest] = vm.Registers[regIndex(src.Value.(string))]
			} else {
				vm.Registers[dest] = src.Value.(int)
			}
		case InstAdd:
			dest := regIndex(inst.Operands[0].Value.(string))
			src := inst.Operands[1]
			var val int
			if src.Type == OpRegister {
				val = vm.Registers[regIndex(src.Value.(string))]
			} else {
				val = src.Value.(int)
			}
			vm.Registers[dest] += val
		case InstSub:
			dest := regIndex(inst.Operands[0].Value.(string))
			src := inst.Operands[1]
			var val int
			if src.Type == OpRegister {
				val = vm.Registers[regIndex(src.Value.(string))]
			} else {
				val = src.Value.(int)
			}
			vm.Registers[dest] -= val
		case InstJmp:
			vm.PC = inst.Operands[0].Value.(int) - 1 // -1 because PC++ at end
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
