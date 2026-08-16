package indexer

import (
	"context"
	"errors"
	"math/big"
	"testing"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/synthlike/smolpho/indexer/internal/bindings"
	"github.com/synthlike/smolpho/indexer/internal/state"
)

type fakeLogFilterer struct {
	logs    []types.Log
	queries []ethereum.FilterQuery
}

func (f *fakeLogFilterer) FilterLogs(_ context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	f.queries = append(f.queries, query)
	var filtered []types.Log
	for _, log := range f.logs {
		if len(query.Addresses) != 1 || query.Addresses[0] != log.Address {
			continue
		}
		if len(query.Topics) == 0 || len(query.Topics[0]) == 0 || len(log.Topics) == 0 {
			continue
		}
		for _, topic := range query.Topics[0] {
			if topic == log.Topics[0] {
				filtered = append(filtered, log)
				break
			}
		}
	}
	return filtered, nil
}

func (*fakeLogFilterer) SubscribeFilterLogs(
	context.Context,
	ethereum.FilterQuery,
	chan<- types.Log,
) (ethereum.Subscription, error) {
	return nil, errors.New("unexpected log subscription")
}

func TestEthereumEventSourceDecodesDebtAndCollateralEvents(t *testing.T) {
	contract := common.HexToAddress("0x00000000000000000000000000000000000000c0")
	user := common.HexToAddress("0x00000000000000000000000000000000000000a1")
	liquidator := common.HexToAddress("0x00000000000000000000000000000000000000b2")
	blockHash := common.HexToHash("0x1234")
	backend := &fakeLogFilterer{logs: []types.Log{
		makeEventLog(t, contract, blockHash, 10, 1, "CollateralWithdrawn",
			[]common.Address{user}, big.NewInt(11)),
		makeEventLog(t, contract, blockHash, 10, 2, "Borrowed",
			[]common.Address{user}, big.NewInt(12), big.NewInt(13)),
		makeEventLog(t, contract, blockHash, 10, 3, "Repaid",
			[]common.Address{user}, big.NewInt(14), big.NewInt(15)),
		makeEventLog(t, contract, blockHash, 10, 4, "Liquidated",
			[]common.Address{liquidator, user}, big.NewInt(16), big.NewInt(17), big.NewInt(18)),
		makeEventLog(t, contract, blockHash, 10, 5, "BadDebtRealized",
			[]common.Address{user}, big.NewInt(19), big.NewInt(20)),
	}}
	filterer, err := bindings.NewSmolphoFilterer(contract, backend)
	if err != nil {
		t.Fatal(err)
	}
	source := &ethereumEventSource{filterer: filterer}

	events, err := source.Events(context.Background(), 10, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 5 {
		t.Fatalf("decoded events = %d, want 5", len(events))
	}
	if len(backend.queries) != 9 {
		t.Fatalf("filter queries = %d, want one for each of 9 supported event types", len(backend.queries))
	}
	for i, query := range backend.queries {
		if query.FromBlock == nil || query.FromBlock.Uint64() != 10 ||
			query.ToBlock == nil || query.ToBlock.Uint64() != 10 {
			t.Fatalf("query %d range = %v-%v, want 10-10", i, query.FromBlock, query.ToBlock)
		}
	}

	wantUser := addressKey(user)
	wantLiquidator := addressKey(liquidator)
	seen := make(map[string]bool)
	for _, decoded := range events {
		if decoded.block != 10 || decoded.hash != blockHash {
			t.Fatalf("event metadata = block %d hash %s, want block 10 hash %s",
				decoded.block, decoded.hash, blockHash)
		}
		switch event := decoded.event.(type) {
		case state.CollateralWithdrawn:
			seen["CollateralWithdrawn"] = true
			assertAddressAndInts(t, event.User, wantUser, []*big.Int{event.Assets}, 11)
		case state.Borrowed:
			seen["Borrowed"] = true
			assertAddressAndInts(t, event.User, wantUser, []*big.Int{event.Assets, event.Shares}, 12, 13)
		case state.Repaid:
			seen["Repaid"] = true
			assertAddressAndInts(t, event.User, wantUser, []*big.Int{event.Assets, event.Shares}, 14, 15)
		case state.Liquidated:
			seen["Liquidated"] = true
			if event.Liquidator != wantLiquidator || event.Borrower != wantUser {
				t.Fatalf("liquidation parties = %s/%s, want %s/%s",
					event.Liquidator, event.Borrower, wantLiquidator, wantUser)
			}
			assertInts(t, []*big.Int{event.RepaidAssets, event.RepaidShares, event.SeizedCollateral}, 16, 17, 18)
		case state.BadDebtRealized:
			seen["BadDebtRealized"] = true
			assertAddressAndInts(t, event.Borrower, wantUser,
				[]*big.Int{event.BadDebtAssets, event.BadDebtShares}, 19, 20)
		default:
			t.Fatalf("unexpected decoded event %T", event)
		}
	}
	for _, name := range []string{
		"CollateralWithdrawn", "Borrowed", "Repaid", "Liquidated", "BadDebtRealized",
	} {
		if !seen[name] {
			t.Fatalf("%s was not decoded", name)
		}
	}
}

func makeEventLog(
	t *testing.T,
	contract common.Address,
	blockHash common.Hash,
	block uint64,
	index uint,
	name string,
	indexed []common.Address,
	values ...any,
) types.Log {
	t.Helper()
	contractABI, err := bindings.SmolphoMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}
	event, ok := contractABI.Events[name]
	if !ok {
		t.Fatalf("event %s is missing from ABI", name)
	}
	data, err := event.Inputs.NonIndexed().Pack(values...)
	if err != nil {
		t.Fatalf("pack %s: %v", name, err)
	}
	topics := []common.Hash{event.ID}
	for _, address := range indexed {
		topics = append(topics, common.BytesToHash(address.Bytes()))
	}
	return types.Log{
		Address:     contract,
		Topics:      topics,
		Data:        data,
		BlockNumber: block,
		BlockHash:   blockHash,
		Index:       index,
	}
}

func assertAddressAndInts(
	t *testing.T,
	gotAddress string,
	wantAddress string,
	got []*big.Int,
	want ...int64,
) {
	t.Helper()
	if gotAddress != wantAddress {
		t.Fatalf("address = %s, want %s", gotAddress, wantAddress)
	}
	assertInts(t, got, want...)
}

func assertInts(t *testing.T, got []*big.Int, want ...int64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("integer count = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].Cmp(big.NewInt(want[i])) != 0 {
			t.Fatalf("integer %d = %s, want %d", i, got[i], want[i])
		}
	}
}
