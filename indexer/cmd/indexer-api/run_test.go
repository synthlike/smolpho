package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	smolphoindexer "github.com/synthlike/smolpho/indexer/internal/indexer"
)

func TestRunValidatesIndexerBeforeOpeningDatabase(t *testing.T) {
	database := filepath.Join(t.TempDir(), "indexer.sqlite")
	deploymentBlock := uint64(1)
	err := run(context.Background(), apiConfig{
		Indexer: smolphoindexer.Config{
			RPC:             "http://not-used.invalid",
			Contract:        "invalid",
			DeploymentBlock: &deploymentBlock,
			Follow:          false,
			Interval:        time.Second,
			BatchSize:       100,
		},
		Database: database,
	})
	if err == nil {
		t.Fatal("run() succeeded with an invalid contract")
	}
	if _, statErr := os.Stat(database); !os.IsNotExist(statErr) {
		t.Fatalf("database was created before validation: %v", statErr)
	}
}

func TestRunRequiresListenAddressBeforeOpeningDatabase(t *testing.T) {
	database := filepath.Join(t.TempDir(), "indexer.sqlite")
	deploymentBlock := uint64(1)
	err := run(context.Background(), apiConfig{
		Indexer: smolphoindexer.Config{
			RPC:             "http://not-used.invalid",
			Contract:        "0x0000000000000000000000000000000000000001",
			DeploymentBlock: &deploymentBlock,
			Follow:          true,
			Interval:        time.Second,
			BatchSize:       100,
		},
		Database: database,
	})
	if err == nil {
		t.Fatal("run() succeeded without a listen address")
	}
	if _, statErr := os.Stat(database); !os.IsNotExist(statErr) {
		t.Fatalf("database was created before validation: %v", statErr)
	}
}
