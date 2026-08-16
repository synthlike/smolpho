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
	// Checkpoint returns indexing progress for the working projection. During
	// a rebuild this is the hidden rebuild checkpoint, not the published one.
	Checkpoint(context.Context) (Checkpoint, error)

	// Commit atomically applies an ordered event batch and advances the
	// checkpoint. Neither change may become visible without the other.
	Commit(context.Context, []state.Event, Checkpoint) error

	// Replace atomically replaces all reconstructed state and its checkpoint.
	// Implementations must not retain mutable aliases to replacement.
	Replace(context.Context, *state.State, Checkpoint) error

	// BeginRebuild starts a hidden replacement projection with an invalid
	// checkpoint. Commit advances this working projection while Snapshot keeps
	// returning the last published projection.
	BeginRebuild(context.Context, *state.State) error

	// PublishRebuild atomically makes a completed hidden rebuild visible. The
	// bool reports whether a rebuild was published; it is false for a no-op.
	PublishRebuild(context.Context) (bool, error)

	// Snapshot returns a deep, point-in-time copy of the published state and
	// checkpoint. Hidden rebuild progress is never exposed.
	Snapshot(context.Context) (Snapshot, error)

	Close() error
}
