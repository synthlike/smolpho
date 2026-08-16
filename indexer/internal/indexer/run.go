// Package indexer reconstructs Smolpho state from canonical Ethereum logs.
package indexer

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/synthlike/smolpho/indexer/internal/bindings"
	"github.com/synthlike/smolpho/indexer/internal/store"
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
func Run(ctx context.Context, config Config) error {
	if err := config.Validate(); err != nil {
		return err
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
	st := store.New(0)

	result, err := syncWithRetry(ctx, client, source, st, deploymentBlock, config.BatchSize)
	if err != nil {
		return fmt.Errorf("backfill: %w", err)
	}
	if result.Pending {
		if !config.Follow {
			return fmt.Errorf("deployment block %d is ahead of chain head %d", deploymentBlock, result.Head)
		}
		log.Printf("waiting for deployment block %d (chain head %d)", deploymentBlock, result.Head)
	} else {
		printState(os.Stdout, st)
	}

	if config.Follow {
		return follow(ctx, client, source, st, config, deploymentBlock)
	}
	return nil
}

func follow(
	ctx context.Context,
	chain chainReader,
	source eventSource,
	st *store.Store,
	config Config,
	deploymentBlock uint64,
) error {
	ticker := time.NewTicker(config.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Fprintln(os.Stdout, "\nshutting down")
			return nil
		case <-ticker.C:
			result, err := syncWithRetry(ctx, chain, source, st, deploymentBlock, config.BatchSize)
			if err != nil {
				log.Printf("sync: %v", err)
				continue
			}
			if result.Pending || !result.Changed {
				continue
			}
			if result.Replayed {
				log.Printf("canonical chain changed; replayed from deployment block %d", deploymentBlock)
			}
			printState(os.Stdout, st)
		}
	}
}
