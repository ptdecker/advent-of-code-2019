/* Problem 02-A
 *
 * An Intcode program is a list of integers separated by commas (like 1,0,0,3,99). To run one,
 * start by looking at the first integer (called position 0). Here, you will find an opcode -
 * either 1, 2, or 99. The opcode indicates what to do; for example, 99 means that the program is
 * finished and should immediately halt. Encountering an unknown opcode means something went wrong.
 *
 * Opcode 1 adds together numbers read from two positions and stores the result in a third position.
 * The three integers immediately after the opcode tell you these three positions - the first two
 * indicate the positions from which you should read the input values, and the third indicates the
 * position at which the output should be stored.
 *
 * For example, if your Intcode computer encounters 1,10,20,30, it should read the values at
 * positions 10 and 20, add those values, and then overwrite the value at position 30 with their sum.
 *
 * Opcode 2 works exactly like opcode 1, except it multiplies the two inputs instead of adding them.
 * Again, the three integers after the opcode indicate where the inputs and outputs are, not their
 * values.
 *
 * Once you're done processing an opcode, move to the next one by stepping forward 4 positions.
 *
 * For example, suppose you have the following program:
 *
 * 1,9,10,3,2,3,11,0,99,30,40,50
 *
 * For the purposes of illustration, here is the same program split into multiple lines:
 *
 * 1,9,10,3,
 * 2,3,11,0,
 * 99,
 * 30,40,50
 *
 * The first four integers, 1,9,10,3, are at positions 0, 1, 2, and 3. Together, they represent the
 * first opcode (1, addition), the positions of the two inputs (9 and 10), and the position of the
 * output (3). To handle this opcode, you first need to get the values at the input positions:
 * position 9 contains 30, and position 10 contains 40. Add these numbers together to get 70. Then,
 * store this value at the output position; here, the output position (3) is at position 3, so it
 * overwrites itself. Afterward, the program looks like this:
 *
 * 1,9,10,70,
 * 2,3,11,0,
 * 99,
 * 30,40,50
 *
 * Step forward 4 positions to reach the next opcode, 2. This opcode works just like the previous,
 * but it multiplies instead of adding. The inputs are at positions 3 and 11; these positions contain
 * 70 and 50 respectively. Multiplying these produces 3500; this is stored at position 0:
 *
 * 3500,9,10,70,
 * 2,3,11,0,
 * 99,
 * 30,40,50
 *
 * Stepping forward 4 more positions arrives at opcode 99, halting the program.
 *
 * Here are the initial and final states of a few more small programs:
 *
 * 1,0,0,0,99 becomes 2,0,0,0,99 (1 + 1 = 2).
 * 2,3,0,3,99 becomes 2,3,0,6,99 (3 * 2 = 6).
 * 2,4,4,5,99,0 becomes 2,4,4,5,99,9801 (99 * 99 = 9801).
 * 1,1,1,4,99,5,6,0,99 becomes 30,1,1,4,2,5,6,0,99.
 *
 * Once you have a working computer, the first step is to restore the gravity assist program (your
 * puzzle input) to the "1202 program alarm" state it had just before the last computer caught fire.
 * To do this, before running the program, replace position 1 with the value 12 and replace position
 * 2 with the value 2. What value is left at position 0 after the program halts?
 *
 * Answer: 4690667
 */

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

func problem02A(fileName string) int {

	vm, err := new(VM).Load(fileName)
	if err != nil {
		log.Fatal(err)
	}

	vm.Write(1, 12)
	vm.Write(2, 2)

	err = vm.Run()
	if err != nil {
		log.Fatal(err)
	}

	return vm.Read(0)
}
