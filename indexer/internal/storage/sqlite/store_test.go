package sqlite

import (
	"context"
	"math/big"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/synthlike/smolpho/indexer/internal/state"
	"github.com/synthlike/smolpho/indexer/internal/storage"
	"github.com/synthlike/smolpho/indexer/internal/storage/storagetest"
)

func TestStoreConformance(t *testing.T) {
	storagetest.Run(t, func(t *testing.T, initialTimestamp uint64) storage.Store {
		store, err := Open(
			context.Background(),
			filepath.Join(t.TempDir(), "indexer.sqlite"),
			initialTimestamp,
		)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = store.Close() })
		return store
	})
}

func TestStatePersistsAcrossReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "indexer.sqlite")
	store, err := Open(context.Background(), path, 100)
	if err != nil {
		t.Fatal(err)
	}
	checkpoint := storage.Checkpoint{
		Number: 9,
		Hash:   common.HexToHash("0x1234"),
		Valid:  true,
	}
	if err := store.Commit(context.Background(), []state.Event{
		state.Supplied{User: "alice", Assets: big.NewInt(25), Shares: big.NewInt(50)},
		state.CollateralSupplied{User: "alice", Assets: big.NewInt(7)},
		state.Borrowed{User: "alice", Assets: big.NewInt(8), Shares: big.NewInt(16)},
	}, checkpoint); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}

	reopened, err := Open(context.Background(), path, 999)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = reopened.Close() })
	snapshot, err := reopened.Snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Checkpoint != checkpoint {
		t.Fatalf("checkpoint = %+v, want %+v", snapshot.Checkpoint, checkpoint)
	}
	if got := snapshot.State.Market.LastUpdate; got.Cmp(big.NewInt(100)) != 0 {
		t.Fatalf("lastUpdate = %s, want persisted value 100", got)
	}
	if got := snapshot.State.Market.TotalSupplyAssets; got.Cmp(big.NewInt(25)) != 0 {
		t.Fatalf("totalSupplyAssets = %s, want 25", got)
	}
	if got := snapshot.State.Market.TotalBorrowAssets; got.Cmp(big.NewInt(8)) != 0 {
		t.Fatalf("totalBorrowAssets = %s, want 8", got)
	}
	position := snapshot.State.Positions["alice"]
	if position == nil {
		t.Fatal("alice position is missing")
	}
	if got := position.BorrowShares; got.Cmp(big.NewInt(16)) != 0 {
		t.Fatalf("borrowShares = %s, want 16", got)
	}
	if got := position.Collateral; got.Cmp(big.NewInt(7)) != 0 {
		t.Fatalf("collateral = %s, want 7", got)
	}
}

func TestFailedCommitRollsBackStateAndCheckpoint(t *testing.T) {
	store, err := Open(
		context.Background(),
		filepath.Join(t.TempDir(), "indexer.sqlite"),
		0,
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })

	checkpoint := storage.Checkpoint{Number: 1, Hash: common.HexToHash("0x01"), Valid: true}
	if err := store.Commit(context.Background(), []state.Event{state.Supplied{
		User: "alice", Assets: big.NewInt(10), Shares: big.NewInt(20),
	}}, checkpoint); err != nil {
		t.Fatal(err)
	}
	if _, err := store.db.Exec(`
		CREATE TRIGGER reject_bob BEFORE INSERT ON positions
		WHEN NEW.user = 'bob'
		BEGIN
			SELECT RAISE(ABORT, 'reject bob');
		END`); err != nil {
		t.Fatal(err)
	}

	err = store.Commit(context.Background(), []state.Event{state.Supplied{
		User: "bob", Assets: big.NewInt(30), Shares: big.NewInt(60),
	}}, storage.Checkpoint{Number: 2, Hash: common.HexToHash("0x02"), Valid: true})
	if err == nil {
		t.Fatal("Commit() succeeded, want trigger failure")
	}

	snapshot, err := store.Snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Checkpoint != checkpoint {
		t.Fatalf("checkpoint = %+v, want unchanged %+v", snapshot.Checkpoint, checkpoint)
	}
	if _, exists := snapshot.State.Positions["bob"]; exists {
		t.Fatal("failed commit persisted bob's position")
	}
	if got := snapshot.State.Market.TotalSupplyAssets; got.Cmp(big.NewInt(10)) != 0 {
		t.Fatalf("totalSupplyAssets = %s, want unchanged 10", got)
	}
}
