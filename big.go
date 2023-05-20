// Package letsgo provides generic useful functions for Go.
package letsgo

import "math/big"

// ToBigInt converts the given string to big.Int.
func ToBigInt(v string) *big.Int {
	intValue := new(big.Int)
	intValue.SetString(v, 10)
	return intValue
}

// ToBigFloat converts the given string to big.Float.
func ToBigFloat(v string) *big.Float {
	floatValue := new(big.Float)
	floatValue.SetString(v)
	return floatValue
}
