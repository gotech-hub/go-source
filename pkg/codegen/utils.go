package codegen

import (
	"crypto/rand"
	"math"
	"math/big"
	"strings"
	"unicode"
)

// repeatStr repeat string with one str (#)
func repeatStr(count uint8, str string) string {
	return strings.Repeat(str, int(count))
}

// numberOfChar returns the number of char in a string
func numberOfChar(str, char string) int {
	return strings.Count(str, char)
}

// isFeasible checks if the code generation is feasible
func isFeasible(charset, pattern, char string, count uint64) bool {
	ls := numberOfChar(pattern, char)
	return math.Pow(float64(len(charset)), float64(ls)) >= float64(count)
}

// isAlphanumeric checks if a string is alphanumeric
func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// randomChar returns a random char from a charset
func randomChar(cs []byte) string {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(cs))))
	return string(cs[n.Int64()])
}
