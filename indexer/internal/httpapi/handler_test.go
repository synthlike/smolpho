package httpapi

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	smolphoindexer "github.com/synthlike/smolpho/indexer/internal/indexer"
	"github.com/synthlike/smolpho/indexer/internal/state"
	"github.com/synthlike/smolpho/indexer/internal/storage"
	"github.com/synthlike/smolpho/indexer/internal/storage/memory"
)

const testAddress = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"

var testMarketConfig = MarketConfig{
	ChainID:              "31337",
	ContractAddress:      "0x0000000000000000000000000000000000000010",
	DeploymentBlock:      "4",
	LoanToken:            "0x0000000000000000000000000000000000000020",
	CollateralToken:      "0x0000000000000000000000000000000000000030",
	Oracle:               "0x0000000000000000000000000000000000000040",
	LLTV:                 "800000000000000000",
	RatePerSecond:        "1",
	LiquidationIncentive: "1100000000000000000",
}

func TestEndpoints(t *testing.T) {
	store := publishedStore(t)
	status := NewStatusTracker()
	status.Observe(smolphoindexer.SyncStatus{Head: 12, HeadKnown: true})
	handler := NewHandler(store, status, testMarketConfig)

	tests := []struct {
		path string
		want []string
	}{
		{path: "/healthz", want: []string{`"status":"ok"`}},
		{path: "/api/v1/config", want: []string{`"chainId":"31337"`, `"deploymentBlock":"4"`, `"lltv":"800000000000000000"`}},
		{path: "/api/v1/status", want: []string{`"caughtUp":true`, `"chainHead":"12"`}},
		{path: "/api/v1/market", want: []string{`"totalSupplyAssets":"50"`, `"totalBorrowAssets":"8"`}},
		{path: "/api/v1/positions/" + testAddress, want: []string{`"supplyAssets":"50"`, `"borrowAssets":"8"`}},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			response := request(handler, http.MethodGet, test.path)
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200: %s", response.Code, response.Body.String())
			}
			for _, want := range test.want {
				if !strings.Contains(response.Body.String(), want) {
					t.Fatalf("body does not contain %q: %s", want, response.Body.String())
				}
			}
		})
	}
}

func TestEndpointStatusCodes(t *testing.T) {
	ready := NewHandler(publishedStore(t), nil, testMarketConfig)
	notReady := NewHandler(memory.New(0), nil, testMarketConfig)
	tests := []struct {
		name    string
		handler http.Handler
		method  string
		path    string
		status  int
	}{
		{name: "initial backfill", handler: notReady, method: http.MethodGet, path: "/api/v1/market", status: http.StatusServiceUnavailable},
		{name: "invalid address", handler: ready, method: http.MethodGet, path: "/api/v1/positions/nope", status: http.StatusBadRequest},
		{name: "missing position", handler: ready, method: http.MethodGet, path: "/api/v1/positions/0x0000000000000000000000000000000000000001", status: http.StatusNotFound},
		{name: "wrong method", handler: ready, method: http.MethodPost, path: "/api/v1/market", status: http.StatusMethodNotAllowed},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := request(test.handler, test.method, test.path)
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d: %s", response.Code, test.status, response.Body.String())
			}
		})
	}
}

func publishedStore(t *testing.T) *memory.Store {
	t.Helper()
	store := memory.New(100)
	t.Cleanup(func() { _ = store.Close() })
	checkpoint := storage.Checkpoint{
		Number: 12, Hash: common.HexToHash("0xaaaa"), Valid: true,
	}
	if err := store.Commit(context.Background(), []state.Event{
		state.Supplied{User: testAddress, Assets: big.NewInt(50), Shares: big.NewInt(50_000_000)},
		state.CollateralSupplied{User: testAddress, Assets: big.NewInt(7)},
		state.Borrowed{User: testAddress, Assets: big.NewInt(8), Shares: big.NewInt(8_000_000)},
	}, checkpoint); err != nil {
		t.Fatal(err)
	}
	return store
}

func request(handler http.Handler, method, target string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}
