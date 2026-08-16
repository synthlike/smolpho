// Package sqlite provides durable SQLite-backed indexer storage.
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	_ "modernc.org/sqlite"

	"github.com/synthlike/smolpho/indexer/internal/state"
	"github.com/synthlike/smolpho/indexer/internal/storage"
)

// Store persists published and optionally staged materialized projections.
type Store struct {
	mu     sync.RWMutex
	db     *sql.DB
	closed bool
}

var _ storage.Store = (*Store)(nil)

// Open opens or creates a SQLite store. initialTimestamp is used only when
// initializing a new database; an existing materialized state is preserved.
func Open(ctx context.Context, dataSourceName string, initialTimestamp uint64) (*Store, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if dataSourceName == "" {
		return nil, fmt.Errorf("sqlite data source name is required")
	}

	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	// A store serializes access through one connection. This also makes :memory:
	// data sources safe, because all operations observe the same database.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	store := &Store{db: db}
	if err := store.initialize(ctx, initialTimestamp); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

// Checkpoint returns progress for the hidden rebuild when one exists, or the
// published generation otherwise.
func (s *Store) Checkpoint(ctx context.Context) (storage.Checkpoint, error) {
	if err := ctx.Err(); err != nil {
		return storage.Checkpoint{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return storage.Checkpoint{}, storage.ErrClosed
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return storage.Checkpoint{}, fmt.Errorf("begin checkpoint transaction: %w", err)
	}
	defer tx.Rollback()
	selection, err := readGenerationSelection(ctx, tx)
	if err != nil {
		return storage.Checkpoint{}, fmt.Errorf("read generation selection: %w", err)
	}
	checkpoint, err := readCheckpoint(ctx, tx, selection.working())
	if err != nil {
		return storage.Checkpoint{}, fmt.Errorf("read checkpoint: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return storage.Checkpoint{}, fmt.Errorf("finish checkpoint transaction: %w", err)
	}
	return checkpoint, nil
}

// Commit applies events to the hidden rebuild when present, or directly to
// the published generation during normal forward indexing.
func (s *Store) Commit(
	ctx context.Context,
	events []state.Event,
	checkpoint storage.Checkpoint,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return storage.ErrClosed
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin commit transaction: %w", err)
	}
	defer tx.Rollback()
	selection, err := readGenerationSelection(ctx, tx)
	if err != nil {
		return fmt.Errorf("read generation selection: %w", err)
	}
	generation := selection.working()
	current, err := readState(ctx, tx, generation)
	if err != nil {
		return fmt.Errorf("read state for commit: %w", err)
	}
	for _, event := range events {
		current.Apply(event)
	}
	if err := writeSnapshot(ctx, tx, generation, current, checkpoint); err != nil {
		return fmt.Errorf("write committed state: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit state transaction: %w", err)
	}
	return nil
}

// Replace immediately replaces the published projection and abandons any
// hidden rebuild. Indexing replays use BeginRebuild and PublishRebuild instead.
func (s *Store) Replace(
	ctx context.Context,
	replacement *state.State,
	checkpoint storage.Checkpoint,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if replacement == nil {
		return fmt.Errorf("replacement state is nil")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return storage.ErrClosed
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin replace transaction: %w", err)
	}
	defer tx.Rollback()
	selection, err := readGenerationSelection(ctx, tx)
	if err != nil {
		return fmt.Errorf("read generation selection: %w", err)
	}
	newGeneration := selection.Next
	if err := createGeneration(ctx, tx, newGeneration, replacement, checkpoint); err != nil {
		return fmt.Errorf("create replacement generation: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE storage_metadata
		SET active_generation = ?, staging_generation = NULL, next_generation = ?
		WHERE id = 1`, newGeneration, newGeneration+1); err != nil {
		return fmt.Errorf("select replacement generation: %w", err)
	}
	if err := deleteOtherGenerations(ctx, tx, newGeneration); err != nil {
		return fmt.Errorf("delete replaced generations: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit replacement transaction: %w", err)
	}
	return nil
}

// BeginRebuild creates a hidden working generation. Starting another rebuild
// atomically discards any previous unpublished generation.
func (s *Store) BeginRebuild(ctx context.Context, replacement *state.State) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if replacement == nil {
		return fmt.Errorf("replacement state is nil")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return storage.ErrClosed
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin rebuild transaction: %w", err)
	}
	defer tx.Rollback()
	selection, err := readGenerationSelection(ctx, tx)
	if err != nil {
		return fmt.Errorf("read generation selection: %w", err)
	}
	newGeneration := selection.Next
	if err := createGeneration(
		ctx, tx, newGeneration, replacement, storage.Checkpoint{},
	); err != nil {
		return fmt.Errorf("create rebuild generation: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE storage_metadata
		SET staging_generation = ?, next_generation = ?
		WHERE id = 1`, newGeneration, newGeneration+1); err != nil {
		return fmt.Errorf("select rebuild generation: %w", err)
	}
	if selection.Staging.Valid {
		if err := deleteGeneration(ctx, tx, selection.Staging.Int64); err != nil {
			return fmt.Errorf("delete abandoned rebuild generation: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit rebuild transaction: %w", err)
	}
	return nil
}

// PublishRebuild atomically switches readers to the completed hidden
// generation and removes the previously published generation.
func (s *Store) PublishRebuild(ctx context.Context) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return false, storage.ErrClosed
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("begin publish transaction: %w", err)
	}
	defer tx.Rollback()
	selection, err := readGenerationSelection(ctx, tx)
	if err != nil {
		return false, fmt.Errorf("read generation selection: %w", err)
	}
	if !selection.Staging.Valid {
		return false, nil
	}
	newActive := selection.Staging.Int64
	if _, err := tx.ExecContext(ctx, `
		UPDATE storage_metadata
		SET active_generation = ?, staging_generation = NULL
		WHERE id = 1`, newActive); err != nil {
		return false, fmt.Errorf("publish rebuild generation: %w", err)
	}
	if err := deleteOtherGenerations(ctx, tx, newActive); err != nil {
		return false, fmt.Errorf("delete superseded generation: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("commit publish transaction: %w", err)
	}
	return true, nil
}

// Snapshot reads only the published generation from one consistent
// transaction. Hidden rebuild progress is not visible.
func (s *Store) Snapshot(ctx context.Context) (storage.Snapshot, error) {
	if err := ctx.Err(); err != nil {
		return storage.Snapshot{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return storage.Snapshot{}, storage.ErrClosed
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return storage.Snapshot{}, fmt.Errorf("begin snapshot transaction: %w", err)
	}
	defer tx.Rollback()
	selection, err := readGenerationSelection(ctx, tx)
	if err != nil {
		return storage.Snapshot{}, fmt.Errorf("read generation selection: %w", err)
	}
	current, err := readState(ctx, tx, selection.Active)
	if err != nil {
		return storage.Snapshot{}, fmt.Errorf("read snapshot state: %w", err)
	}
	checkpoint, err := readCheckpoint(ctx, tx, selection.Active)
	if err != nil {
		return storage.Snapshot{}, fmt.Errorf("read snapshot checkpoint: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return storage.Snapshot{}, fmt.Errorf("finish snapshot transaction: %w", err)
	}
	return storage.Snapshot{State: current, Checkpoint: checkpoint}, nil
}

// Close releases the database. It is safe to call more than once.
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	return s.db.Close()
}

type generationSelection struct {
	Active  int64
	Staging sql.NullInt64
	Next    int64
}

func (s generationSelection) working() int64 {
	if s.Staging.Valid {
		return s.Staging.Int64
	}
	return s.Active
}

type queryer interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func readGenerationSelection(ctx context.Context, db queryer) (generationSelection, error) {
	var selection generationSelection
	err := db.QueryRowContext(ctx, `
		SELECT active_generation, staging_generation, next_generation
		FROM storage_metadata WHERE id = 1`,
	).Scan(&selection.Active, &selection.Staging, &selection.Next)
	return selection, err
}

func readCheckpoint(
	ctx context.Context,
	db queryer,
	generation int64,
) (storage.Checkpoint, error) {
	var (
		numberText string
		hashBytes  []byte
		valid      bool
	)
	if err := db.QueryRowContext(ctx, `
		SELECT checkpoint_number, checkpoint_hash, checkpoint_valid
		FROM generations WHERE id = ?`, generation,
	).Scan(&numberText, &hashBytes, &valid); err != nil {
		return storage.Checkpoint{}, err
	}
	number, err := strconv.ParseUint(numberText, 10, 64)
	if err != nil {
		return storage.Checkpoint{}, fmt.Errorf("parse checkpoint number %q: %w", numberText, err)
	}
	if len(hashBytes) != common.HashLength {
		return storage.Checkpoint{}, fmt.Errorf("checkpoint hash has %d bytes, want %d", len(hashBytes), common.HashLength)
	}
	return storage.Checkpoint{
		Number: number,
		Hash:   common.BytesToHash(hashBytes),
		Valid:  valid,
	}, nil
}

type stateQueryer interface {
	queryer
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func readState(ctx context.Context, db stateQueryer, generation int64) (*state.State, error) {
	var marketValues [5]string
	if err := db.QueryRowContext(ctx, `
		SELECT total_supply_assets, total_supply_shares, total_borrow_assets,
		       total_borrow_shares, last_update
		FROM market WHERE generation = ?`, generation,
	).Scan(
		&marketValues[0], &marketValues[1], &marketValues[2],
		&marketValues[3], &marketValues[4],
	); err != nil {
		return nil, err
	}
	marketInts, err := parseBigInts(marketValues[:])
	if err != nil {
		return nil, fmt.Errorf("parse market: %w", err)
	}
	current := &state.State{
		Market: state.Market{
			TotalSupplyAssets: marketInts[0],
			TotalSupplyShares: marketInts[1],
			TotalBorrowAssets: marketInts[2],
			TotalBorrowShares: marketInts[3],
			LastUpdate:        marketInts[4],
		},
		Positions: make(map[string]*state.Position),
	}

	rows, err := db.QueryContext(ctx, `
		SELECT user, supply_shares, borrow_shares, collateral
		FROM positions WHERE generation = ?`, generation,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user string
		var values [3]string
		if err := rows.Scan(&user, &values[0], &values[1], &values[2]); err != nil {
			return nil, err
		}
		ints, err := parseBigInts(values[:])
		if err != nil {
			return nil, fmt.Errorf("parse position %q: %w", user, err)
		}
		current.Positions[user] = &state.Position{
			SupplyShares: ints[0],
			BorrowShares: ints[1],
			Collateral:   ints[2],
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return current, nil
}

func parseBigInts(values []string) ([]*big.Int, error) {
	parsed := make([]*big.Int, len(values))
	for i, value := range values {
		integer, ok := new(big.Int).SetString(value, 10)
		if !ok {
			return nil, fmt.Errorf("invalid integer %q", value)
		}
		parsed[i] = integer
	}
	return parsed, nil
}

func createGeneration(
	ctx context.Context,
	tx *sql.Tx,
	generation int64,
	current *state.State,
	checkpoint storage.Checkpoint,
) error {
	zeroHash := common.Hash{}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO generations (
			id, checkpoint_number, checkpoint_hash, checkpoint_valid
		) VALUES (?, '0', ?, 0)`, generation, zeroHash[:]); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO market (
			generation, total_supply_assets, total_supply_shares,
			total_borrow_assets, total_borrow_shares, last_update
		) VALUES (?, '0', '0', '0', '0', '0')`, generation); err != nil {
		return err
	}
	return writeSnapshot(ctx, tx, generation, current, checkpoint)
}

func writeSnapshot(
	ctx context.Context,
	tx *sql.Tx,
	generation int64,
	current *state.State,
	checkpoint storage.Checkpoint,
) error {
	market := current.Market
	if err := requireBigInts(
		market.TotalSupplyAssets,
		market.TotalSupplyShares,
		market.TotalBorrowAssets,
		market.TotalBorrowShares,
		market.LastUpdate,
	); err != nil {
		return fmt.Errorf("market: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE market SET
			total_supply_assets = ?, total_supply_shares = ?,
			total_borrow_assets = ?, total_borrow_shares = ?, last_update = ?
		WHERE generation = ?`,
		market.TotalSupplyAssets.String(), market.TotalSupplyShares.String(),
		market.TotalBorrowAssets.String(), market.TotalBorrowShares.String(),
		market.LastUpdate.String(), generation,
	); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx,
		"DELETE FROM positions WHERE generation = ?", generation,
	); err != nil {
		return err
	}
	statement, err := tx.PrepareContext(ctx, `
		INSERT INTO positions (
			generation, user, supply_shares, borrow_shares, collateral
		) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer statement.Close()
	for user, position := range current.Positions {
		if position == nil {
			return fmt.Errorf("position %q is nil", user)
		}
		if err := requireBigInts(position.SupplyShares, position.BorrowShares, position.Collateral); err != nil {
			return fmt.Errorf("position %q: %w", user, err)
		}
		if _, err := statement.ExecContext(
			ctx, generation, user, position.SupplyShares.String(),
			position.BorrowShares.String(), position.Collateral.String(),
		); err != nil {
			return err
		}
	}
	hash := checkpoint.Hash
	if _, err := tx.ExecContext(ctx, `
		UPDATE generations SET
			checkpoint_number = ?, checkpoint_hash = ?, checkpoint_valid = ?
		WHERE id = ?`,
		strconv.FormatUint(checkpoint.Number, 10), hash[:], checkpoint.Valid, generation,
	); err != nil {
		return err
	}
	return nil
}

func deleteGeneration(ctx context.Context, tx *sql.Tx, generation int64) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM generations WHERE id = ?", generation)
	return err
}

func deleteOtherGenerations(ctx context.Context, tx *sql.Tx, generation int64) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM generations WHERE id <> ?", generation)
	return err
}

func requireBigInts(values ...*big.Int) error {
	for i, value := range values {
		if value == nil {
			return fmt.Errorf("integer %d is nil", i)
		}
	}
	return nil
}
