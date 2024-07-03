package str

import (
	"unicode"
)

func IsAlphaNumeric(v string) bool {
	for _, char := range v {
		if !unicode.IsNumber(char) && !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

func IsAlpha(v string) bool {
	for _, char := range v {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

func IsNumeric(v string) bool {
	for _, char := range v {
		if !unicode.IsNumber(char) {
			return false
		}
	}
	return true
}
