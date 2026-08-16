# Smolpho indexer

The indexer reconstructs Smolpho's market and user positions from contract
events, including supply, collateral, borrowing, repayment, liquidation, and
bad-debt accounting. The shared indexing engine is used by two applications:

- `indexer-cli` keeps state in memory and prints snapshots for demos and manual
  inspection;
- `indexer-api` persists state and checkpoints in SQLite and exposes the
  published projection through a read-only HTTP API.

See the [indexer design notes](../INDEXING.md) for its design history, event
pipeline, reorganization handling, atomicity model, and storage trade-offs.

## CLI usage

```sh
go run ./cmd/indexer-cli \
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

## API worker usage

```sh
go run ./cmd/indexer-api \
  -rpc http://localhost:8545 \
  -contract 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  -deployment-block 1 \
  -database smolpho-indexer.sqlite \
  -listen 127.0.0.1:8080
```

The API worker follows by default and resumes from its durable SQLite
checkpoint after a restart. Pass `-follow=false` for a one-shot backfill. Use a
separate database file for each chain and Smolpho deployment; deployment
identity metadata is not stored in the schema yet.

For the deterministic Anvil deployment, run it from the repository root with:

```sh
task demo:api
```

The SQLite schema is intentionally allowed to break while this is a prototype.
There are no migrations: after an incompatible schema change, delete the local
database and let the worker reconstruct it from chain events.

Initial backfills and reorg replays are built in a hidden SQLite generation.
The last complete generation remains available to readers until replay reaches
the sampled canonical head, when the worker atomically publishes the rebuilt
state. Hidden progress also survives a process restart.

## HTTP API

The API serves:

- `GET /healthz` — process liveness;
- `GET /api/v1/status` — chain head, lag, sync/rebuild flags, published and
  working checkpoints, last error, and last successful synchronization;
- `GET /api/v1/market` — published market totals; and
- `GET /api/v1/positions/{address}` — published shares, derived asset values,
  and collateral for one Ethereum address.

For example:

```sh
curl http://127.0.0.1:8080/api/v1/status | jq
curl http://127.0.0.1:8080/api/v1/market | jq
curl http://127.0.0.1:8080/api/v1/positions/0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266 | jq
```

All block numbers and accounting values are decimal JSON strings. This keeps
Solidity-sized integers exact for JavaScript and TypeScript clients. State
endpoints return `503` until the first complete projection is published.
During a reorg they continue serving the previous complete generation while
`/api/v1/status` reports `rebuilding: true` and exposes hidden working progress.

## Chain reorganizations

Every committed range records the canonical ending block hash. Before indexing
new blocks, the indexer compares that checkpoint with the node's canonical
header. If they differ, it rebuilds reconstructed state from the deployment
block in a hidden generation and publishes it atomically when complete. It also
checks block hashes before and after each log query and retries when the
canonical chain changes during a query.
