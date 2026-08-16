package memory

import (
	"testing"

	"github.com/synthlike/smolpho/indexer/internal/storage"
	"github.com/synthlike/smolpho/indexer/internal/storage/storagetest"
)

func TestStoreConformance(t *testing.T) {
	storagetest.Run(t, func(t *testing.T, initialTimestamp uint64) storage.Store {
		store := New(initialTimestamp)
		t.Cleanup(func() { _ = store.Close() })
		return store
	})
}
