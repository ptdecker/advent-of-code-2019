/*
 * Advent of Code - 2019
 */

package main

import (
	"fmt"
	"strconv"
	"strings"
)

func tryProblem(name string, expected int, actual int) {
	if expected == actual {
		fmt.Printf("Problem %s answer: %v (pass)\n", name, actual)
	} else {
		fmt.Printf("%s answer: %v (fail, expected: %v)\n", name, actual, expected)
	}
}

func runProblems() {
	tryProblem("01-A", problem01A("./data/day01.txt"), 3297866)
	tryProblem("01-B", problem01B("./data/day01.txt"), 4943923)
	tryProblem("02-A", problem02A("./data/day02.txt"), 4690667)
	tryProblem("02-B", problem02B("./data/day02.txt", 19690720), 6255)
	tryProblem("03-A", problem03A("./data/day03.txt"), 227)
	tryProblem("03-B", problem03B("./data/day03.txt"), 20286)
	tryProblem("04-A", problem04A(171309, 643603), 1625)
	tryProblem("04-B", problem04B(171309, 643603), 1111)
}

func loadConsole() {

	var vm VM
	var err error

consoleloop:
	for {
		val := prompt("", "$")
		tokens := strings.Split(val, " ")
		if len(tokens) < 1 || len(tokens[0]) < 2 {
			continue
		}
		switch strings.ToUpper(tokens[0])[:2] {
		case "LO":
			if len(tokens) < 2 {
				fmt.Println("Please provide the name of a file to load")
				break
			}
			vm, err = vm.Load(tokens[1])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			fmt.Printf("%s loaded\n", tokens[1])
		case "WR":
			if len(tokens) < 3 {
				fmt.Println("Please provide an address and data")
				break
			}
			if vm == nil {
				fmt.Println("Please load the VM first using 'LOAD <file name>'")
				break
			}
			addr, err := strconv.Atoi(tokens[1])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			val, err := strconv.Atoi(tokens[2])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			vm.ModeWrite(addr, val, ImmediateMode)
			fmt.Printf("%d written to %d\n", val, addr)
		case "RE":
			if len(tokens) < 2 {
				fmt.Println("Please provide an address to read")
				break
			}
			if vm == nil {
				fmt.Println("Please load the VM first using 'LOAD <file name>'")
				break
			}
			addr, err := strconv.Atoi(tokens[1])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			val := vm.ModeRead(addr, ImmediateMode)
			fmt.Printf("%d contains %d\n", addr, val)
		case "RU":
			if vm == nil {
				fmt.Println("Please load the VM first using 'LOAD <file name>'")
				break
			}
			err := vm.Run(false)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case "HE":
			fmt.Println("Command options:")
			fmt.Println("\tLOAD <file name>")
			fmt.Println("\tWRITE <address> <value>")
			fmt.Println("\tREAD <address>")
			fmt.Println("\tRUN")
			fmt.Println("\tQUIT")
		case "QU":
			break consoleloop
		default:
			fmt.Printf("Unrecognized command '%s'--try 'HELP'\n", val)
		}
	}
}

func main() {

mainloop:
	for {
		val := prompt("Select an option--(R)un Problems, Intcode (C)onsole, (Q)uit:", ">")
		switch strings.ToUpper(val) {
		case "R":
			runProblems()
		case "C":
			loadConsole()
		case "Q":
			break mainloop
		default:
			fmt.Println("Please enter either an 'R', 'C', or 'Q'")
		}
	}
}
