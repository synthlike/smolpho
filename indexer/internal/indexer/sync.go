package indexer

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/synthlike/smolpho/indexer/internal/state"
	"github.com/synthlike/smolpho/indexer/internal/storage"
)

const maxCanonicalRetries = 3

var errChainChanged = errors.New("canonical chain changed while indexing")

type chainReader interface {
	HeaderByNumber(context.Context, *big.Int) (*types.Header, error)
}

type syncResult struct {
	Head     uint64
	Changed  bool
	Replayed bool
	Pending  bool
}

func syncWithRetry(
	ctx context.Context,
	chain chainReader,
	source eventSource,
	st storage.Store,
	deploymentBlock uint64,
	batchSize uint64,
) (syncResult, error) {
	var result syncResult
	var err error
	var replayed bool
	for range maxCanonicalRetries {
		result, err = syncOnce(ctx, chain, source, st, deploymentBlock, batchSize)
		replayed = replayed || result.Replayed
		result.Replayed = replayed
		if !errors.Is(err, errChainChanged) {
			return result, err
		}
	}
	return result, fmt.Errorf("chain changed during %d consecutive attempts: %w", maxCanonicalRetries, err)
}

func syncOnce(
	ctx context.Context,
	chain chainReader,
	source eventSource,
	st storage.Store,
	deploymentBlock uint64,
	batchSize uint64,
) (syncResult, error) {
	head, err := chain.HeaderByNumber(ctx, nil)
	if err != nil {
		return syncResult{}, fmt.Errorf("get chain head: %w", err)
	}
	result := syncResult{Head: head.Number.Uint64()}
	checkpoint, err := st.Checkpoint(ctx)
	if err != nil {
		return result, fmt.Errorf("read checkpoint: %w", err)
	}
	if result.Head < deploymentBlock {
		if checkpoint.Valid {
			if err := st.BeginRebuild(ctx, state.New(0)); err != nil {
				return result, fmt.Errorf("reset state before deployment: %w", err)
			}
			if _, err := st.PublishRebuild(ctx); err != nil {
				return result, fmt.Errorf("publish reset before deployment: %w", err)
			}
			result.Changed = true
			result.Replayed = true
		}
		result.Pending = true
		return result, nil
	}

	deploymentHeader, err := headerByNumber(ctx, chain, deploymentBlock)
	if err != nil {
		return result, fmt.Errorf("get deployment block %d: %w", deploymentBlock, err)
	}
	deploymentHash := deploymentHeader.Hash()

	if checkpoint.Valid {
		canonical := checkpoint.Number <= result.Head
		if canonical {
			header, headerErr := headerByNumber(ctx, chain, checkpoint.Number)
			if headerErr != nil {
				return result, fmt.Errorf("verify checkpoint block %d: %w", checkpoint.Number, headerErr)
			}
			canonical = header.Hash() == checkpoint.Hash
		}
		if !canonical {
			checkpoint = storage.Checkpoint{}
			result.Changed = true
			result.Replayed = true
		}
	}

	if !checkpoint.Valid {
		if err := st.BeginRebuild(ctx, state.New(deploymentHeader.Time)); err != nil {
			return result, fmt.Errorf("initialize state: %w", err)
		}
	}

	start := deploymentBlock
	if checkpoint.Valid {
		if checkpoint.Number == ^uint64(0) {
			published, err := st.PublishRebuild(ctx)
			if err != nil {
				return result, fmt.Errorf("publish rebuilt state: %w", err)
			}
			result.Changed = result.Changed || published
			return result, nil
		}
		start = checkpoint.Number + 1
	}
	for start <= result.Head {
		end := result.Head
		if batchSize-1 < result.Head-start {
			end = start + batchSize - 1
		}
		if err := indexCanonicalRange(
			ctx, chain, source, st, deploymentBlock, deploymentHash, start, end,
		); err != nil {
			return result, err
		}
		result.Changed = true
		if end == ^uint64(0) {
			break
		}
		start = end + 1
	}
	published, err := st.PublishRebuild(ctx)
	if err != nil {
		return result, fmt.Errorf("publish rebuilt state: %w", err)
	}
	result.Changed = result.Changed || published
	return result, nil
}

func indexCanonicalRange(
	ctx context.Context,
	chain chainReader,
	source eventSource,
	st storage.Store,
	deploymentBlock uint64,
	deploymentHash common.Hash,
	from uint64,
	to uint64,
) error {
	endBefore, err := headerByNumber(ctx, chain, to)
	if err != nil {
		return fmt.Errorf("get range end block %d: %w", to, err)
	}
	events, err := source.Events(ctx, from, to)
	if err != nil {
		return fmt.Errorf("get events %d-%d: %w", from, to, err)
	}

	headerCache := make(map[uint64]*types.Header)
	for i := range events {
		event := &events[i]
		if event.block < from || event.block > to {
			return fmt.Errorf("event block %d outside requested range %d-%d", event.block, from, to)
		}
		header := headerCache[event.block]
		if header == nil {
			header, err = headerByNumber(ctx, chain, event.block)
			if err != nil {
				return fmt.Errorf("get event block %d: %w", event.block, err)
			}
			headerCache[event.block] = header
		}
		if event.hash != (common.Hash{}) && event.hash != header.Hash() {
			return errChainChanged
		}
		if accrued, ok := event.event.(state.InterestAccrued); ok {
			accrued.Timestamp = header.Time
			event.event = accrued
		}
	}

	sortOrderedEvents(events)

	endAfter, err := headerByNumber(ctx, chain, to)
	if err != nil {
		return fmt.Errorf("recheck range end block %d: %w", to, err)
	}
	deploymentAfter, err := headerByNumber(ctx, chain, deploymentBlock)
	if err != nil {
		return fmt.Errorf("recheck deployment block %d: %w", deploymentBlock, err)
	}
	if endBefore.Hash() != endAfter.Hash() || deploymentHash != deploymentAfter.Hash() {
		return errChainChanged
	}

	decoded := make([]state.Event, len(events))
	for i := range events {
		decoded[i] = events[i].event
	}
	checkpoint := storage.Checkpoint{Number: to, Hash: endAfter.Hash(), Valid: true}
	if err := st.Commit(ctx, decoded, checkpoint); err != nil {
		return fmt.Errorf("commit range %d-%d: %w", from, to, err)
	}
	return nil
}

func sortOrderedEvents(events []orderedEvent) {
	slices.SortStableFunc(events, func(a, b orderedEvent) int {
		if a.block != b.block {
			return cmp.Compare(a.block, b.block)
		}
		return cmp.Compare(a.index, b.index)
	})
}

func headerByNumber(ctx context.Context, chain chainReader, number uint64) (*types.Header, error) {
	return chain.HeaderByNumber(ctx, new(big.Int).SetUint64(number))
}
