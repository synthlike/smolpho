package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

const schemaVersion = 2

func (s *Store) initialize(ctx context.Context, initialTimestamp uint64) error {
	pragmas := []struct {
		statement string
		name      string
	}{
		{statement: "PRAGMA busy_timeout = 5000", name: "busy timeout"},
		{statement: "PRAGMA journal_mode = WAL", name: "journal mode"},
		{statement: "PRAGMA foreign_keys = ON", name: "foreign keys"},
	}
	for _, pragma := range pragmas {
		if _, err := s.db.ExecContext(ctx, pragma.statement); err != nil {
			return fmt.Errorf("configure sqlite %s: %w", pragma.name, err)
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin schema transaction: %w", err)
	}
	defer tx.Rollback()

	exists, err := tableExists(ctx, tx, "storage_metadata")
	if err != nil {
		return fmt.Errorf("inspect sqlite schema: %w", err)
	}
	if !exists {
		if err := createSchema(ctx, tx, initialTimestamp); err != nil {
			return fmt.Errorf("create sqlite schema: %w", err)
		}
	} else {
		version, err := readSchemaVersion(ctx, tx)
		if err != nil {
			return fmt.Errorf("read sqlite schema version: %w", err)
		}
		if version != schemaVersion {
			return fmt.Errorf("unsupported sqlite schema version %d (want %d)", version, schemaVersion)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit schema transaction: %w", err)
	}
	return nil
}

func tableExists(ctx context.Context, tx *sql.Tx, name string) (bool, error) {
	var found string
	err := tx.QueryRowContext(ctx, `
		SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`, name,
	).Scan(&found)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func readSchemaVersion(ctx context.Context, tx *sql.Tx) (int, error) {
	var version int
	err := tx.QueryRowContext(ctx,
		"SELECT schema_version FROM storage_metadata WHERE id = 1",
	).Scan(&version)
	return version, err
}

func createSchema(ctx context.Context, tx *sql.Tx, initialTimestamp uint64) error {
	for _, statement := range schemaStatements() {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	zeroHash := common.Hash{}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO storage_metadata (
			id, schema_version, active_generation, staging_generation, next_generation
		) VALUES (1, ?, 1, NULL, 2)`, schemaVersion); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO generations (
			id, checkpoint_number, checkpoint_hash, checkpoint_valid
		) VALUES (1, '0', ?, 0)`, zeroHash[:]); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO market (
			generation, total_supply_assets, total_supply_shares,
			total_borrow_assets, total_borrow_shares, last_update
		) VALUES (1, '0', '0', '0', '0', ?)`, strconv.FormatUint(initialTimestamp, 10)); err != nil {
		return err
	}
	return nil
}

func schemaStatements() []string {
	return []string{
		`CREATE TABLE storage_metadata (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			schema_version INTEGER NOT NULL,
			active_generation INTEGER NOT NULL,
			staging_generation INTEGER,
			next_generation INTEGER NOT NULL
		)`,
		`CREATE TABLE generations (
			id INTEGER PRIMARY KEY,
			checkpoint_number TEXT NOT NULL,
			checkpoint_hash BLOB NOT NULL,
			checkpoint_valid INTEGER NOT NULL CHECK (checkpoint_valid IN (0, 1))
		)`,
		`CREATE TABLE market (
			generation INTEGER PRIMARY KEY REFERENCES generations(id) ON DELETE CASCADE,
			total_supply_assets TEXT NOT NULL,
			total_supply_shares TEXT NOT NULL,
			total_borrow_assets TEXT NOT NULL,
			total_borrow_shares TEXT NOT NULL,
			last_update TEXT NOT NULL
		)`,
		`CREATE TABLE positions (
			generation INTEGER NOT NULL REFERENCES generations(id) ON DELETE CASCADE,
			user TEXT NOT NULL,
			supply_shares TEXT NOT NULL,
			borrow_shares TEXT NOT NULL,
			collateral TEXT NOT NULL,
			PRIMARY KEY (generation, user)
		)`,
	}
}
