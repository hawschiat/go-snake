package utilities

import "fmt"

// Max returns the maximum of two integers
func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

// Center returns a string centered in a block
func Center(s string, w int) string {
	// https://stackoverflow.com/questions/41133006/how-to-fmt-printprint-this-on-the-center
	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(s))/2, s))
}
