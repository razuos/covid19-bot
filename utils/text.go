package utils

import (
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

// RemoveSpecialChars removes special characters from a string.
func RemoveSpecialChars(str string) string {
	transformer := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

	result := strings.ToLower(str)
	result, _, _ = transform.String(transformer, result)

	return result
}

// Reverse reverses a string
func Reverse(input string) string {
	n := 0
	rune := make([]rune, len(input))
	for _, r := range input {
		rune[n] = r
		n++
	}
	rune = rune[0:n]
	for i := 0; i < n/2; i++ {
		rune[i], rune[n-1-i] = rune[n-1-i], rune[i]
	}
	return string(rune)
}
