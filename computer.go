package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
)

// VM a virtual machine that can load and run Intcode
type VM []int // VM's memory

// Size returns the current size of the memory in the VM
func (vm VM) Size() int {
	return len(vm)
}

// Write attempts to write a value to the VM's memory
func (vm VM) Write(address, val int) {
	if address < 0 {
		log.Fatalf("attempt to write to a negative address of %d", address)
	}
	if address > vm.Size() {
		log.Fatalf("attempt to write to address %d but memory stops at %d", address, vm.Size())
	}
	vm[address] = val
}

// Read attempts to retreive a value from VM's memory
func (vm VM) Read(address int) int {
	if address < 0 {
		log.Fatalf("attempt to read to a negative address of %d", address)
	}
	if address > vm.Size() {
		log.Fatalf("attempt to read from address %d but memory stops at %d", address, vm.Size())
	}
	return vm[address]
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
	vm.Write(resultAddress, vm.Read(termAddress1)+vm.Read(termAddress2))
}

// mul implements the 'mul' opcode for the VM
func (vm VM) mul(termAddress1, termAddress2, resultAddress int) {
	vm.Write(resultAddress, vm.Read(termAddress1)*vm.Read(termAddress2))
}

// Run attempts to execute the loaded Intcode program in the VM
func (vm VM) Run() error {
	if vm.Size() == 0 {
		return fmt.Errorf("no program loaded")
	}
	ip := 0 // instruction pointer
execLoop:
	for {
		opcode := vm.Read(ip)
		switch opcode {
		case 1: // addition
			vm.add(vm.Read(ip+1), vm.Read(ip+2), vm.Read(ip+3))
		case 2: // multiplication
			vm.mul(vm.Read(ip+1), vm.Read(ip+2), vm.Read(ip+3))
		case 99: // halt
			break execLoop
		default:
			return fmt.Errorf("Invalid opcode %v encountered at position %v", opcode, ip)
		}
		ip = ip + 4
		if ip > vm.Size() {
			return fmt.Errorf("no halt instruction occured before end of memory")
		}
	}
	return nil
}

// A helper function that complies with the scanner function needed by scanner.Split
func scanCommas(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, ','); i >= 0 {
		// We have a value up to a comma
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
