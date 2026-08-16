package indexer

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/synthlike/smolpho/indexer/internal/state"
	"github.com/synthlike/smolpho/indexer/internal/storage"
	"github.com/synthlike/smolpho/indexer/internal/storage/memory"
)

type fakeChain struct {
	head    uint64
	headers map[uint64]*types.Header
}

func (c *fakeChain) HeaderByNumber(_ context.Context, number *big.Int) (*types.Header, error) {
	n := c.head
	if number != nil {
		n = number.Uint64()
	}
	header := c.headers[n]
	if header == nil {
		return nil, errors.New("header not found")
	}
	return header, nil
}

type blockRange struct {
	from uint64
	to   uint64
}

type fakeEventSource struct {
	events []orderedEvent
	err    error
	calls  []blockRange
	after  func()
}

func (s *fakeEventSource) Events(_ context.Context, from, to uint64) ([]orderedEvent, error) {
	s.calls = append(s.calls, blockRange{from: from, to: to})
	var events []orderedEvent
	for _, event := range s.events {
		if event.block >= from && event.block <= to {
			events = append(events, event)
		}
	}
	if s.after != nil {
		s.after()
	}
	return events, s.err
}

func newHeader(number, timestamp uint64, branch byte) *types.Header {
	return &types.Header{
		Number: new(big.Int).SetUint64(number),
		Time:   timestamp,
		Extra:  []byte{branch},
	}
}

func supplied(block uint64, index uint, header *types.Header, user string, assets, shares int64) orderedEvent {
	return orderedEvent{
		block: block,
		index: index,
		hash:  header.Hash(),
		event: state.Supplied{
			User: user, Assets: big.NewInt(assets), Shares: big.NewInt(shares),
		},
	}
}

