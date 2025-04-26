package utils

import "strings"

func CamelToSnake(s string) string {
	var builder strings.Builder
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			builder.WriteRune('_')
		}
		builder.WriteRune(r)
	}
	return strings.ToLower(builder.String())
}
