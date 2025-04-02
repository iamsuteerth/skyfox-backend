package utils

import (
	"strings"
	"unicode"
)

func ToCamelCase(name string) string {
	words := strings.Fields(strings.ToLower(name))
	for i, word := range words {
		if len(word) > 0 {
			r := []rune(word)
			r[0] = unicode.ToUpper(r[0])
			words[i] = string(r)
		}
	}
	return strings.Join(words, " ")
}
