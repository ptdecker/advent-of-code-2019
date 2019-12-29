package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

/*
 * Utility functions utilized by all problems
 */

// scanCommas is a helper function that complies with the scanner function needed by scanner.Split
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

// prompt displays a prompt on the console and then returns the string input
func prompt(query, promptstr string) string {
	reader := bufio.NewReader(os.Stdin)
	if len(query) > 0 {
		fmt.Println(query)
	}
	fmt.Printf("%s ", promptstr)
	text, _ := reader.ReadString('\n')
	return strings.Replace(text, "\n", "", -1)
}
