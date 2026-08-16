package indexer

import (
	"fmt"
	"io"
	"math/big"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/synthlike/smolpho/indexer/internal/state"
	"github.com/synthlike/smolpho/indexer/internal/store"
)

func addressKey(address common.Address) string {
	return strings.ToLower(address.Hex())
}

// printState renders a snapshot of the reconstructed market and positions.
func printState(output io.Writer, st *store.Store) {
	st.Read(func(s *state.State, checkpoint store.Checkpoint) {
		m := s.Market
		fmt.Fprintf(output, "\n=== market @ block %d (%s) ===\n", checkpoint.Number, checkpoint.Hash.Hex())
		fmt.Fprintf(output, "  totalSupplyAssets: %s\n", m.TotalSupplyAssets)
		fmt.Fprintf(output, "  totalSupplyShares: %s\n", m.TotalSupplyShares)
		fmt.Fprintf(output, "  totalBorrowAssets: %s\n", m.TotalBorrowAssets)
		fmt.Fprintf(output, "  totalBorrowShares: %s\n", m.TotalBorrowShares)
		fmt.Fprintf(output, "  lastUpdate: %s\n", m.LastUpdate)
		fmt.Fprintf(output, "  supply share price: %s\n", supplySharePrice(m))
		if len(s.Positions) == 0 {
			fmt.Fprintln(output, "  (no positions)")
			return
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
				user, p.SupplyShares, s.SupplyAssets(user), p.Collateral)
		}
	})
}

func supplySharePrice(m state.Market) string {
	num := new(big.Float).SetInt(new(big.Int).Add(m.TotalSupplyAssets, state.VirtualAssets))
	den := new(big.Float).SetInt(new(big.Int).Add(m.TotalSupplyShares, state.VirtualShares))
	return new(big.Float).Quo(num, den).Text('g', 12)
}
