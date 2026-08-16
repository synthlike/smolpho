package indexer

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/synthlike/smolpho/indexer/internal/bindings"
	"github.com/synthlike/smolpho/indexer/internal/state"
)

type eventSource interface {
	Events(context.Context, uint64, uint64) ([]orderedEvent, error)
}

type ethereumEventSource struct {
	filterer *bindings.SmolphoFilterer
}

// orderedEvent pairs a decoded event with its on-chain position so the reducer
// receives events in exactly the order the contract emitted them.
type orderedEvent struct {
	block uint64
	index uint
	hash  common.Hash
	event state.Event
}

// Events retrieves every supported Smolpho event in [from, to]. Canonical
// validation, timestamp enrichment, sorting, and committing happen separately.
func (s *ethereumEventSource) Events(ctx context.Context, from, to uint64) ([]orderedEvent, error) {
	opts := &bind.FilterOpts{Start: from, End: &to, Context: ctx}
	var events []orderedEvent

	itAcc, err := s.filterer.FilterInterestAccrued(opts)
	if err != nil {
		return nil, fmt.Errorf("filter InterestAccrued: %w", err)
	}
	for itAcc.Next() {
		e := itAcc.Event
		events = append(events, orderedEvent{e.Raw.BlockNumber, e.Raw.Index, e.Raw.BlockHash,
			state.InterestAccrued{Elapsed: e.Elapsed, Interest: e.Interest}})
	}
	iterErr := itAcc.Error()
	_ = itAcc.Close()
	if iterErr != nil {
		return nil, fmt.Errorf("iterate InterestAccrued: %w", iterErr)
	}

	itSup, err := s.filterer.FilterSupplied(opts, nil)
	if err != nil {
		return nil, fmt.Errorf("filter Supplied: %w", err)
	}
	for itSup.Next() {
		e := itSup.Event
		events = append(events, orderedEvent{e.Raw.BlockNumber, e.Raw.Index, e.Raw.BlockHash,
			state.Supplied{User: addressKey(e.User), Assets: e.Assets, Shares: e.Shares}})
	}
	iterErr = itSup.Error()
	_ = itSup.Close()
	if iterErr != nil {
		return nil, fmt.Errorf("iterate Supplied: %w", iterErr)
	}

	itWd, err := s.filterer.FilterWithdrawn(opts, nil)
	if err != nil {
		return nil, fmt.Errorf("filter Withdrawn: %w", err)
	}
	for itWd.Next() {
		e := itWd.Event
		events = append(events, orderedEvent{e.Raw.BlockNumber, e.Raw.Index, e.Raw.BlockHash,
			state.Withdrawn{User: addressKey(e.User), Assets: e.Assets, Shares: e.Shares}})
	}
	iterErr = itWd.Error()
	_ = itWd.Close()
	if iterErr != nil {
		return nil, fmt.Errorf("iterate Withdrawn: %w", iterErr)
	}

	itCol, err := s.filterer.FilterCollateralSupplied(opts, nil)
	if err != nil {
		return nil, fmt.Errorf("filter CollateralSupplied: %w", err)
	}
	for itCol.Next() {
		e := itCol.Event
		events = append(events, orderedEvent{e.Raw.BlockNumber, e.Raw.Index, e.Raw.BlockHash,
			state.CollateralSupplied{User: addressKey(e.User), Assets: e.Assets}})
	}
	iterErr = itCol.Error()
	_ = itCol.Close()
	if iterErr != nil {
		return nil, fmt.Errorf("iterate CollateralSupplied: %w", iterErr)
	}

	itColWd, err := s.filterer.FilterCollateralWithdrawn(opts, nil)
	if err != nil {
		return nil, fmt.Errorf("filter CollateralWithdrawn: %w", err)
	}
	for itColWd.Next() {
		e := itColWd.Event
		events = append(events, orderedEvent{e.Raw.BlockNumber, e.Raw.Index, e.Raw.BlockHash,
			state.CollateralWithdrawn{User: addressKey(e.User), Assets: e.Assets}})
	}
	iterErr = itColWd.Error()
	_ = itColWd.Close()
	if iterErr != nil {
		return nil, fmt.Errorf("iterate CollateralWithdrawn: %w", iterErr)
	}

	itBorrowed, err := s.filterer.FilterBorrowed(opts, nil)
	if err != nil {
		return nil, fmt.Errorf("filter Borrowed: %w", err)
	}
	for itBorrowed.Next() {
		e := itBorrowed.Event
		events = append(events, orderedEvent{e.Raw.BlockNumber, e.Raw.Index, e.Raw.BlockHash,
			state.Borrowed{User: addressKey(e.User), Assets: e.Assets, Shares: e.Shares}})
	}
	iterErr = itBorrowed.Error()
	_ = itBorrowed.Close()
	if iterErr != nil {
		return nil, fmt.Errorf("iterate Borrowed: %w", iterErr)
	}

	itRepaid, err := s.filterer.FilterRepaid(opts, nil)
	if err != nil {
		return nil, fmt.Errorf("filter Repaid: %w", err)
	}
	for itRepaid.Next() {
		e := itRepaid.Event
		events = append(events, orderedEvent{e.Raw.BlockNumber, e.Raw.Index, e.Raw.BlockHash,
			state.Repaid{User: addressKey(e.User), Assets: e.Assets, Shares: e.Shares}})
	}
	iterErr = itRepaid.Error()
	_ = itRepaid.Close()
	if iterErr != nil {
		return nil, fmt.Errorf("iterate Repaid: %w", iterErr)
	}

	itLiquidated, err := s.filterer.FilterLiquidated(opts, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("filter Liquidated: %w", err)
	}
	for itLiquidated.Next() {
		e := itLiquidated.Event
		events = append(events, orderedEvent{e.Raw.BlockNumber, e.Raw.Index, e.Raw.BlockHash,
			state.Liquidated{
				Liquidator:       addressKey(e.Liquidator),
				Borrower:         addressKey(e.Borrower),
				RepaidAssets:     e.RepaidAssets,
				RepaidShares:     e.RepaidShares,
				SeizedCollateral: e.SeizedCollateral,
			}})
	}
	iterErr = itLiquidated.Error()
	_ = itLiquidated.Close()
	if iterErr != nil {
		return nil, fmt.Errorf("iterate Liquidated: %w", iterErr)
	}

	itBadDebt, err := s.filterer.FilterBadDebtRealized(opts, nil)
	if err != nil {
		return nil, fmt.Errorf("filter BadDebtRealized: %w", err)
	}
	for itBadDebt.Next() {
		e := itBadDebt.Event
		events = append(events, orderedEvent{e.Raw.BlockNumber, e.Raw.Index, e.Raw.BlockHash,
			state.BadDebtRealized{
				Borrower:      addressKey(e.Borrower),
				BadDebtAssets: e.BadDebtAssets,
				BadDebtShares: e.BadDebtShares,
			}})
	}
	iterErr = itBadDebt.Error()
	_ = itBadDebt.Close()
	if iterErr != nil {
		return nil, fmt.Errorf("iterate BadDebtRealized: %w", iterErr)
	}

	return events, nil
}
