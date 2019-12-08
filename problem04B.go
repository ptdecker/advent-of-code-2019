/* Problem 04-B
 *
 * An Elf just remembered one more important detail: the two adjacent matching digits are not part
 * of a larger group of matching digits.
 *
 * Given this additional criterion, but still ignoring the range rule, the following are now true:
 *
 * 112233 meets these criteria because the digits never decrease and all repeated digits are exactly
 * two digits long.
 *
 * 123444 no longer meets the criteria (the repeated 44 is part of a larger group of 444).
 *
 * 111122 meets the criteria (even though 1 is repeated more than twice, it still contains
 * a double 22).
 *
 * How many different passwords within the range given in your puzzle input meet all of the criteria?
 *
 * Answer:
 */

package main

import (
	"log"
	"strconv"
	"strings"
)

func altPassesRules(guess int) bool {

	guessSplit := strings.Split(strconv.Itoa(guess), "")

	// Must be a six-digit number
	if len(guessSplit) != 6 {
		return false
	}

	// Develop histgram of repeating digits
	digitRepeats := map[int]int{}
	for index := 0; index < 6; index++ {
		val, err := strconv.Atoi(guessSplit[index])
		if err != nil {
			log.Fatalf("could not convert '%v' to an integer", val)
		}
		digitRepeats[val]++
	}

	// Determine if we have the repeating digits rule satisfied
	noDoubles := true
	for _, val := range digitRepeats {
		noDoubles = noDoubles && (val != 2)
	}
	if noDoubles {
		return false
	}

	// Going from left to right, the digits never decrease; they only ever increase or stay the same
	increases := true
	for index := 1; index < 6 && increases; index++ {
		val1, err := strconv.Atoi(guessSplit[index-1])
		if err != nil {
			log.Fatalf("could not convert '%v' to an integer", val1)
		}
		val2, err := strconv.Atoi(guessSplit[index])
		if err != nil {
			log.Fatalf("could not convert '%v' to an integer", val1)
		}
		increases = (val1 <= val2)
	}
	if !increases {
		return false
	}

	return true
}

func problem04B(start, end int) int {

	count := 0
	for guess := start; guess <= end; guess++ {
		if altPassesRules(guess) {
			count++
		}
	}

	return count
}
