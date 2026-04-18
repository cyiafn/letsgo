package letsgo

import (
	"math/big"
	"testing"
)

func TestToBigInt(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want *big.Int
	}{
		{"zero", "0", big.NewInt(0)},
		{"positive", "12345", big.NewInt(12345)},
		{"negative", "-99", big.NewInt(-99)},
		{"large", "123456789012345678901234567890", func() *big.Int {
			b, _ := new(big.Int).SetString("123456789012345678901234567890", 10)
			return b
		}()},
		{"invalid returns zero", "not-a-number", big.NewInt(0)},
		{"empty returns zero", "", big.NewInt(0)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ToBigInt(tc.in)
			if got.Cmp(tc.want) != 0 {
				t.Fatalf("ToBigInt(%q) = %s, want %s", tc.in, got.String(), tc.want.String())
			}
		})
	}
}

func TestToBigFloat(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"zero", "0", "0"},
		{"positive", "3.14", "3.14"},
		{"negative", "-2.5", "-2.5"},
		{"integer", "100", "100"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ToBigFloat(tc.in)
			want, _, _ := big.ParseFloat(tc.want, 10, got.Prec(), big.ToNearestEven)
			if got.Cmp(want) != 0 {
				t.Fatalf("ToBigFloat(%q) = %s, want %s", tc.in, got.String(), want.String())
			}
		})
	}
}
