/* Problem 04-A
 *
 * You arrive at the Venus fuel depot only to discover it's protected by a password. The Elves had
 * written the password on a sticky note, but someone threw it out.
 *
 * However, they do remember a few key facts about the password:
 *
 * It is a six-digit number.
 * The value is within the range given in your puzzle input.
 * Two adjacent digits are the same (like 22 in 122345).
 * Going from left to right, the digits never decrease; they only ever increase or stay the same (like 111123 or 135679).
 * Other than the range rule, the following are true:
 *
 * 111111 meets these criteria (double 11, never decreases).
 * 223450 does not meet these criteria (decreasing pair of digits 50).
 * 123789 does not meet these criteria (no double).
 *
 * How many different passwords within the range given in your puzzle input meet these criteria?
 */

package main

import (
	"log"
	"strconv"
	"strings"
)

func passesRules(guess int) bool {

	guessSplit := strings.Split(strconv.Itoa(guess), "")

	// Must be a six-digit number
	if len(guessSplit) != 6 {
		return false
	}

	// Two adjacent digits are the same
	noRepeats := true
	for index := 1; index < 6 && noRepeats; index++ {
		noRepeats = (guessSplit[index-1] != guessSplit[index])
	}
	if noRepeats {
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

func problem04A(start, end int) int {

	count := 0
	for guess := start; guess < end; guess++ {
		if passesRules(guess) {
			count++
		}
	}

	return count
}
