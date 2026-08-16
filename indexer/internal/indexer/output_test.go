package indexer

import (
	"math/big"
	"testing"
)

func TestFormatInt(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "0", want: "0"},
		{input: "12", want: "12"},
		{input: "123", want: "123"},
		{input: "1234", want: "1_234"},
		{input: "50000000000000000000", want: "50_000_000_000_000_000_000"},
		{input: "50000000000000000000000000", want: "50_000_000_000_000_000_000_000_000"},
		{input: "-1234567", want: "-1_234_567"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			input, ok := new(big.Int).SetString(test.input, 10)
			if !ok {
				t.Fatalf("invalid test input %q", test.input)
			}
			if got := formatInt(input); got != test.want {
				t.Fatalf("formatInt(%s) = %q, want %q", input, got, test.want)
			}
		})
	}
}
