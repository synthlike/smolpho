package httpapi

import (
	"sync"
	"time"

	smolphoindexer "github.com/synthlike/smolpho/indexer/internal/indexer"
)

// RuntimeStatus is an immutable snapshot of indexing-loop health.
type RuntimeStatus struct {
	Syncing            bool
	Head               uint64
	HeadKnown          bool
	Pending            bool
	LastError          string
	LastSuccessfulSync time.Time
	ReorgCount         uint64
}

// StatusTracker receives indexing-loop updates and makes them safe for HTTP
// handlers to read concurrently.
type StatusTracker struct {
	mu     sync.RWMutex
	status RuntimeStatus
}

func NewStatusTracker() *StatusTracker {
	return &StatusTracker{}
}

// Observe implements the indexer's synchronization status callback.
func (t *StatusTracker) Observe(update smolphoindexer.SyncStatus) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if update.Syncing {
		t.status.Syncing = true
		return
	}

	t.status.Syncing = false
	if update.HeadKnown {
		t.status.Head = update.Head
		t.status.HeadKnown = true
	}
	t.status.Pending = update.Pending
	if update.Err != nil {
		t.status.LastError = update.Err.Error()
		return
	}
	t.status.LastError = ""
	t.status.LastSuccessfulSync = time.Now().UTC()
	if update.Replayed {
		t.status.ReorgCount++
	}
}

func (t *StatusTracker) Snapshot() RuntimeStatus {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.status
}
