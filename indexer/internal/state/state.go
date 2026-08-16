// Package state reconstructs Smolpho market and position accounting from the
// contract's emitted events. It is intentionally free of any chain/RPC
// dependency: it consumes plain decoded event values and applies them in log
// order, mirroring the on-chain storage transitions exactly. This makes the
// reducer directly unit-testable against the spec's deterministic scenarios.
package state

import "math/big"

// Virtual offsets mirror Smolpho's constants. They define the initial exchange
// rate and keep empty-pool conversions valid.
var (
	VirtualAssets = big.NewInt(1)
	VirtualShares = big.NewInt(1e6)
)

// Market mirrors the on-chain Market struct.
type Market struct {
	TotalSupplyAssets *big.Int
	TotalSupplyShares *big.Int
	TotalBorrowAssets *big.Int
	TotalBorrowShares *big.Int
	LastUpdate        *big.Int
}

// Position mirrors the on-chain Position struct.
type Position struct {
	SupplyShares *big.Int
	BorrowShares *big.Int
	Collateral   *big.Int
}

// State is the reconstructed market plus every known user position, keyed by
// lowercased hex address.
type State struct {
	Market    Market
	Positions map[string]*Position
}

// New returns a zeroed state, ready to have events applied in log order. The
// initial timestamp is the deployment block timestamp, which is also the
// contract's initial market.lastUpdate value.
func New(initialTimestamp uint64) *State {
	return &State{
		Market: Market{
			TotalSupplyAssets: new(big.Int),
			TotalSupplyShares: new(big.Int),
			TotalBorrowAssets: new(big.Int),
			TotalBorrowShares: new(big.Int),
			LastUpdate:        new(big.Int).SetUint64(initialTimestamp),
		},
		Positions: make(map[string]*Position),
	}
}

// Clone returns a deep copy that can be read or mutated independently.
func (s *State) Clone() *State {
	clone := &State{
		Market: Market{
			TotalSupplyAssets: cloneInt(s.Market.TotalSupplyAssets),
			TotalSupplyShares: cloneInt(s.Market.TotalSupplyShares),
			TotalBorrowAssets: cloneInt(s.Market.TotalBorrowAssets),
			TotalBorrowShares: cloneInt(s.Market.TotalBorrowShares),
			LastUpdate:        cloneInt(s.Market.LastUpdate),
		},
		Positions: make(map[string]*Position, len(s.Positions)),
	}
	for user, position := range s.Positions {
		clone.Positions[user] = &Position{
			SupplyShares: cloneInt(position.SupplyShares),
			BorrowShares: cloneInt(position.BorrowShares),
			Collateral:   cloneInt(position.Collateral),
		}
	}
	return clone
}

func cloneInt(value *big.Int) *big.Int {
	if value == nil {
		return nil
	}
	return new(big.Int).Set(value)
}

// Event is a sealed set of the events the reducer understands. Concrete types
// carry already-decoded values; the caller is responsible for delivering them
// in (blockNumber, logIndex) order, matching on-chain emission order.
type Event interface{ isEvent() }

// InterestAccrued adds interest to both the borrow and supply asset totals
// without changing any shares — the mechanism by which shares appreciate.
type InterestAccrued struct {
	Elapsed   *big.Int
	Interest  *big.Int
	Timestamp uint64
}

// Supplied credits supply shares to a user and grows the supply totals.
type Supplied struct {
	User   string
	Assets *big.Int
	Shares *big.Int
}

// Withdrawn burns supply shares from a user and shrinks the supply totals.
type Withdrawn struct {
	User   string
	Assets *big.Int
	Shares *big.Int
}

// CollateralSupplied credits collateral assets directly to a user position.
type CollateralSupplied struct {
	User   string
	Assets *big.Int
}

// CollateralWithdrawn removes collateral assets from a user position.
type CollateralWithdrawn struct {
	User   string
	Assets *big.Int
}

// Borrowed credits debt shares to a user and grows the borrow totals.
type Borrowed struct {
	User   string
	Assets *big.Int
	Shares *big.Int
}

// Repaid burns a user's debt shares and shrinks the borrow totals.
type Repaid struct {
	User   string
	Assets *big.Int
	Shares *big.Int
}

// Liquidated repays part of a borrower's debt and removes seized collateral.
// Liquidator is retained for event fidelity but has no protocol position
// transition: seized collateral leaves the protocol rather than becoming
// deposited collateral for the liquidator.
type Liquidated struct {
	Liquidator       string
	Borrower         string
	RepaidAssets     *big.Int
	RepaidShares     *big.Int
	SeizedCollateral *big.Int
}

// BadDebtRealized removes a borrower's remaining unpayable debt and charges
// the corresponding asset loss to suppliers.
type BadDebtRealized struct {
	Borrower      string
	BadDebtAssets *big.Int
	BadDebtShares *big.Int
}

func (InterestAccrued) isEvent()     {}
func (Supplied) isEvent()            {}
func (Withdrawn) isEvent()           {}
func (CollateralSupplied) isEvent()  {}
func (CollateralWithdrawn) isEvent() {}
func (Borrowed) isEvent()            {}
func (Repaid) isEvent()              {}
func (Liquidated) isEvent()          {}
func (BadDebtRealized) isEvent()     {}

