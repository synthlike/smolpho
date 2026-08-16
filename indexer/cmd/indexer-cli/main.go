// Command indexer-cli reconstructs a Smolpho deployment in memory and prints
// snapshots to the terminal.
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
	rpc := flag.String("rpc", "http://localhost:8545", "Ethereum RPC endpoint")
	contract := flag.String("contract", "", "Smolpho contract address (required)")
	deploymentBlock := flag.Uint64("deployment-block", 0, "Smolpho deployment block (required)")
	follow := flag.Bool("follow", false, "keep polling for new blocks after backfill")
	interval := flag.Duration("interval", 2*time.Second, "poll interval when following")
	batchSize := flag.Uint64("batch-size", 2_000, "maximum blocks per log query")
	flag.Parse()

	var deploymentBlockValue *uint64
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "deployment-block" {
			deploymentBlockValue = deploymentBlock
		}
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := run(ctx, smolphoindexer.Config{
		RPC:             *rpc,
		Contract:        *contract,
		DeploymentBlock: deploymentBlockValue,
		Follow:          *follow,
		Interval:        *interval,
		BatchSize:       *batchSize,
	})
	if err != nil {
		log.Fatal(err)
	}
}
