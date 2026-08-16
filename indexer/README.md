# Smolpho indexer

The indexer reconstructs Smolpho's implemented supply-side market and user
positions from contract events. Its state and checkpoints are held in memory;
restarting the process replays the chain from the deployment block.

## Usage

```sh
go run ./cmd/indexer \
  -rpc http://localhost:8545 \
  -contract 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  -deployment-block 1 \
  -follow
```

`-contract` and `-deployment-block` are required. The deployment block is the
block containing the contract deployment, not an arbitrary query start. Event
replay requires the full history because the indexer has no persisted position
snapshot from which a partial replay could continue.

The indexer queries at most `-batch-size` blocks at a time (default `2000`). In
follow mode it polls every `-interval` (default `2s`). If the deployment block
is in the future, follow mode waits for it; a one-shot invocation reports an
error.

## Chain reorganizations

Every committed range records the canonical ending block hash. Before indexing
new blocks, the indexer compares that checkpoint with the node's canonical
header. If they differ, it discards its in-memory state and replays from the
deployment block. It also checks block hashes before and after each log query
and retries when the canonical chain changes during a query.
