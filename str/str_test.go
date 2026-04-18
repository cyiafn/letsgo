package str

import "testing"

func TestIsAlphaNumeric(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"abc123", true},
		{"ABCxyz", true},
		{"123456", true},
		{"", true},
		{"abc 123", false},
		{"hello!", false},
		{"a-b", false},
		{"ünîcødé", true},
		{"漢字42", true},
	}
	for _, tc := range tests {
		if got := IsAlphaNumeric(tc.in); got != tc.want {
			t.Errorf("IsAlphaNumeric(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"abc", true},
		{"ABC", true},
		{"", true},
		{"abc1", false},
		{"a b", false},
		{"ünîcødé", true},
		{"漢字", true},
		{"42", false},
	}
	for _, tc := range tests {
		if got := IsAlpha(tc.in); got != tc.want {
			t.Errorf("IsAlpha(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"123", true},
		{"", true},
		{"12a", false},
		{"1.5", false},
		{"-1", false},
		{"0", true},
		{"१२३", true},
	}
	for _, tc := range tests {
		if got := IsNumeric(tc.in); got != tc.want {
			t.Errorf("IsNumeric(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
