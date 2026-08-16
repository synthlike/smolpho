// Package httpapi exposes published indexer state through a read-only HTTP API.
package httpapi

import (
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/synthlike/smolpho/indexer/internal/storage"
)

type api struct {
	store  storage.Store
	status *StatusTracker
}

// NewHandler constructs the complete read-only API handler.
func NewHandler(store storage.Store, status *StatusTracker) http.Handler {
	if status == nil {
		status = NewStatusTracker()
	}
	a := &api{store: store, status: status}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", a.health)
	mux.HandleFunc("GET /api/v1/status", a.indexerStatus)
	mux.HandleFunc("GET /api/v1/market", a.market)
	mux.HandleFunc("GET /api/v1/positions/{address}", a.position)
	return mux
}

type checkpointResponse struct {
	BlockNumber string `json:"blockNumber"`
	BlockHash   string `json:"blockHash"`
	Valid       bool   `json:"valid"`
}

type statusResponse struct {
	Ready                   bool               `json:"ready"`
	CaughtUp                bool               `json:"caughtUp"`
	Syncing                 bool               `json:"syncing"`
	Rebuilding              bool               `json:"rebuilding"`
	WaitingForDeployment    bool               `json:"waitingForDeployment"`
	ChainHead               *string            `json:"chainHead"`
	Lag                     *string            `json:"lag"`
	PublishedCheckpoint     checkpointResponse `json:"publishedCheckpoint"`
	WorkingCheckpoint       checkpointResponse `json:"workingCheckpoint"`
	LastSuccessfulSync      *string            `json:"lastSuccessfulSync"`
	LastError               string             `json:"lastError,omitempty"`
	DetectedReorganizations string             `json:"detectedReorganizations"`
}

type marketResponse struct {
	Checkpoint        checkpointResponse `json:"checkpoint"`
	TotalSupplyAssets string             `json:"totalSupplyAssets"`
	TotalSupplyShares string             `json:"totalSupplyShares"`
	TotalBorrowAssets string             `json:"totalBorrowAssets"`
	TotalBorrowShares string             `json:"totalBorrowShares"`
	LastUpdate        string             `json:"lastUpdate"`
}

type positionResponse struct {
	Checkpoint   checkpointResponse `json:"checkpoint"`
	Address      string             `json:"address"`
	SupplyShares string             `json:"supplyShares"`
	SupplyAssets string             `json:"supplyAssets"`
	BorrowShares string             `json:"borrowShares"`
	BorrowAssets string             `json:"borrowAssets"`
	Collateral   string             `json:"collateral"`
}

func (a *api) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *api) indexerStatus(w http.ResponseWriter, r *http.Request) {
	projection, err := a.store.ProjectionStatus(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "read projection status")
		return
	}
	runtime := a.status.Snapshot()
	ready := projection.Published.Valid
	caughtUp := ready && runtime.HeadKnown && !runtime.Syncing &&
		!projection.Rebuilding && runtime.LastError == "" &&
		projection.Published.Number >= runtime.Head

	var chainHead, lag, lastSuccessfulSync *string
	if runtime.HeadKnown {
		headValue := strconv.FormatUint(runtime.Head, 10)
		chainHead = &headValue
		if ready && !projection.Rebuilding && runtime.Head >= projection.Published.Number {
			lagValue := strconv.FormatUint(runtime.Head-projection.Published.Number, 10)
			lag = &lagValue
		}
	}
	if !runtime.LastSuccessfulSync.IsZero() {
		value := runtime.LastSuccessfulSync.Format(time.RFC3339Nano)
		lastSuccessfulSync = &value
	}
	writeJSON(w, http.StatusOK, statusResponse{
		Ready:                   ready,
		CaughtUp:                caughtUp,
		Syncing:                 runtime.Syncing,
		Rebuilding:              projection.Rebuilding,
		WaitingForDeployment:    runtime.Pending,
		ChainHead:               chainHead,
		Lag:                     lag,
		PublishedCheckpoint:     checkpointJSON(projection.Published),
		WorkingCheckpoint:       checkpointJSON(projection.Working),
		LastSuccessfulSync:      lastSuccessfulSync,
		LastError:               runtime.LastError,
		DetectedReorganizations: strconv.FormatUint(runtime.ReorgCount, 10),
	})
}

func (a *api) market(w http.ResponseWriter, r *http.Request) {
	snapshot, ok := a.publishedSnapshot(w, r)
	if !ok {
		return
	}
	market := snapshot.State.Market
	writeJSON(w, http.StatusOK, marketResponse{
		Checkpoint:        checkpointJSON(snapshot.Checkpoint),
		TotalSupplyAssets: decimal(market.TotalSupplyAssets),
		TotalSupplyShares: decimal(market.TotalSupplyShares),
		TotalBorrowAssets: decimal(market.TotalBorrowAssets),
		TotalBorrowShares: decimal(market.TotalBorrowShares),
		LastUpdate:        decimal(market.LastUpdate),
	})
}

func (a *api) position(w http.ResponseWriter, r *http.Request) {
	rawAddress := r.PathValue("address")
	if !common.IsHexAddress(rawAddress) {
		writeError(w, http.StatusBadRequest, "invalid Ethereum address")
		return
	}
	address := common.HexToAddress(rawAddress)
	key := strings.ToLower(address.Hex())
	snapshot, ok := a.publishedSnapshot(w, r)
	if !ok {
		return
	}
	position := snapshot.State.Positions[key]
	if position == nil {
		writeError(w, http.StatusNotFound, "position not found")
		return
	}
	writeJSON(w, http.StatusOK, positionResponse{
		Checkpoint:   checkpointJSON(snapshot.Checkpoint),
		Address:      address.Hex(),
		SupplyShares: decimal(position.SupplyShares),
		SupplyAssets: decimal(snapshot.State.SupplyAssets(key)),
		BorrowShares: decimal(position.BorrowShares),
		BorrowAssets: decimal(snapshot.State.BorrowAssets(key)),
		Collateral:   decimal(position.Collateral),
	})
}

func (a *api) publishedSnapshot(
	w http.ResponseWriter,
	r *http.Request,
) (storage.Snapshot, bool) {
	snapshot, err := a.store.Snapshot(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "read published state")
		return storage.Snapshot{}, false
	}
	if !snapshot.Checkpoint.Valid {
		writeError(w, http.StatusServiceUnavailable, "indexer has not published its initial backfill")
		return storage.Snapshot{}, false
	}
	return snapshot, true
}

func checkpointJSON(checkpoint storage.Checkpoint) checkpointResponse {
	return checkpointResponse{
		BlockNumber: strconv.FormatUint(checkpoint.Number, 10),
		BlockHash:   checkpoint.Hash.Hex(),
		Valid:       checkpoint.Valid,
	}
}

func decimal(value *big.Int) string {
	return value.String()
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
