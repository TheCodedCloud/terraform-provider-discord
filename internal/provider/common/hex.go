package common

import "fmt"

// StrHex converts an integer to a hex string
func StrHex(i int) string {
	return fmt.Sprintf("#%06X", i)
}

// IntHex converts a hex string to an integer
func IntHex(s string) int {
	var i int
	fmt.Sscanf(s, "#%06X", &i)
	return i
}
