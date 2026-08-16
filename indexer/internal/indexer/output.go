package indexer

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/synthlike/smolpho/indexer/internal/state"
	"github.com/synthlike/smolpho/indexer/internal/storage"
)

func addressKey(address common.Address) string {
	return strings.ToLower(address.Hex())
}

// printState renders a snapshot of the reconstructed market and positions.
func printState(ctx context.Context, output io.Writer, st storage.Store) error {
	snapshot, err := st.Snapshot(ctx)
	if err != nil {
		return err
	}
	s := snapshot.State
	checkpoint := snapshot.Checkpoint
	m := s.Market
	fmt.Fprintf(output, "\n=== market @ block %d (%s) ===\n", checkpoint.Number, checkpoint.Hash.Hex())
	fmt.Fprintf(output, "  totalSupplyAssets: %s\n", formatInt(m.TotalSupplyAssets))
	fmt.Fprintf(output, "  totalSupplyShares: %s\n", formatInt(m.TotalSupplyShares))
	fmt.Fprintf(output, "  totalBorrowAssets: %s\n", formatInt(m.TotalBorrowAssets))
	fmt.Fprintf(output, "  totalBorrowShares: %s\n", formatInt(m.TotalBorrowShares))
	fmt.Fprintf(output, "  lastUpdate: %s\n", formatInt(m.LastUpdate))
	fmt.Fprintf(output, "  supply share price: %s\n", supplySharePrice(m))
	if len(s.Positions) == 0 {
		fmt.Fprintln(output, "  (no positions)")
		return nil
	}
	fmt.Fprintln(output, "  positions:")
	users := make([]string, 0, len(s.Positions))
	for user := range s.Positions {
		users = append(users, user)
	}
	slices.Sort(users)
	for _, user := range users {
		p := s.Positions[user]
		fmt.Fprintf(output, "    %s  supplyShares=%s  supplyAssets=%s  collateral=%s\n",
			user,
			formatInt(p.SupplyShares),
			formatInt(s.SupplyAssets(user)),
			formatInt(p.Collateral),
		)
	}
	return nil
}

func formatInt(value *big.Int) string {
	digits := value.String()
	sign := ""
	if digits[0] == '-' {
		sign = "-"
		digits = digits[1:]
	}

	firstGroup := len(digits) % 3
	if firstGroup == 0 {
		firstGroup = 3
	}

	var formatted strings.Builder
	formatted.Grow(len(sign) + len(digits) + (len(digits)-1)/3)
	formatted.WriteString(sign)
	formatted.WriteString(digits[:firstGroup])
	for offset := firstGroup; offset < len(digits); offset += 3 {
		formatted.WriteByte('_')
		formatted.WriteString(digits[offset : offset+3])
	}
	return formatted.String()
}

func supplySharePrice(m state.Market) string {
	num := new(big.Float).SetInt(new(big.Int).Add(m.TotalSupplyAssets, state.VirtualAssets))
	den := new(big.Float).SetInt(new(big.Int).Add(m.TotalSupplyShares, state.VirtualShares))
	return new(big.Float).Quo(num, den).Text('g', 12)
}
