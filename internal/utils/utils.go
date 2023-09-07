package utils // Package 'utils' provides a set of utility functions

import (
	"strings" // Importing the 'strings' standard library for string manipulations
)

// Function 'CountLines' takes in a string 'text' and returns the number of lines in that text.
// Achieved by splitting the text at each newline character '\n' and returning the length of the resulting slice.
func CountLines(text string) int {
	return len(strings.Split(text, "\n"))
}

// Function 'Min' takes in two integers ('a' and 'b') and returns the smaller one.
func Min(a, b int) int {
	if a < b {
		return a // If 'a' is less than 'b', return 'a'
	}
	return b // Else, return 'b'
}

// Function 'Max' takes in two integers ('a' and 'b') and returns the larger one.
func Max(a, b int) int {
	if a > b {
		return a // If 'a' is greater than 'b', return 'a'
	}
	return b // Else, return 'b'
}
