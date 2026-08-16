# smolpho

Smolpho is a deliberately small lending protocol inspired by Morpho Blue. It
is an educational project for exploring lending-market accounting, Solidity
contract design, and event-driven indexing - all in small package.

## Components

### Solidity contracts

The [contracts](contracts/README.md) component contains one immutable isolated
market per `Smolpho` deployment. The current implementation includes:

- supply shares and asset withdrawals;
- direct collateral deposits and withdrawals;
- fixed-rate interest accrual;
- liquidity, preview, debt-value, and health views;
- bounded accounting values and reentrancy-safe token transfers.

The full intended protocol behavior and accounting rules are described in the
[Smolpho specification](SMOLPHO-SPEC.md).

### Indexer

The [indexer](indexer/README.md) reconstructs market and user state from
Smolpho events. It processes logs in canonical order, commits state in batches,
tracks block-hash checkpoints, and replays from the deployment block after a
chain reorganization. 
`indexer-cli` holds state in memory and prints it;
`indexer-api` persists the same projection in SQLite; it will be served via HTTP API.

### Project tooling

The [Taskfile](Taskfile.yml) builds the contracts and indexer, refreshes the Go
ABI bindings, runs contract tests, and provides a deterministic local demo.

## Prerequisites

- [Foundry](https://book.getfoundry.sh/) (`forge`, `cast`, and `anvil`)
- Go
- [Task](https://taskfile.dev/)
- `jq`
- `abigen`

## Build and test

Build the contracts, regenerate the ABI and Go bindings, and build the indexer:

```sh
task build
```

Run the Solidity and Go test suites:

```sh
task contracts:test
cd indexer && go test ./...
```

## Demo

The repository includes a demo that deploys the contracts,
follows them with the indexer, and sends a paced sequence of protocol actions.
Run it from the repository root in three terminals.

Terminal 1 — start a fresh local chain:

```sh
task demo:anvil
```

Terminal 2 — deploy the contracts and start the indexer:

```sh
task demo:deploy
task demo:indexer
```

Terminal 3 — run the scenario:

```sh
task demo:scenario
```

The scenario supplies loan assets, supplies collateral, and withdraws half of
the minted supply shares. Pauses between actions allow the indexer terminal to
print a separate reconstructed-state snapshot for each step.

The demo depends on Anvil's deterministic account zero, deployment nonces, and
contract addresses. `demo:deploy` verifies that the account nonce is zero and
refuses to deploy onto a reused chain. To run the demo again, stop Anvil, start
a fresh `task demo:anvil`, and deploy again. The private key in the Taskfile is
Anvil's public test-only key and must never be used on a real network.
