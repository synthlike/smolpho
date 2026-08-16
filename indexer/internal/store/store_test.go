package store

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/synthlike/smolpho/indexer/internal/state"
)

func TestCommitAndReset(t *testing.T) {
	st := New(100)
	hash := common.HexToHash("0x1234")
	st.Commit([]state.Event{state.Supplied{
		User: "alice", Assets: big.NewInt(10), Shares: big.NewInt(20),
	}}, 7, hash)

	checkpoint := st.Checkpoint()
	if !checkpoint.Valid || checkpoint.Number != 7 || checkpoint.Hash != hash {
		t.Fatalf("checkpoint = %+v, want block 7 at %s", checkpoint, hash)
	}
	st.Read(func(s *state.State, _ Checkpoint) {
		if got := s.Market.TotalSupplyAssets; got.Cmp(big.NewInt(10)) != 0 {
			t.Fatalf("totalSupplyAssets = %s, want 10", got)
		}
	})

	st.Reset(200)
	if checkpoint := st.Checkpoint(); checkpoint.Valid {
		t.Fatalf("checkpoint remained valid after reset: %+v", checkpoint)
	}
	st.Read(func(s *state.State, _ Checkpoint) {
		if got := s.Market.TotalSupplyAssets; got.Sign() != 0 {
			t.Fatalf("totalSupplyAssets after reset = %s, want 0", got)
		}
		if got := s.Market.LastUpdate; got.Cmp(big.NewInt(200)) != 0 {
			t.Fatalf("lastUpdate after reset = %s, want 200", got)
		}
	})
}