func TestSyncWaitsForFutureDeploymentBlock(t *testing.T) {
	chain := &fakeChain{
		head:    5,
		headers: map[uint64]*types.Header{5: newHeader(5, 50, 1)},
	}
	source := &fakeEventSource{}
	st := memory.New(0)

	result, err := syncWithRetry(context.Background(), chain, source, st, 10, 100)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Pending || result.Head != 5 {
		t.Fatalf("result = %+v, want pending at head 5", result)
	}
	checkpoint, err := st.Checkpoint(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if checkpoint.Valid || len(source.calls) != 0 {
		t.Fatal("future deployment unexpectedly advanced or queried events")
	}

	deployment := newHeader(10, 100, 1)
	chain.head = 10
	chain.headers[10] = deployment
	source.events = []orderedEvent{supplied(10, 1, deployment, "alice", 10, 20)}
	result, err = syncWithRetry(context.Background(), chain, source, st, 10, 100)
	if err != nil {
		t.Fatal(err)
	}
	if result.Pending || len(source.calls) != 1 || source.calls[0] != (blockRange{10, 10}) {
		t.Fatalf("result/calls = %+v/%+v, want one deployment-block query", result, source.calls)
	}
	assertPositionShares(t, st, "alice", 20)
}

func TestSyncStartsAtDeploymentAndBatchesRanges(t *testing.T) {
	chain := &fakeChain{head: 14, headers: make(map[uint64]*types.Header)}
	for block := uint64(10); block <= chain.head; block++ {
		chain.headers[block] = newHeader(block, block*10, 1)
	}
	source := &fakeEventSource{}
	st := memory.New(0)

	if _, err := syncWithRetry(context.Background(), chain, source, st, 10, 2); err != nil {
		t.Fatal(err)
	}
	want := []blockRange{{10, 11}, {12, 13}, {14, 14}}
	if len(source.calls) != len(want) {
		t.Fatalf("calls = %+v, want %+v", source.calls, want)
	}
	for i := range want {
		if source.calls[i] != want[i] {
			t.Fatalf("call %d = %+v, want %+v", i, source.calls[i], want[i])
		}
	}
	checkpoint, err := st.Checkpoint(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if checkpoint.Number != 14 || checkpoint.Hash != chain.headers[14].Hash() {
		t.Fatalf("checkpoint = %+v, want canonical block 14", checkpoint)
	}
}

func TestSyncReplaysSameHeightReorg(t *testing.T) {
	deployment := newHeader(10, 100, 1)
	oldHead := newHeader(11, 110, 1)
	chain := &fakeChain{
		head: 11,
		headers: map[uint64]*types.Header{
			10: deployment,
			11: oldHead,
		},
	}
	source := &fakeEventSource{events: []orderedEvent{
		supplied(11, 1, oldHead, "alice", 10, 20),
	}}
	st := memory.New(0)
	if _, err := syncWithRetry(context.Background(), chain, source, st, 10, 100); err != nil {
		t.Fatal(err)
	}

	newHead := newHeader(11, 111, 2)
	chain.headers[11] = newHead
	source.events = []orderedEvent{supplied(11, 1, newHead, "bob", 30, 40)}
	result, err := syncWithRetry(context.Background(), chain, source, st, 10, 100)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Replayed {
		t.Fatalf("result = %+v, want reorg replay", result)
	}
	assertPositionShares(t, st, "alice", 0)
	assertPositionShares(t, st, "bob", 40)
}

func TestSyncRetriesWhenChainChangesDuringQuery(t *testing.T) {
	deployment := newHeader(10, 100, 1)
	oldHead := newHeader(11, 110, 1)
	newHead := newHeader(11, 111, 2)
	chain := &fakeChain{
		head: 11,
		headers: map[uint64]*types.Header{
			10: deployment,
			11: oldHead,
		},
	}
	source := &fakeEventSource{events: []orderedEvent{
		supplied(11, 1, oldHead, "alice", 10, 20),
	}}
	source.after = func() {
		chain.headers[11] = newHead
		source.events = []orderedEvent{supplied(11, 1, newHead, "bob", 30, 40)}
		source.after = nil
	}
	st := memory.New(0)

	if _, err := syncWithRetry(context.Background(), chain, source, st, 10, 100); err != nil {
		t.Fatal(err)
	}
	if len(source.calls) != 2 {
		t.Fatalf("event source calls = %d, want one failed attempt and one retry", len(source.calls))
	}
	assertPositionShares(t, st, "alice", 0)
	assertPositionShares(t, st, "bob", 40)
}

func TestCanonicalRangeIsAtomicOnSourceError(t *testing.T) {
	header := newHeader(10, 100, 1)
	chain := &fakeChain{head: 10, headers: map[uint64]*types.Header{10: header}}
	st := memory.New(100)
	initialHash := common.HexToHash("0x1234")
	initialCheckpoint := storage.Checkpoint{Number: 9, Hash: initialHash, Valid: true}
	if err := st.Commit(context.Background(), []state.Event{state.Supplied{
		User: "alice", Assets: big.NewInt(10), Shares: big.NewInt(20),
	}}, initialCheckpoint); err != nil {
		t.Fatal(err)
	}
	source := &fakeEventSource{
		events: []orderedEvent{supplied(10, 1, header, "bob", 30, 40)},
		err:    errors.New("rpc failed"),
	}

	err := indexCanonicalRange(context.Background(), chain, source, st, 10, header.Hash(), 10, 10)
	if err == nil {
		t.Fatal("indexCanonicalRange succeeded despite source error")
	}
	checkpoint, checkpointErr := st.Checkpoint(context.Background())
	if checkpointErr != nil {
		t.Fatal(checkpointErr)
	}
	if checkpoint.Number != 9 || checkpoint.Hash != initialHash {
		t.Fatalf("checkpoint changed after failed range: %+v", checkpoint)
	}
	assertPositionShares(t, st, "alice", 20)
	assertPositionShares(t, st, "bob", 0)
}

func TestInterestAccrualUsesEventBlockTimestamp(t *testing.T) {
	header := newHeader(10, 777, 1)
	chain := &fakeChain{head: 10, headers: map[uint64]*types.Header{10: header}}
	source := &fakeEventSource{events: []orderedEvent{{
		block: 10,
		index: 1,
		hash:  header.Hash(),
		event: state.InterestAccrued{Elapsed: big.NewInt(7), Interest: big.NewInt(3)},
	}}}
	st := memory.New(0)
	if _, err := syncWithRetry(context.Background(), chain, source, st, 10, 100); err != nil {
		t.Fatal(err)
	}
	snapshot, err := st.Snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := snapshot.State.Market.LastUpdate; got.Cmp(big.NewInt(777)) != 0 {
		t.Fatalf("lastUpdate = %s, want 777", got)
	}
}

func TestSyncAppliesCompleteMarketLifecycleInLogOrder(t *testing.T) {
	header := newHeader(10, 777, 1)
	chain := &fakeChain{head: 10, headers: map[uint64]*types.Header{10: header}}
	event := func(index uint, value state.Event) orderedEvent {
		return orderedEvent{block: 10, index: index, hash: header.Hash(), event: value}
	}
	// Deliberately scramble the result returned by the event source. The sync
	// layer must restore EVM log order before applying dependent transitions.
	source := &fakeEventSource{events: []orderedEvent{
		event(7, state.BadDebtRealized{
			Borrower: "borrower", BadDebtAssets: big.NewInt(2_800), BadDebtShares: big.NewInt(2_800_000_000),
		}),
		event(4, state.Borrowed{
			User: "borrower", Assets: big.NewInt(3_000), Shares: big.NewInt(3_000_000_000),
		}),
		event(1, state.Supplied{
			User: "lender", Assets: big.NewInt(5_000), Shares: big.NewInt(5_000_000_000),
		}),
		event(6, state.Liquidated{
			Liquidator:       "liquidator",
			Borrower:         "borrower",
			RepaidAssets:     big.NewInt(100),
			RepaidShares:     big.NewInt(100_000_000),
			SeizedCollateral: big.NewInt(80),
		}),
		event(2, state.CollateralSupplied{User: "borrower", Assets: big.NewInt(100)}),
		event(5, state.Repaid{
			User: "borrower", Assets: big.NewInt(100), Shares: big.NewInt(100_000_000),
		}),
		event(3, state.CollateralWithdrawn{User: "borrower", Assets: big.NewInt(20)}),
	}}
	store := memory.New(0)

	if _, err := syncWithRetry(context.Background(), chain, source, store, 10, 100); err != nil {
		t.Fatal(err)
	}
	snapshot, err := store.Snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	market := snapshot.State.Market
	if market.TotalSupplyAssets.Cmp(big.NewInt(2_200)) != 0 {
		t.Fatalf("totalSupplyAssets = %s, want 2200", market.TotalSupplyAssets)
	}
	if market.TotalBorrowAssets.Sign() != 0 || market.TotalBorrowShares.Sign() != 0 {
		t.Fatalf("borrow totals = (%s assets, %s shares), want zero",
			market.TotalBorrowAssets, market.TotalBorrowShares)
	}
	borrower := snapshot.State.Positions["borrower"]
	if borrower == nil {
		t.Fatal("borrower position is missing")
	}
	if borrower.BorrowShares.Sign() != 0 || borrower.Collateral.Sign() != 0 {
		t.Fatalf("borrower = %+v, want zero debt and collateral", borrower)
	}
}

func TestSortOrderedEvents(t *testing.T) {
	events := []orderedEvent{
		{block: 2, index: 1},
		{block: 1, index: 9},
		{block: 1, index: 3},
	}
	sortOrderedEvents(events)
	want := []blockRange{{1, 3}, {1, 9}, {2, 1}}
	for i := range want {
		if events[i].block != want[i].from || uint64(events[i].index) != want[i].to {
			t.Fatalf("event %d = (%d,%d), want (%d,%d)",
				i, events[i].block, events[i].index, want[i].from, want[i].to)
		}
	}
}

func assertPositionShares(t *testing.T, st storage.Store, user string, want int64) {
	t.Helper()
	snapshot, err := st.Snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	position := snapshot.State.Positions[user]
	if position == nil {
		if want == 0 {
			return
		}
		t.Fatalf("position %q missing, want %d shares", user, want)
	}
	if got := position.SupplyShares; got.Cmp(big.NewInt(want)) != 0 {
		t.Fatalf("position %q shares = %s, want %d", user, got, want)
	}
}