// position returns the user's position, creating a zeroed one on first sight.
func (s *State) position(user string) *Position {
	p, ok := s.Positions[user]
	if !ok {
		p = &Position{
			SupplyShares: new(big.Int),
			BorrowShares: new(big.Int),
			Collateral:   new(big.Int),
		}
		s.Positions[user] = p
	}
	return p
}

// Apply advances the state by one event. Applying events in log order
// reproduces on-chain state exactly, including the InterestAccrued that
// precedes every Supplied/Withdrawn (state-changing calls accrue first).
func (s *State) Apply(ev Event) {
	switch e := ev.(type) {
	case InterestAccrued:
		s.Market.TotalBorrowAssets.Add(s.Market.TotalBorrowAssets, e.Interest)
		s.Market.TotalSupplyAssets.Add(s.Market.TotalSupplyAssets, e.Interest)
		s.Market.LastUpdate.SetUint64(e.Timestamp)
	case Supplied:
		p := s.position(e.User)
		p.SupplyShares.Add(p.SupplyShares, e.Shares)
		s.Market.TotalSupplyShares.Add(s.Market.TotalSupplyShares, e.Shares)
		s.Market.TotalSupplyAssets.Add(s.Market.TotalSupplyAssets, e.Assets)
	case Withdrawn:
		p := s.position(e.User)
		p.SupplyShares.Sub(p.SupplyShares, e.Shares)
		s.Market.TotalSupplyShares.Sub(s.Market.TotalSupplyShares, e.Shares)
		s.Market.TotalSupplyAssets.Sub(s.Market.TotalSupplyAssets, e.Assets)
	case CollateralSupplied:
		p := s.position(e.User)
		p.Collateral.Add(p.Collateral, e.Assets)
	case CollateralWithdrawn:
		p := s.position(e.User)
		p.Collateral.Sub(p.Collateral, e.Assets)
	case Borrowed:
		p := s.position(e.User)
		p.BorrowShares.Add(p.BorrowShares, e.Shares)
		s.Market.TotalBorrowShares.Add(s.Market.TotalBorrowShares, e.Shares)
		s.Market.TotalBorrowAssets.Add(s.Market.TotalBorrowAssets, e.Assets)
	case Repaid:
		p := s.position(e.User)
		p.BorrowShares.Sub(p.BorrowShares, e.Shares)
		s.Market.TotalBorrowShares.Sub(s.Market.TotalBorrowShares, e.Shares)
		subFloorZero(s.Market.TotalBorrowAssets, e.Assets)
	case Liquidated:
		p := s.position(e.Borrower)
		p.BorrowShares.Sub(p.BorrowShares, e.RepaidShares)
		p.Collateral.Sub(p.Collateral, e.SeizedCollateral)
		s.Market.TotalBorrowShares.Sub(s.Market.TotalBorrowShares, e.RepaidShares)
		subFloorZero(s.Market.TotalBorrowAssets, e.RepaidAssets)
	case BadDebtRealized:
		p := s.position(e.Borrower)
		p.BorrowShares.Sub(p.BorrowShares, e.BadDebtShares)
		s.Market.TotalBorrowShares.Sub(s.Market.TotalBorrowShares, e.BadDebtShares)
		subFloorZero(s.Market.TotalBorrowAssets, e.BadDebtAssets)
		subFloorZero(s.Market.TotalSupplyAssets, e.BadDebtAssets)
	}
}

func subFloorZero(value, amount *big.Int) {
	if amount.Cmp(value) >= 0 {
		value.SetInt64(0)
		return
	}
	value.Sub(value, amount)
}

// SupplyAssets returns the loan-token value of a user's supply shares, using
// the same rounded-down virtual-offset conversion as the contract's
// supplyAssets view: floor(shares * (assets+VA) / (shares_total+VS)).
func (s *State) SupplyAssets(user string) *big.Int {
	p, ok := s.Positions[user]
	if !ok {
		return new(big.Int)
	}
	return toAssetsDown(p.SupplyShares, s.Market.TotalSupplyAssets, s.Market.TotalSupplyShares)
}

// BorrowAssets returns the loan-token debt represented by a user's borrow
// shares, using the contract's rounded-up virtual-offset conversion.
func (s *State) BorrowAssets(user string) *big.Int {
	p, ok := s.Positions[user]
	if !ok {
		return new(big.Int)
	}
	return toAssetsUp(p.BorrowShares, s.Market.TotalBorrowAssets, s.Market.TotalBorrowShares)
}

// toAssetsDown = floor(shares * (totalAssets + VIRTUAL_ASSETS) /
//
//	(totalShares + VIRTUAL_SHARES)).
func toAssetsDown(shares, totalAssets, totalShares *big.Int) *big.Int {
	num := new(big.Int).Mul(shares, new(big.Int).Add(totalAssets, VirtualAssets))
	den := new(big.Int).Add(totalShares, VirtualShares)
	return num.Div(num, den)
}

func toAssetsUp(shares, totalAssets, totalShares *big.Int) *big.Int {
	num := new(big.Int).Mul(shares, new(big.Int).Add(totalAssets, VirtualAssets))
	den := new(big.Int).Add(totalShares, VirtualShares)
	quotient, remainder := new(big.Int), new(big.Int)
	quotient.QuoRem(num, den, remainder)
	if remainder.Sign() != 0 {
		quotient.Add(quotient, big.NewInt(1))
	}
	return quotient
}
