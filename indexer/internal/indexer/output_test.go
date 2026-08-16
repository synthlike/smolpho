package indexer

import (
	"bytes"
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/synthlike/smolpho/indexer/internal/state"
	"github.com/synthlike/smolpho/indexer/internal/storage"
	"github.com/synthlike/smolpho/indexer/internal/storage/memory"
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

func TestPrintStateIncludesBorrowPosition(t *testing.T) {
	store := memory.New(100)
	if err := store.Commit(context.Background(), []state.Event{state.Borrowed{
		User: "alice", Assets: big.NewInt(1_000), Shares: big.NewInt(1_000_000_000),
	}}, storage.Checkpoint{
		Number: 7,
		Hash:   common.HexToHash("0x1234"),
		Valid:  true,
	}); err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	if err := printState(context.Background(), &output, store); err != nil {
		t.Fatal(err)
	}
	want := "borrowShares=1_000_000_000  borrowAssets=1_000"
	if !strings.Contains(output.String(), want) {
		t.Fatalf("output does not contain %q:\n%s", want, output.String())
	}
}
