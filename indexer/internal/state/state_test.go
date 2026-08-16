package state

import (
	"math/big"
	"testing"
)

func bi(n int64) *big.Int { return big.NewInt(n) }

// Applying a first Supplied credits the position and grows the supply totals.
func TestSupplyUpdatesTotalsAndPosition(t *testing.T) {
	s := New(0)
	// A first 1000-asset supply into an empty market mints 1000 * VIRTUAL_SHARES
	// shares under the contract's toSharesDown; the reducer applies the emitted
	// pair verbatim.
	s.Apply(Supplied{User: "alice", Assets: bi(1000), Shares: bi(1_000_000_000)})

	if got := s.Market.TotalSupplyAssets; got.Cmp(bi(1000)) != 0 {
		t.Fatalf("totalSupplyAssets = %s, want 1000", got)
	}
	if got := s.Market.TotalSupplyShares; got.Cmp(bi(1_000_000_000)) != 0 {
		t.Fatalf("totalSupplyShares = %s, want 1e9", got)
	}
	if got := s.SupplyAssets("alice"); got.Cmp(bi(1000)) != 0 {
		t.Fatalf("supplyAssets(alice) = %s, want 1000", got)
	}
}

// Interest raises a supplier's asset value without changing any share balance.
func TestInterestGrowsAssetsNotShares(t *testing.T) {
	s := New(100)
	s.Apply(Supplied{User: "alice", Assets: bi(1000), Shares: bi(1_000_000_000)})

	sharesBefore := new(big.Int).Set(s.Positions["alice"].SupplyShares)
	s.Apply(InterestAccrued{Elapsed: bi(3600), Interest: bi(100), Timestamp: 3700})

	if got := s.Positions["alice"].SupplyShares; got.Cmp(sharesBefore) != 0 {
		t.Fatalf("shares changed on interest: %s -> %s", sharesBefore, got)
	}
	// Value rises by ~the accrued interest (rounded down via virtual offset):
	// floor(1e9 * (1100+1) / (1e9 + 1e6)) = 1099.
	if got := s.SupplyAssets("alice"); got.Cmp(bi(1099)) != 0 {
		t.Fatalf("supplyAssets(alice) = %s, want 1099", got)
	}
	if got := s.Market.LastUpdate; got.Cmp(bi(3700)) != 0 {
		t.Fatalf("lastUpdate = %s, want 3700", got)
	}
}

// Share conservation: per-user supply shares sum to the market total.
func TestShareConservation(t *testing.T) {
	s := New(0)
	s.Apply(Supplied{User: "alice", Assets: bi(1000), Shares: bi(1_000_000_000)})
	s.Apply(InterestAccrued{Elapsed: bi(3600), Interest: bi(100)})
	s.Apply(Supplied{User: "bob", Assets: bi(500), Shares: bi(450_000_000)})
	s.Apply(Withdrawn{User: "alice", Assets: bi(200), Shares: bi(180_000_000)})

	sum := new(big.Int)
	for _, p := range s.Positions {
		sum.Add(sum, p.SupplyShares)
	}
	if sum.Cmp(s.Market.TotalSupplyShares) != 0 {
		t.Fatalf("sum(shares) = %s, totalSupplyShares = %s", sum, s.Market.TotalSupplyShares)
	}
}

// Withdraw reduces the position and totals by the emitted amounts.
func TestWithdrawReducesTotals(t *testing.T) {
	s := New(0)
	s.Apply(Supplied{User: "alice", Assets: bi(1000), Shares: bi(1_000_000_000)})
	s.Apply(Withdrawn{User: "alice", Assets: bi(400), Shares: bi(400_000_000)})

	if got := s.Market.TotalSupplyAssets; got.Cmp(bi(600)) != 0 {
		t.Fatalf("totalSupplyAssets = %s, want 600", got)
	}
	if got := s.Positions["alice"].SupplyShares; got.Cmp(bi(600_000_000)) != 0 {
		t.Fatalf("alice shares = %s, want 6e8", got)
	}
}

// SupplyAssets for an unknown user is zero and does not panic.
func TestSupplyAssetsUnknownUser(t *testing.T) {
	s := New(0)
	if got := s.SupplyAssets("nobody"); got.Sign() != 0 {
		t.Fatalf("supplyAssets(nobody) = %s, want 0", got)
	}
}

func TestCollateralSupplyUpdatesPosition(t *testing.T) {
	s := New(0)
	s.Apply(CollateralSupplied{User: "alice", Assets: bi(250)})
	s.Apply(CollateralSupplied{User: "alice", Assets: bi(50)})

	if got := s.Positions["alice"].Collateral; got.Cmp(bi(300)) != 0 {
		t.Fatalf("alice collateral = %s, want 300", got)
	}
}

func TestCollateralWithdrawReducesPosition(t *testing.T) {
	s := New(0)
	s.Apply(CollateralSupplied{User: "alice", Assets: bi(300)})
	s.Apply(CollateralWithdrawn{User: "alice", Assets: bi(125)})

	if got := s.Positions["alice"].Collateral; got.Cmp(bi(175)) != 0 {
		t.Fatalf("alice collateral = %s, want 175", got)
	}
}

