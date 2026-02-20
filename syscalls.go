// syscalls.go
package main

import (
	"fmt"
	"os"
)

func handleSyscall(vm *VM) error {
	syscallNum := vm.Registers[0] // Assume r0 holds syscall number
	switch syscallNum {
	case 0: // Exit
		os.Exit(vm.Registers[1])
	case 1: // Print int
		fmt.Println(vm.Registers[1])
	default:
		return fmt.Errorf("unknown syscall: %d", syscallNum)
	}
	return nil
}
