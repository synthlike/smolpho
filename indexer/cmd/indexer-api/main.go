// Command indexer-api durably reconstructs and serves a Smolpho deployment.
package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	smolphoindexer "github.com/synthlike/smolpho/indexer/internal/indexer"
)

func main() {
	var corsOrigins []string
	rpc := flag.String("rpc", "http://localhost:8545", "Ethereum RPC endpoint")
	contract := flag.String("contract", "", "Smolpho contract address (required)")
	deploymentBlock := flag.Uint64("deployment-block", 0, "Smolpho deployment block (required)")
	database := flag.String("database", "smolpho-indexer.sqlite", "SQLite database path")
	listen := flag.String("listen", "127.0.0.1:8080", "HTTP listen address")
	follow := flag.Bool("follow", true, "keep polling for new blocks after backfill")
	interval := flag.Duration("interval", 2*time.Second, "poll interval when following")
	batchSize := flag.Uint64("batch-size", 2_000, "maximum blocks per log query")
	flag.Func("cors-origin", "allowed browser origin; may be repeated", func(origin string) error {
		corsOrigins = append(corsOrigins, origin)
		return nil
	})
	flag.Parse()

	var deploymentBlockValue *uint64
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "deployment-block" {
			deploymentBlockValue = deploymentBlock
		}
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := run(ctx, apiConfig{
		Indexer: smolphoindexer.Config{
			RPC:             *rpc,
			Contract:        *contract,
			DeploymentBlock: deploymentBlockValue,
			Follow:          *follow,
			Interval:        *interval,
			BatchSize:       *batchSize,
		},
		Database:    *database,
		Listen:      *listen,
		CORSOrigins: corsOrigins,
	})
	if err != nil {
		log.Fatal(err)
	}
}
