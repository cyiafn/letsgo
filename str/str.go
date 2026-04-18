// Package str provides small unicode-aware string classification helpers.
package str

import (
	"unicode"
)

// IsAlphaNumeric reports whether every rune in v is a Unicode letter or
// Unicode number.
//
// Semantics:
//   - Pure: v is not mutated.
//   - UNICODE-AWARE, not ASCII-only. Non-Latin letters and non-Arabic digits
//     count (e.g. "漢字", "ünîcødé", "१२३"). If you need ASCII-only
//     classification, pre-filter or write an explicit check.
//   - The empty string returns true (vacuous truth — there is no rune that
//     violates the predicate).
//   - Any whitespace, punctuation, symbol, or control character causes
//     false. That includes the space character.
//   - Iterates by rune via `range`, so the comparison is per-codepoint, not
//     per-byte — multibyte UTF-8 sequences are handled correctly.
//
// Example:
//
//	IsAlphaNumeric("abc123")   // -> true
//	IsAlphaNumeric("ünîcødé")  // -> true
//	IsAlphaNumeric("漢字42")    // -> true
//	IsAlphaNumeric("")         // -> true (vacuously)
//	IsAlphaNumeric("abc 123")  // -> false (space)
//	IsAlphaNumeric("hello!")   // -> false (punctuation)
func IsAlphaNumeric(v string) bool {
	for _, char := range v {
		if !unicode.IsNumber(char) && !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

// IsAlpha reports whether every rune in v is a Unicode letter.
//
// Semantics:
//   - Pure: v is not mutated.
//   - UNICODE-AWARE. Non-Latin letters are accepted (e.g. "漢字", "ünîcødé").
//   - The empty string returns true (vacuous truth).
//   - Digits, whitespace, punctuation, and symbols all cause false.
//   - Iterates by rune, so multibyte UTF-8 is handled correctly.
//
// Example:
//
//	IsAlpha("abc")      // -> true
//	IsAlpha("ünîcødé")  // -> true
//	IsAlpha("漢字")      // -> true
//	IsAlpha("")         // -> true (vacuously)
//	IsAlpha("abc1")     // -> false (digit)
//	IsAlpha("a b")      // -> false (space)
func IsAlpha(v string) bool {
	for _, char := range v {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

// IsNumeric reports whether every rune in v is a Unicode number.
//
// Semantics:
//   - Pure: v is not mutated.
//   - UNICODE-AWARE: accepts any digit script recognized by unicode.IsNumber
//     (e.g. Devanagari "१२३"), not just ASCII '0'..'9'. If you need strict
//     ASCII-decimal, use strconv.Atoi or check r >= '0' && r <= '9' yourself.
//   - DOES NOT parse numeric literals: a leading sign ("-1", "+1"), a
//     decimal point ("1.5"), or scientific notation ("1e3") ALL return false
//     because those characters are not Unicode number runes.
//   - The empty string returns true (vacuous truth).
//
// Example:
//
//	IsNumeric("123")    // -> true
//	IsNumeric("0")      // -> true
//	IsNumeric("१२३")    // -> true (Devanagari digits)
//	IsNumeric("")       // -> true (vacuously)
//	IsNumeric("12a")    // -> false (letter)
//	IsNumeric("-1")     // -> false (sign is not a number rune)
//	IsNumeric("1.5")    // -> false (dot is not a number rune)
func IsNumeric(v string) bool {
	for _, char := range v {
		if !unicode.IsNumber(char) {
			return false
		}
	}
	return true
}