func TestBorrowAndRepayUpdateDebt(t *testing.T) {
	s := New(0)
	s.Apply(Borrowed{User: "alice", Assets: bi(1_000), Shares: bi(1_000_000_000)})

	if got := s.BorrowAssets("alice"); got.Cmp(bi(1_000)) != 0 {
		t.Fatalf("borrowAssets(alice) = %s, want 1000", got)
	}
	s.Apply(Repaid{User: "alice", Assets: bi(400), Shares: bi(400_000_000)})

	if got := s.Market.TotalBorrowAssets; got.Cmp(bi(600)) != 0 {
		t.Fatalf("totalBorrowAssets = %s, want 600", got)
	}
	if got := s.Market.TotalBorrowShares; got.Cmp(bi(600_000_000)) != 0 {
		t.Fatalf("totalBorrowShares = %s, want 600000000", got)
	}
	if got := s.Positions["alice"].BorrowShares; got.Cmp(bi(600_000_000)) != 0 {
		t.Fatalf("alice borrow shares = %s, want 600000000", got)
	}
}

func TestRepayAssetsSaturateAtZero(t *testing.T) {
	s := New(0)
	s.Apply(Borrowed{User: "alice", Assets: bi(1), Shares: bi(1)})
	s.Apply(Repaid{User: "alice", Assets: bi(2), Shares: bi(1)})

	if got := s.Market.TotalBorrowAssets; got.Sign() != 0 {
		t.Fatalf("totalBorrowAssets = %s, want 0", got)
	}
}

func TestLiquidationAndBadDebtMirrorContractPhases(t *testing.T) {
	s := New(0)
	s.Apply(Supplied{User: "lender", Assets: bi(5_000), Shares: bi(5_000_000_000)})
	s.Apply(CollateralSupplied{User: "borrower", Assets: bi(100)})
	s.Apply(Borrowed{User: "borrower", Assets: bi(3_000), Shares: bi(3_000_000_000)})
	s.Apply(Liquidated{
		Liquidator:       "liquidator",
		Borrower:         "borrower",
		RepaidAssets:     bi(100),
		RepaidShares:     bi(100_000_000),
		SeizedCollateral: bi(100),
	})

	position := s.Positions["borrower"]
	if got := position.BorrowShares; got.Cmp(bi(2_900_000_000)) != 0 {
		t.Fatalf("borrower shares after liquidation = %s, want 2900000000", got)
	}
	if got := position.Collateral; got.Sign() != 0 {
		t.Fatalf("borrower collateral after liquidation = %s, want 0", got)
	}
	if _, exists := s.Positions["liquidator"]; exists {
		t.Fatal("liquidation created a protocol position for the liquidator")
	}

	s.Apply(BadDebtRealized{
		Borrower:      "borrower",
		BadDebtAssets: bi(2_900),
		BadDebtShares: bi(2_900_000_000),
	})
	if got := position.BorrowShares; got.Sign() != 0 {
		t.Fatalf("borrower shares after bad debt = %s, want 0", got)
	}
	if got := s.Market.TotalBorrowAssets; got.Sign() != 0 {
		t.Fatalf("totalBorrowAssets after bad debt = %s, want 0", got)
	}
	if got := s.Market.TotalBorrowShares; got.Sign() != 0 {
		t.Fatalf("totalBorrowShares after bad debt = %s, want 0", got)
	}
	if got := s.Market.TotalSupplyAssets; got.Cmp(bi(2_100)) != 0 {
		t.Fatalf("totalSupplyAssets after bad debt = %s, want 2100", got)
	}
}

func TestBorrowShareConservation(t *testing.T) {
	s := New(0)
	s.Apply(Borrowed{User: "alice", Assets: bi(100), Shares: bi(200)})
	s.Apply(Borrowed{User: "bob", Assets: bi(300), Shares: bi(400)})
	s.Apply(Repaid{User: "alice", Assets: bi(25), Shares: bi(50)})
	s.Apply(Liquidated{
		Borrower: "bob", RepaidAssets: bi(75), RepaidShares: bi(100), SeizedCollateral: bi(0),
	})

	sum := new(big.Int)
	for _, position := range s.Positions {
		sum.Add(sum, position.BorrowShares)
	}
	if sum.Cmp(s.Market.TotalBorrowShares) != 0 {
		t.Fatalf("sum(borrow shares) = %s, totalBorrowShares = %s", sum, s.Market.TotalBorrowShares)
	}
}

func TestBorrowAssetsUnknownUser(t *testing.T) {
	s := New(0)
	if got := s.BorrowAssets("nobody"); got.Sign() != 0 {
		t.Fatalf("borrowAssets(nobody) = %s, want 0", got)
	}
}

func TestBorrowAssetsRoundsUp(t *testing.T) {
	s := New(0)
	s.Apply(Borrowed{User: "alice", Assets: bi(1), Shares: bi(1)})
	if got := s.BorrowAssets("alice"); got.Cmp(bi(1)) != 0 {
		t.Fatalf("borrowAssets(alice) = %s, want rounded-up 1", got)
	}
}
