package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/synthlike/smolpho/indexer/internal/httpapi"
	smolphoindexer "github.com/synthlike/smolpho/indexer/internal/indexer"
	"github.com/synthlike/smolpho/indexer/internal/storage"
	"github.com/synthlike/smolpho/indexer/internal/storage/sqlite"
)

type apiConfig struct {
	Indexer  smolphoindexer.Config
	Database string
	Listen   string
}

func run(ctx context.Context, config apiConfig) error {
	if err := config.Indexer.Validate(); err != nil {
		return err
	}
	if config.Listen == "" {
		return fmt.Errorf("-listen is required")
	}
	store, err := sqlite.Open(ctx, config.Database, 0)
	if err != nil {
		return fmt.Errorf("open indexer database: %w", err)
	}
	defer store.Close()

	listener, err := net.Listen("tcp", config.Listen)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", config.Listen, err)
	}
	defer listener.Close()

	runContext, cancel := context.WithCancel(ctx)
	defer cancel()
	logger := log.Default()
	status := httpapi.NewStatusTracker()
	server := &http.Server{
		Handler:           httpapi.NewHandler(store, status),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return runContext
		},
	}

	workerErrors := make(chan error, 1)
	go func() {
		status.Observe(smolphoindexer.SyncStatus{Syncing: true})
		workerErrors <- smolphoindexer.Run(
			runContext,
			config.Indexer,
			smolphoindexer.Dependencies{
				Store:            store,
				Logger:           logger,
				HandleSyncStatus: status.Observe,
				HandleSnapshot: func(_ context.Context, snapshot storage.Snapshot) error {
					logger.Printf("indexed canonical block %d (%s)",
						snapshot.Checkpoint.Number, snapshot.Checkpoint.Hash.Hex())
					return nil
				},
			},
		)
	}()
	serverErrors := make(chan error, 1)
	go func() {
		logger.Printf("serving indexer API on http://%s", listener.Addr())
		serverErrors <- server.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		cancel()
		shutdownErr := shutdownServer(server)
		<-workerErrors
		return shutdownErr
	case workerErr := <-workerErrors:
		cancel()
		shutdownErr := shutdownServer(server)
		if ctx.Err() != nil {
			return shutdownErr
		}
		if workerErr != nil {
			return workerErr
		}
		return shutdownErr
	case serverErr := <-serverErrors:
		cancel()
		workerErr := <-workerErrors
		if serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
			return fmt.Errorf("serve indexer API: %w", serverErr)
		}
		return workerErr
	}
}

func shutdownServer(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("shut down indexer API: %w", err)
	}
	return nil
}
