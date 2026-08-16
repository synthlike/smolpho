// Package memory provides ephemeral in-memory indexer storage.
package memory

import (
	"context"
	"sync"

	"github.com/synthlike/smolpho/indexer/internal/state"
	"github.com/synthlike/smolpho/indexer/internal/storage"
)

// Store guards one reconstructed state and canonical checkpoint in memory.
type Store struct {
	mu         sync.RWMutex
	state      *state.State
	checkpoint storage.Checkpoint
	closed     bool
}

var _ storage.Store = (*Store)(nil)

// New returns an empty store whose last-update value is initialTimestamp.
func New(initialTimestamp uint64) *Store {
	return &Store{state: state.New(initialTimestamp)}
}

func (s *Store) Checkpoint(ctx context.Context) (storage.Checkpoint, error) {
	if err := ctx.Err(); err != nil {
		return storage.Checkpoint{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return storage.Checkpoint{}, storage.ErrClosed
	}
	return s.checkpoint, nil
}

func (s *Store) Commit(
	ctx context.Context,
	events []state.Event,
	checkpoint storage.Checkpoint,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return storage.ErrClosed
	}
	for _, event := range events {
		s.state.Apply(event)
	}
	s.checkpoint = checkpoint
	return nil
}

func (s *Store) Replace(
	ctx context.Context,
	replacement *state.State,
	checkpoint storage.Checkpoint,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	clone := replacement.Clone()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return storage.ErrClosed
	}
	s.state = clone
	s.checkpoint = checkpoint
	return nil
}

func (s *Store) Snapshot(ctx context.Context) (storage.Snapshot, error) {
	if err := ctx.Err(); err != nil {
		return storage.Snapshot{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return storage.Snapshot{}, storage.ErrClosed
	}
	return storage.Snapshot{
		State:      s.state.Clone(),
		Checkpoint: s.checkpoint,
	}, nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}
