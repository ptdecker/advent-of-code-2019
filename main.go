/*
 * Advent of Code - 2019
 */

package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

/* Problem 01-A
 *
 * The Elves quickly load you into a spacecraft and prepare to launch.
 *
 * At the first Go / No Go poll, every Elf is Go until the Fuel Counter-Upper. They haven't
 * determined the amount of fuel required yet.
 *
 * Fuel required to launch a given module is based on its mass. Specifically, to find the fuel
 * required for a module, take its mass, divide by three, round down, and subtract 2.
 *
 * For example:
 *
 * For a mass of 12, divide by 3 and round down to get 4, then subtract 2 to get 2.
 * For a mass of 14, dividing by 3 and rounding down still yields 4, so the fuel required is also 2.
 * For a mass of 1969, the fuel required is 654.
 * For a mass of 100756, the fuel required is 33583.
 * The Fuel Counter-Upper needs to know the total fuel requirement. To find it, individually
 * calculate the fuel needed for the mass of each module (your puzzle input), then add together
 * all the fuel values.
 *
 * What is the sum of the fuel requirements for all of the modules on your spacecraft?
 *
 * Answer:  3297866
 */

func fuelRequired(mass int64) int64 {
	fuel := int64(math.Floor(float64(mass)/3)) - 2
	if fuel < 0 {
		return 0
	}
	return fuel
}

func problem01A(fileName string) int64 {

	// Open data file containing the masses of each module
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read data file containing the masses of each module line by line accumulating the total fuel
	// required for all the modules processed
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	totalFuel := int64(0) // total fuel accumulator
	for scanner.Scan() {

		// Determine the mass from the line read from the file
		mass, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		// Calculate the fuel needed to cover the mass and add it to the
		// total fule required for the flight
		totalFuel += fuelRequired(int64(mass))
	}
	return totalFuel
}

/* Problem 01-B
 *
 * The Elves quickly load you into a spacecraft and prepare to launch.
 *
 * During the second Go / No Go poll, the Elf in charge of the Rocket Equation Double-Checker
 * stops the launch sequence. Apparently, you forgot to include additional fuel for the fuel
 * you just added.
 *
 * Fuel itself requires fuel just like a module - take its mass, divide by three, round down, and
 * subtract 2. However, that fuel also requires fuel, and that fuel requires fuel, and so on. Any
 * mass that would require negative fuel should instead be treated as if it requires zero fuel; the
 * remaining mass, if any, is instead handled by wishing really hard, which has no mass and is
 * outside the scope of this calculation.
 *
 * So, for each module mass, calculate its fuel and add it to the total. Then, treat the fuel amount
 * you just calculated as the input mass and repeat the process, continuing until a fuel requirement
 * is zero or negative. For example:
 *
 * A module of mass 14 requires 2 fuel. This fuel requires no further fuel (2 divided by 3 and
 * rounded down is 0, which would call for a negative fuel), so the total fuel required is still
 * just 2.
 *
 * At first, a module of mass 1969 requires 654 fuel. Then, this fuel requires 216 more fuel
 * (654 / 3 - 2). 216 then requires 70 more fuel, which requires 21 fuel, which requires 5 fuel,
 * which requires no further fuel. So, the total fuel required for a module of mass 1969 is
 * 654 + 216 + 70 + 21 + 5 = 966. The fuel required by a module of mass 100756 and its fuel is:
 * 33583 + 11192 + 3728 + 1240 + 411 + 135 + 43 + 12 + 2 = 50346.
 *
 * What is the sum of the fuel requirements for all of the modules on your spacecraft when also
 * taking into account the mass of the added fuel? (Calculate the fuel requirements for each module
 * separately, then add them all up at the end.)
 *
 * Answer: 4943923
 */

func problem01B(fileName string) int64 {

	// Open data file containing the masses of each module
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read data file containing the masses of each module line by line accumulating the total fuel
	// required for all the modules processed
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	totalFuel := int64(0) // total fuel accumulator
	for scanner.Scan() {

		// Determine the mass from the line read from the file
		mass, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		// Calculate the fuel needed to cover the mass and take into
		// consideration the fule needed to cover the fuel.  Add this fuel
		// to the total amount of fuel required for the trip
		massRemaining := int64(mass)      // The remaining amount of mass we need fuel to cover for this module
		incrementalFuelNeeded := int64(0) // The amount of fuel needed to cover that mass
		for {
			incrementalFuelNeeded = fuelRequired(massRemaining)
			totalFuel += incrementalFuelNeeded
			if incrementalFuelNeeded == 0 {
				break
			}
			massRemaining = incrementalFuelNeeded
		}

	}

	return totalFuel
}

func main() {

	fmt.Printf("Problem 01-A answer: %v\n", problem01A("day01.txt"))
	fmt.Printf("Problem 01-B answer: %v\n", problem01B("day01.txt"))

}
