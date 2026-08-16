package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/synthlike/smolpho/indexer/internal/bindings"
	"github.com/synthlike/smolpho/indexer/internal/httpapi"
	smolphoindexer "github.com/synthlike/smolpho/indexer/internal/indexer"
	"github.com/synthlike/smolpho/indexer/internal/storage"
	"github.com/synthlike/smolpho/indexer/internal/storage/sqlite"
)

type apiConfig struct {
	Indexer     smolphoindexer.Config
	Database    string
	Listen      string
	CORSOrigins []string
}

func run(ctx context.Context, config apiConfig) error {
	if err := config.Indexer.Validate(); err != nil {
		return err
	}
	if config.Listen == "" {
		return fmt.Errorf("-listen is required")
	}
	if err := httpapi.ValidateAllowedOrigins(config.CORSOrigins); err != nil {
		return err
	}
	marketConfig, err := loadMarketConfig(ctx, config.Indexer)
	if err != nil {
		return err
	}
	store, err := sqlite.Open(ctx, config.Database, 0)
	if err != nil {
		return fmt.Errorf("open indexer database: %w", err)
	}
	defer store.Close()

	listener, err := net.Listen("tcp", config.Listen)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", config.Listen, err)
	}
	defer listener.Close()

	runContext, cancel := context.WithCancel(ctx)
	defer cancel()
	logger := log.Default()
	status := httpapi.NewStatusTracker()
	handler := httpapi.NewHandler(store, status, marketConfig)
	server := &http.Server{
		Handler:           httpapi.WithCORS(handler, config.CORSOrigins),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return runContext
		},
	}

	workerErrors := make(chan error, 1)
	go func() {
		status.Observe(smolphoindexer.SyncStatus{Syncing: true})
		workerErrors <- smolphoindexer.Run(
			runContext,
			config.Indexer,
			smolphoindexer.Dependencies{
				Store:            store,
				Logger:           logger,
				HandleSyncStatus: status.Observe,
				HandleSnapshot: func(_ context.Context, snapshot storage.Snapshot) error {
					logger.Printf("indexed canonical block %d (%s)",
						snapshot.Checkpoint.Number, snapshot.Checkpoint.Hash.Hex())
					return nil
				},
			},
		)
	}()
	serverErrors := make(chan error, 1)
	go func() {
		logger.Printf("serving indexer API on http://%s", listener.Addr())
		serverErrors <- server.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		cancel()
		shutdownErr := shutdownServer(server)
		<-workerErrors
		return shutdownErr
	case workerErr := <-workerErrors:
		cancel()
		shutdownErr := shutdownServer(server)
		if ctx.Err() != nil {
			return shutdownErr
		}
		if workerErr != nil {
			return workerErr
		}
		return shutdownErr
	case serverErr := <-serverErrors:
		cancel()
		workerErr := <-workerErrors
		if serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
			return fmt.Errorf("serve indexer API: %w", serverErr)
		}
		return workerErr
	}
}

func loadMarketConfig(
	ctx context.Context,
	config smolphoindexer.Config,
) (httpapi.MarketConfig, error) {
	client, err := ethclient.DialContext(ctx, config.RPC)
	if err != nil {
		return httpapi.MarketConfig{}, fmt.Errorf("dial %s for market config: %w", config.RPC, err)
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return httpapi.MarketConfig{}, fmt.Errorf("read chain ID: %w", err)
	}
	contractAddress := common.HexToAddress(config.Contract)
	contract, err := bindings.NewSmolphoCaller(contractAddress, client)
	if err != nil {
		return httpapi.MarketConfig{}, fmt.Errorf("bind contract for market config: %w", err)
	}
	call := &bind.CallOpts{Context: ctx}
	loanToken, err := contract.LoanToken(call)
	if err != nil {
		return httpapi.MarketConfig{}, fmt.Errorf("read loan token: %w", err)
	}
	collateralToken, err := contract.CollateralToken(call)
	if err != nil {
		return httpapi.MarketConfig{}, fmt.Errorf("read collateral token: %w", err)
	}
	oracle, err := contract.Oracle(call)
	if err != nil {
		return httpapi.MarketConfig{}, fmt.Errorf("read oracle: %w", err)
	}
	lltv, err := contract.Lltv(call)
	if err != nil {
		return httpapi.MarketConfig{}, fmt.Errorf("read LLTV: %w", err)
	}
	ratePerSecond, err := contract.RatePerSecond(call)
	if err != nil {
		return httpapi.MarketConfig{}, fmt.Errorf("read rate per second: %w", err)
	}
	liquidationIncentive, err := contract.LiquidationIncentive(call)
	if err != nil {
		return httpapi.MarketConfig{}, fmt.Errorf("read liquidation incentive: %w", err)
	}

	return httpapi.MarketConfig{
		ChainID:              chainID.String(),
		ContractAddress:      contractAddress.Hex(),
		DeploymentBlock:      strconv.FormatUint(*config.DeploymentBlock, 10),
		LoanToken:            loanToken.Hex(),
		CollateralToken:      collateralToken.Hex(),
		Oracle:               oracle.Hex(),
		LLTV:                 lltv.String(),
		RatePerSecond:        strconv.FormatUint(ratePerSecond, 10),
		LiquidationIncentive: liquidationIncentive.String(),
	}, nil
}

func shutdownServer(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("shut down indexer API: %w", err)
	}
	return nil
}
