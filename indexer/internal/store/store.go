// Package store holds the reconstructed state in memory. It is a thin
// concurrency-safe wrapper so a background log tailer can apply events while
// the CLI reads a consistent snapshot for printing. State is rebuilt from a
// backfill on each start; nothing is persisted.
package store

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"github.com/synthlike/smolpho/indexer/internal/state"
)

// Checkpoint identifies the last canonical block whose logs were fully
// applied.
type Checkpoint struct {
	Number uint64
	Hash   common.Hash
	Valid  bool
}

// Store guards a single reconstructed State and its canonical checkpoint.
type Store struct {
	mu         sync.RWMutex
	state      *state.State
	checkpoint Checkpoint
}

// New returns an empty in-memory store.
func New(initialTimestamp uint64) *Store {
	return &Store{state: state.New(initialTimestamp)}
}

// Reset discards reconstructed state and starts again from the deployment
// block timestamp.
func (s *Store) Reset(initialTimestamp uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = state.New(initialTimestamp)
	s.checkpoint = Checkpoint{}
}

// Commit atomically applies a canonical range and advances its checkpoint.
func (s *Store) Commit(events []state.Event, number uint64, hash common.Hash) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, event := range events {
		s.state.Apply(event)
	}
	s.checkpoint = Checkpoint{Number: number, Hash: hash, Valid: true}
}

// Checkpoint returns the last fully-applied canonical block, if any.
func (s *Store) Checkpoint() Checkpoint {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.checkpoint
}

// Read runs fn against the live state and checkpoint under the read lock.
// Callers must not retain references beyond the callback.
func (s *Store) Read(fn func(*state.State, Checkpoint)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fn(s.state, s.checkpoint)
}
