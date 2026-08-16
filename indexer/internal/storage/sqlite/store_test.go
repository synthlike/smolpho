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

func TestHiddenRebuildPersistsAcrossReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "indexer.sqlite")
	store, err := Open(context.Background(), path, 100)
	if err != nil {
		t.Fatal(err)
	}
	publishedCheckpoint := storage.Checkpoint{
		Number: 1, Hash: common.HexToHash("0x1111"), Valid: true,
	}
	if err := store.Commit(context.Background(), []state.Event{state.Supplied{
		User: "alice", Assets: big.NewInt(10), Shares: big.NewInt(20),
	}}, publishedCheckpoint); err != nil {
		t.Fatal(err)
	}
	if err := store.BeginRebuild(context.Background(), state.New(200)); err != nil {
		t.Fatal(err)
	}
	workingCheckpoint := storage.Checkpoint{
		Number: 2, Hash: common.HexToHash("0x2222"), Valid: true,
	}
	if err := store.Commit(context.Background(), []state.Event{state.CollateralSupplied{
		User: "bob", Assets: big.NewInt(30),
	}}, workingCheckpoint); err != nil {
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
	published, err := reopened.Snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if published.Checkpoint != publishedCheckpoint {
		t.Fatalf("published checkpoint = %+v, want %+v", published.Checkpoint, publishedCheckpoint)
	}
	if _, exists := published.State.Positions["bob"]; exists {
		t.Fatal("reopen exposed hidden rebuild state")
	}
	working, err := reopened.Checkpoint(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if working != workingCheckpoint {
		t.Fatalf("working checkpoint = %+v, want %+v", working, workingCheckpoint)
	}
	publishedRebuild, err := reopened.PublishRebuild(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !publishedRebuild {
		t.Fatal("PublishRebuild() reported no publication")
	}
	published, err = reopened.Snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if published.Checkpoint != workingCheckpoint {
		t.Fatalf("published checkpoint = %+v, want %+v", published.Checkpoint, workingCheckpoint)
	}
	if _, exists := published.State.Positions["alice"]; exists {
		t.Fatal("published rebuild retained old state")
	}
	if got := published.State.Positions["bob"].Collateral; got.Cmp(big.NewInt(30)) != 0 {
		t.Fatalf("bob collateral = %s, want 30", got)
	}
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

func TestFailedRebuildCommitLeavesPublishedGenerationVisible(t *testing.T) {
	store, err := Open(
		context.Background(),
		filepath.Join(t.TempDir(), "indexer.sqlite"),
		0,
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })

	publishedCheckpoint := storage.Checkpoint{
		Number: 1, Hash: common.HexToHash("0x01"), Valid: true,
	}
	if err := store.Commit(context.Background(), []state.Event{state.Supplied{
		User: "alice", Assets: big.NewInt(10), Shares: big.NewInt(20),
	}}, publishedCheckpoint); err != nil {
		t.Fatal(err)
	}
	if err := store.BeginRebuild(context.Background(), state.New(100)); err != nil {
		t.Fatal(err)
	}
	if _, err := store.db.Exec(`
		CREATE TRIGGER reject_bob_rebuild BEFORE INSERT ON positions
		WHEN NEW.user = 'bob'
		BEGIN
			SELECT RAISE(ABORT, 'reject bob rebuild');
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
	if snapshot.Checkpoint != publishedCheckpoint {
		t.Fatalf("published checkpoint = %+v, want %+v", snapshot.Checkpoint, publishedCheckpoint)
	}
	if _, exists := snapshot.State.Positions["bob"]; exists {
		t.Fatal("failed rebuild exposed bob's position")
	}
	if got := snapshot.State.Positions["alice"].SupplyShares; got.Cmp(big.NewInt(20)) != 0 {
		t.Fatalf("published alice shares = %s, want 20", got)
	}
	working, err := store.Checkpoint(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if working.Valid {
		t.Fatalf("working checkpoint = %+v, want unchanged invalid checkpoint", working)
	}
}
