// Package storagetest contains behavioral tests shared by storage backends.
package storagetest

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/synthlike/smolpho/indexer/internal/state"
	"github.com/synthlike/smolpho/indexer/internal/storage"
)

// Factory creates a fresh store for one test case.
type Factory func(*testing.T, uint64) storage.Store

// Run executes the storage contract against a backend implementation.
func Run(t *testing.T, factory Factory) {
	t.Helper()
	t.Run("initial snapshot", func(t *testing.T) {
		store := factory(t, 100)
		snapshot := mustSnapshot(t, store)
		if snapshot.Checkpoint.Valid {
			t.Fatalf("initial checkpoint = %+v, want invalid", snapshot.Checkpoint)
		}
		if got := snapshot.State.Market.LastUpdate; got.Cmp(big.NewInt(100)) != 0 {
			t.Fatalf("initial lastUpdate = %s, want 100", got)
		}
	})

	t.Run("atomic commit", func(t *testing.T) {
		store := factory(t, 100)
		checkpoint := storage.Checkpoint{
			Number: 7,
			Hash:   common.HexToHash("0x1234"),
			Valid:  true,
		}
		err := store.Commit(context.Background(), []state.Event{state.Supplied{
			User: "alice", Assets: big.NewInt(10), Shares: big.NewInt(20),
		}}, checkpoint)
		if err != nil {
			t.Fatal(err)
		}

		snapshot := mustSnapshot(t, store)
		if snapshot.Checkpoint != checkpoint {
			t.Fatalf("checkpoint = %+v, want %+v", snapshot.Checkpoint, checkpoint)
		}
		if got := snapshot.State.Market.TotalSupplyAssets; got.Cmp(big.NewInt(10)) != 0 {
			t.Fatalf("totalSupplyAssets = %s, want 10", got)
		}
	})

	t.Run("snapshot isolation", func(t *testing.T) {
		store := factory(t, 100)
		checkpoint := storage.Checkpoint{Number: 1, Valid: true}
		if err := store.Commit(context.Background(), []state.Event{state.Supplied{
			User: "alice", Assets: big.NewInt(10), Shares: big.NewInt(20),
		}}, checkpoint); err != nil {
			t.Fatal(err)
		}

		snapshot := mustSnapshot(t, store)
		snapshot.State.Market.TotalSupplyAssets.SetInt64(999)
		snapshot.State.Positions["alice"].SupplyShares.SetInt64(999)
		delete(snapshot.State.Positions, "alice")

		fresh := mustSnapshot(t, store)
		if got := fresh.State.Market.TotalSupplyAssets; got.Cmp(big.NewInt(10)) != 0 {
			t.Fatalf("stored assets mutated through snapshot: %s", got)
		}
		if got := fresh.State.Positions["alice"].SupplyShares; got.Cmp(big.NewInt(20)) != 0 {
			t.Fatalf("stored shares mutated through snapshot: %s", got)
		}
	})

	t.Run("atomic replace", func(t *testing.T) {
		store := factory(t, 100)
		if err := store.Commit(context.Background(), []state.Event{state.Supplied{
			User: "alice", Assets: big.NewInt(10), Shares: big.NewInt(20),
		}}, storage.Checkpoint{Number: 1, Valid: true}); err != nil {
			t.Fatal(err)
		}

		replacement := state.New(200)
		replacement.Apply(state.CollateralSupplied{User: "bob", Assets: big.NewInt(30)})
		checkpoint := storage.Checkpoint{
			Number: 2,
			Hash:   common.HexToHash("0x5678"),
			Valid:  true,
		}
		if err := store.Replace(context.Background(), replacement, checkpoint); err != nil {
			t.Fatal(err)
		}
		replacement.Positions["bob"].Collateral.SetInt64(999)

		snapshot := mustSnapshot(t, store)
		if snapshot.Checkpoint != checkpoint {
			t.Fatalf("checkpoint = %+v, want %+v", snapshot.Checkpoint, checkpoint)
		}
		if _, exists := snapshot.State.Positions["alice"]; exists {
			t.Fatal("replace retained an old position")
		}
		if got := snapshot.State.Positions["bob"].Collateral; got.Cmp(big.NewInt(30)) != 0 {
			t.Fatalf("replacement retained mutable alias: %s", got)
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		store := factory(t, 0)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		checkpoint := storage.Checkpoint{
			Number: 12,
			Hash:   common.HexToHash("0xcafe"),
			Valid:  true,
		}
		err := store.Commit(ctx, []state.Event{state.Supplied{
			User: "alice", Assets: big.NewInt(42), Shares: big.NewInt(84),
		}}, checkpoint)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Commit() error = %v, want context.Canceled", err)
		}

		snapshot := mustSnapshot(t, store)
		if snapshot.Checkpoint.Valid {
			t.Fatalf("checkpoint = %+v, want invalid", snapshot.Checkpoint)
		}
		if got := snapshot.State.Market.TotalSupplyAssets; got.Sign() != 0 {
			t.Fatalf("totalSupplyAssets = %s, want 0", got)
		}

		if _, err = store.Snapshot(ctx); !errors.Is(err, context.Canceled) {
			t.Fatalf("Snapshot() error = %v, want context.Canceled", err)
		}
	})

	t.Run("closed store", func(t *testing.T) {
		store := factory(t, 0)
		if err := store.Close(); err != nil {
			t.Fatal(err)
		}
		if _, err := store.Snapshot(context.Background()); !errors.Is(err, storage.ErrClosed) {
			t.Fatalf("Snapshot() error = %v, want ErrClosed", err)
		}
	})
}

func mustSnapshot(t *testing.T, store storage.Store) storage.Snapshot {
	t.Helper()
	snapshot, err := store.Snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}
