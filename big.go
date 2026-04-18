// Package letsgo provides generic useful functions for Go.
package letsgo

import "math/big"

// ToBigInt parses a base-10 numeric string into a *big.Int.
//
// Semantics:
//   - Pure: the input string is not mutated and no global state is touched.
//   - Never returns nil; a fresh *big.Int is always allocated.
//   - If v is empty or not a valid base-10 integer, the underlying
//     big.Int.SetString call leaves the value as its zero (0). No error
//     is surfaced — callers that need to distinguish "0" from "invalid"
//     should validate v themselves before calling.
//   - Leading '+' / '-' signs are accepted; whitespace and underscores are not.
//
// Example:
//
//	ToBigInt("12345")                            // -> *big.Int value 12345
//	ToBigInt("-99")                              // -> *big.Int value -99
//	ToBigInt("123456789012345678901234567890")   // -> arbitrary precision value
//	ToBigInt("not-a-number")                     // -> *big.Int value 0
//	ToBigInt("")                                 // -> *big.Int value 0
func ToBigInt(v string) *big.Int {
	intValue := new(big.Int)
	intValue.SetString(v, 10)
	return intValue
}

// ToBigFloat parses a decimal string into a *big.Float using base-10.
//
// Semantics:
//   - Pure: the input string is not mutated.
//   - Never returns nil; a fresh *big.Float is always allocated.
//   - Precision is whatever big.Float.SetString picks for the literal
//     (usually enough to represent the mantissa exactly). This is NOT
//     float64's 53-bit precision — comparing against a big.NewFloat(x)
//     value will often fail even when the decimals look identical.
//     When comparing, parse the expected value via big.ParseFloat with
//     the same precision, e.g. got.Prec().
//   - Invalid input leaves the value as 0 (no error is surfaced).
//   - Accepts standard decimal notation and scientific notation
//     (e.g. "1.5e3"), with optional leading sign.
//
// Example:
//
//	ToBigFloat("3.14")    // -> *big.Float ≈ 3.14
//	ToBigFloat("-2.5")    // -> *big.Float = -2.5
//	ToBigFloat("1e2")     // -> *big.Float = 100
//	ToBigFloat("nope")    // -> *big.Float = 0
func ToBigFloat(v string) *big.Float {
	floatValue := new(big.Float)
	floatValue.SetString(v)
	return floatValue
}
