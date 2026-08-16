// Package storage defines persistence semantics for reconstructed indexer
// state. Implementations may be ephemeral or durable.
package storage

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"

	"github.com/synthlike/smolpho/indexer/internal/state"
)

// ErrClosed is returned after a store has been closed.
var ErrClosed = errors.New("storage is closed")

// Checkpoint identifies the last canonical block whose logs were fully
// applied.
type Checkpoint struct {
	Number uint64
	Hash   common.Hash
	Valid  bool
}

// Snapshot is an immutable point-in-time view of state and its checkpoint.
// Store implementations must return a deep copy that callers may retain.
type Snapshot struct {
	State      *state.State
	Checkpoint Checkpoint
}

// Store persists reconstructed state and canonical indexing progress.
type Store interface {
	// Checkpoint returns the last fully-applied canonical block, if any.
	Checkpoint(context.Context) (Checkpoint, error)

	// Commit atomically applies an ordered event batch and advances the
	// checkpoint. Neither change may become visible without the other.
	Commit(context.Context, []state.Event, Checkpoint) error

	// Replace atomically replaces all reconstructed state and its checkpoint.
	// Implementations must not retain mutable aliases to replacement.
	Replace(context.Context, *state.State, Checkpoint) error

	// Snapshot returns a deep, point-in-time copy of state and checkpoint.
	Snapshot(context.Context) (Snapshot, error)

	Close() error
}
