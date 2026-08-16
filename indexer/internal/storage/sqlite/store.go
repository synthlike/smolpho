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

const schemaVersion = 1

// Store persists one materialized indexer state and its canonical checkpoint.
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

func (s *Store) initialize(ctx context.Context, initialTimestamp uint64) error {
	if _, err := s.db.ExecContext(ctx, "PRAGMA busy_timeout = 5000"); err != nil {
		return fmt.Errorf("configure sqlite busy timeout: %w", err)
	}
	if _, err := s.db.ExecContext(ctx, "PRAGMA journal_mode = WAL"); err != nil {
		return fmt.Errorf("configure sqlite journal mode: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin schema transaction: %w", err)
	}
	defer tx.Rollback()

	statements := []string{
		`CREATE TABLE IF NOT EXISTS storage_metadata (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			schema_version INTEGER NOT NULL,
			checkpoint_number TEXT NOT NULL,
			checkpoint_hash BLOB NOT NULL,
			checkpoint_valid INTEGER NOT NULL CHECK (checkpoint_valid IN (0, 1))
		)`,
		`CREATE TABLE IF NOT EXISTS market (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			total_supply_assets TEXT NOT NULL,
			total_supply_shares TEXT NOT NULL,
			total_borrow_assets TEXT NOT NULL,
			total_borrow_shares TEXT NOT NULL,
			last_update TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS positions (
			user TEXT PRIMARY KEY,
			supply_shares TEXT NOT NULL,
			borrow_shares TEXT NOT NULL,
			collateral TEXT NOT NULL
		)`,
	}
	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("create sqlite schema: %w", err)
		}
	}

	zeroHash := common.Hash{}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO storage_metadata (
			id, schema_version, checkpoint_number, checkpoint_hash, checkpoint_valid
		) VALUES (1, ?, '0', ?, 0)
		ON CONFLICT(id) DO NOTHING`, schemaVersion, zeroHash[:]); err != nil {
		return fmt.Errorf("initialize storage metadata: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO market (
			id, total_supply_assets, total_supply_shares,
			total_borrow_assets, total_borrow_shares, last_update
		) VALUES (1, '0', '0', '0', '0', ?)
		ON CONFLICT(id) DO NOTHING`, strconv.FormatUint(initialTimestamp, 10)); err != nil {
		return fmt.Errorf("initialize market: %w", err)
	}

	var version int
	if err := tx.QueryRowContext(ctx,
		"SELECT schema_version FROM storage_metadata WHERE id = 1",
	).Scan(&version); err != nil {
		return fmt.Errorf("read sqlite schema version: %w", err)
	}
	if version != schemaVersion {
		return fmt.Errorf("unsupported sqlite schema version %d (want %d)", version, schemaVersion)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit schema transaction: %w", err)
	}
	return nil
}

// Checkpoint returns the last fully committed canonical block.
func (s *Store) Checkpoint(ctx context.Context) (storage.Checkpoint, error) {
	if err := ctx.Err(); err != nil {
		return storage.Checkpoint{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return storage.Checkpoint{}, storage.ErrClosed
	}
	checkpoint, err := readCheckpoint(ctx, s.db)
	if err != nil {
		return storage.Checkpoint{}, fmt.Errorf("read checkpoint: %w", err)
	}
	return checkpoint, nil
}

// Commit applies events and advances the checkpoint in one transaction.
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

	current, err := readState(ctx, tx)
	if err != nil {
		return fmt.Errorf("read state for commit: %w", err)
	}
	for _, event := range events {
		current.Apply(event)
	}
	if err := writeSnapshot(ctx, tx, current, checkpoint); err != nil {
		return fmt.Errorf("write committed state: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit state transaction: %w", err)
	}
	return nil
}

// Replace swaps all reconstructed state and checkpoint in one transaction.
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
	if err := writeSnapshot(ctx, tx, replacement, checkpoint); err != nil {
		return fmt.Errorf("write replacement state: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit replacement transaction: %w", err)
	}
	return nil
}

// Snapshot reads state and checkpoint from one consistent transaction.
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
	current, err := readState(ctx, tx)
	if err != nil {
		return storage.Snapshot{}, fmt.Errorf("read snapshot state: %w", err)
	}
	checkpoint, err := readCheckpoint(ctx, tx)
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

type queryer interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func readCheckpoint(ctx context.Context, db queryer) (storage.Checkpoint, error) {
	var (
		numberText string
		hashBytes  []byte
		valid      bool
	)
	if err := db.QueryRowContext(ctx, `
		SELECT checkpoint_number, checkpoint_hash, checkpoint_valid
		FROM storage_metadata WHERE id = 1`,
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

func readState(ctx context.Context, db stateQueryer) (*state.State, error) {
	var marketValues [5]string
	if err := db.QueryRowContext(ctx, `
		SELECT total_supply_assets, total_supply_shares, total_borrow_assets,
		       total_borrow_shares, last_update
		FROM market WHERE id = 1`,
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
		SELECT user, supply_shares, borrow_shares, collateral FROM positions`,
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

func writeSnapshot(
	ctx context.Context,
	tx *sql.Tx,
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
		WHERE id = 1`,
		market.TotalSupplyAssets.String(), market.TotalSupplyShares.String(),
		market.TotalBorrowAssets.String(), market.TotalBorrowShares.String(),
		market.LastUpdate.String(),
	); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM positions"); err != nil {
		return err
	}
	statement, err := tx.PrepareContext(ctx, `
		INSERT INTO positions (user, supply_shares, borrow_shares, collateral)
		VALUES (?, ?, ?, ?)`)
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
			ctx, user, position.SupplyShares.String(),
			position.BorrowShares.String(), position.Collateral.String(),
		); err != nil {
			return err
		}
	}
	hash := checkpoint.Hash
	if _, err := tx.ExecContext(ctx, `
		UPDATE storage_metadata SET
			checkpoint_number = ?, checkpoint_hash = ?, checkpoint_valid = ?
		WHERE id = 1`,
		strconv.FormatUint(checkpoint.Number, 10), hash[:], checkpoint.Valid,
	); err != nil {
		return err
	}
	return nil
}

func requireBigInts(values ...*big.Int) error {
	for i, value := range values {
		if value == nil {
			return fmt.Errorf("integer %d is nil", i)
		}
	}
	return nil
}
