package main

import (
	"context"
	"fmt"
	"log"

	smolphoindexer "github.com/synthlike/smolpho/indexer/internal/indexer"
	"github.com/synthlike/smolpho/indexer/internal/storage"
	"github.com/synthlike/smolpho/indexer/internal/storage/sqlite"
)

type apiConfig struct {
	Indexer  smolphoindexer.Config
	Database string
}

func run(ctx context.Context, config apiConfig) error {
	if err := config.Indexer.Validate(); err != nil {
		return err
	}
	store, err := sqlite.Open(ctx, config.Database, 0)
	if err != nil {
		return fmt.Errorf("open indexer database: %w", err)
	}
	defer store.Close()

	logger := log.Default()
	return smolphoindexer.Run(ctx, config.Indexer, smolphoindexer.Dependencies{
		Store:  store,
		Logger: logger,
		HandleSnapshot: func(_ context.Context, snapshot storage.Snapshot) error {
			logger.Printf("indexed canonical block %d (%s)",
				snapshot.Checkpoint.Number, snapshot.Checkpoint.Hash.Hex())
			return nil
		},
	})
}
