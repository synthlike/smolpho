// Package indexer reconstructs Smolpho state from canonical Ethereum logs.
package indexer

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/synthlike/smolpho/indexer/internal/bindings"
	"github.com/synthlike/smolpho/indexer/internal/storage"
)

// Config controls one indexer process.
type Config struct {
	RPC             string
	Contract        string
	DeploymentBlock *uint64
	Follow          bool
	Interval        time.Duration
	BatchSize       uint64
}

// Logger receives operational messages from the indexing loop.
type Logger interface {
	Printf(string, ...any)
}

// SnapshotHandler is called after a successful backfill or follow update.
// The supplied snapshot is isolated from subsequent storage mutations.
type SnapshotHandler func(context.Context, storage.Snapshot) error

// Dependencies are process-owned services used by the shared indexing engine.
// Run does not close Store; the command that constructs it owns its lifecycle.
type Dependencies struct {
	Store          storage.Store
	Logger         Logger
	HandleSnapshot SnapshotHandler
}

// Validate checks configuration before any RPC connection is opened.
func (c Config) Validate() error {
	if c.Contract == "" || !common.IsHexAddress(c.Contract) {
		return fmt.Errorf("-contract must be a valid hex address")
	}
	if c.DeploymentBlock == nil {
		return fmt.Errorf("-deployment-block is required")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("-interval must be greater than zero")
	}
	if c.BatchSize == 0 {
		return fmt.Errorf("-batch-size must be greater than zero")
	}
	return nil
}

// Run backfills a Smolpho deployment and optionally follows new canonical
// blocks until ctx is cancelled.
func Run(ctx context.Context, config Config, dependencies Dependencies) error {
	if err := config.Validate(); err != nil {
		return err
	}
	if dependencies.Store == nil {
		return fmt.Errorf("indexer storage is required")
	}
	deploymentBlock := *config.DeploymentBlock

	client, err := ethclient.DialContext(ctx, config.RPC)
	if err != nil {
		return fmt.Errorf("dial %s: %w", config.RPC, err)
	}
	defer client.Close()

	filterer, err := bindings.NewSmolphoFilterer(common.HexToAddress(config.Contract), client)
	if err != nil {
		return fmt.Errorf("bind contract: %w", err)
	}

	source := &ethereumEventSource{filterer: filterer}
	result, err := syncWithRetry(
		ctx, client, source, dependencies.Store, deploymentBlock, config.BatchSize,
	)
	if err != nil {
		return fmt.Errorf("backfill: %w", err)
	}
	if result.Pending {
		if !config.Follow {
			return fmt.Errorf("deployment block %d is ahead of chain head %d", deploymentBlock, result.Head)
		}
		logf(dependencies.Logger,
			"waiting for deployment block %d (chain head %d)", deploymentBlock, result.Head)
	} else {
		if err := handleSnapshot(ctx, dependencies); err != nil {
			return fmt.Errorf("handle snapshot: %w", err)
		}
	}

	if config.Follow {
		return follow(ctx, client, source, dependencies, config, deploymentBlock)
	}
	return nil
}

func follow(
	ctx context.Context,
	chain chainReader,
	source eventSource,
	dependencies Dependencies,
	config Config,
	deploymentBlock uint64,
) error {
	ticker := time.NewTicker(config.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			result, err := syncWithRetry(
				ctx, chain, source, dependencies.Store, deploymentBlock, config.BatchSize,
			)
			if err != nil {
				logf(dependencies.Logger, "sync: %v", err)
				continue
			}
			if result.Pending || !result.Changed {
				continue
			}
			if result.Replayed {
				logf(dependencies.Logger,
					"canonical chain changed; replayed from deployment block %d", deploymentBlock)
			}
			if err := handleSnapshot(ctx, dependencies); err != nil {
				return fmt.Errorf("handle snapshot: %w", err)
			}
		}
	}
}

func handleSnapshot(ctx context.Context, dependencies Dependencies) error {
	if dependencies.HandleSnapshot == nil {
		return nil
	}
	snapshot, err := dependencies.Store.Snapshot(ctx)
	if err != nil {
		return err
	}
	return dependencies.HandleSnapshot(ctx, snapshot)
}

func logf(logger Logger, format string, args ...any) {
	if logger != nil {
		logger.Printf(format, args...)
	}
}
