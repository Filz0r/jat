package utils

import (
	"strings"
	"unicode"
)

func ToSnakeCase(text string) string {
	fields := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	return strings.ToLower(strings.Join(fields, "_"))
}
