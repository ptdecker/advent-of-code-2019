/*
 * Ship's computer
 */

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// VM a virtual machine that can load and run Intcode
type VM []int // VM's memory

// ParamMode is an enum that defines opcode parameter mode
type ParamMode int

const (
	// PositionMode 'position' parameter mode
	PositionMode ParamMode = 0
	// ImmediateMode 'immediate' parameter mode
	ImmediateMode ParamMode = 1
	// ValueMode 'value' mode
	ValueMode ParamMode = 9
)

// Size returns the current size of the memory in the VM
func (vm VM) Size() int {
	return len(vm)
}

// immediateWrite attempts to write a value to the VM's memory using the 'immediate' mode
// where the passed address is the address of the desired data
func (vm VM) immediateWrite(address, val int) {
	if address < 0 {
		log.Fatalf("attempt to write to a negative address of %d", address)
	}
	if address > vm.Size() {
		log.Fatalf("attempt to write to address %d but memory stops at %d", address, vm.Size())
	}
	vm[address] = val
}

// positionWrite attempts to write a value to the VM's memory using the 'position' mode
// where the location pointed to by the passed address contains the address of the desired
// data
func (vm VM) positionWrite(address, val int) {
	vm.immediateWrite(vm[address], val)
}

// ModeWrite writes to memory using the mode passed
func (vm VM) ModeWrite(address, val int, mode ParamMode) {
	if mode == ImmediateMode {
		vm.immediateWrite(address, val)
	} else {
		vm.positionWrite(address, val)
	}
}

// immediateRead attempts to retreive a value from VM's memory using the 'immediate' mode
// where the passed address is the address of the desired data
func (vm VM) immediateRead(address int) int {
	if address < 0 {
		log.Fatalf("attempt to read to a negative address of %d", address)
	}
	if address > vm.Size() {
		log.Fatalf("attempt to read from address %d but memory stops at %d", address, vm.Size())
	}
	return vm[address]
}

// positionRead attempts to retreive a value from VM's memory using the 'position' mode
// where the location pointed to by the passed address contains the address of the desired
// data.
func (vm VM) positionRead(address int) int {
	return vm.immediateRead(vm[address])
}

// ModeRead reads from memory using the mode passed
func (vm VM) ModeRead(address int, mode ParamMode) int {
	if mode == ImmediateMode {
		return vm.immediateRead(address)
	}
	return vm.positionRead(address)
}

// Load attempts to load VM's memory with Intcode from a file
func (vm VM) Load(fileName string) (VM, error) {

	// Open data file containing a program
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot load memory from %s: %v", fileName, err)
	}
	defer file.Close()

	// Read data file containing the program and load it into core memory
	scanner := bufio.NewScanner(file)
	scanner.Split(scanCommas)
	address := 0
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("at address %d: %v", address, val)
		}
		vm = append(vm, val)
		address++
	}

	return vm, nil
}

// add implments the 'add' opcode for the VM
func (vm VM) add(termAddress1, termAddress2, resultAddress int) {
	vm.ModeWrite(resultAddress, vm.ModeRead(termAddress1, PositionMode)+vm.ModeRead(termAddress2, PositionMode), PositionMode)
}

// mul implements the 'mul' opcode for the VM
func (vm VM) mul(termAddress1, termAddress2, resultAddress int) {
	vm.ModeWrite(resultAddress, vm.ModeRead(termAddress1, PositionMode)*vm.ModeRead(termAddress2, PositionMode), PositionMode)
}

// fmtParam formats a parameter for verbose display depending upon the mode
func (vm VM) fmtParam(val int, mode ParamMode) string {
	switch mode {
	case ImmediateMode:
		return fmt.Sprintf("$%4d  (%d)", val, vm.ModeRead(val, mode))
	case PositionMode:
		return fmt.Sprintf("[%4d] (%d)", val, vm.ModeRead(val, mode))
	case ValueMode:
		return fmt.Sprintf(" %4d  (%d)", val, vm.ModeRead(val, mode))
	default:
		return "Err"
	}
}

// Run attempts to execute the loaded Intcode program in the VM
func (vm VM) Run(verbose bool) error {
	if vm.Size() == 0 {
		return fmt.Errorf("no program loaded")
	}
	ip := 0 // instruction pointer
execLoop:
	for {
		opcode := vm.immediateRead(ip)
		switch opcode {
		case 1: // addition
			if verbose {
				p1 := vm.fmtParam(ip+1, PositionMode)
				p2 := vm.fmtParam(ip+2, PositionMode)
				p3 := vm.fmtParam(ip+3, PositionMode)
				fmt.Printf("%4d:\tADD\t%s\t%s\t%s\n", ip, p1, p2, p3)
			}
			vm.add(ip+1, ip+2, ip+3)
			ip = ip + 4
		case 2: // multiplication
			if verbose {
				p1 := vm.fmtParam(ip+1, PositionMode)
				p2 := vm.fmtParam(ip+2, PositionMode)
				p3 := vm.fmtParam(ip+3, PositionMode)
				fmt.Printf("%4d:\tMUL\t%s\t%s\t%s\n", ip, p1, p2, p3)
			}
			vm.mul(ip+1, ip+2, ip+3)
			ip = ip + 4
		case 99: // halt
			if verbose {
				fmt.Printf("%4d:\tHLT\n", ip)
			}
			ip = ip + 1
			break execLoop
		default:
			return fmt.Errorf("Invalid opcode %v encountered at position %v", opcode, ip)
		}
		if ip > vm.Size() {
			return fmt.Errorf("no halt instruction occured before end of memory")
		}
	}
	return nil
}
