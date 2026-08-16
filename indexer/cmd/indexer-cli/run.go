package main

import (
	"context"
	"log"
	"os"

	smolphoindexer "github.com/synthlike/smolpho/indexer/internal/indexer"
	"github.com/synthlike/smolpho/indexer/internal/storage/memory"
)

func run(ctx context.Context, config smolphoindexer.Config) error {
	store := memory.New(0)
	defer store.Close()
	return smolphoindexer.Run(ctx, config, smolphoindexer.Dependencies{
		Store:          store,
		Logger:         log.Default(),
		HandleSnapshot: smolphoindexer.TextSnapshotHandler(os.Stdout),
	})
}
