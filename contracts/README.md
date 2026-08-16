# Smolpho contracts

Smolpho is a minimal, educational lending protocol inspired by Morpho Blue. One `Smolpho` deployment represents one immutable isolated market with a loan token, collateral token, oracle, LLTV, fixed interest rate, and liquidation incentive.

It is not intended for production use.

## Features

- supply and withdrawal using supply shares;
- direct collateral supply and withdrawal;
- oracle-based borrower health checks;
- collateralized borrowing and share-based repayment;
- lazy fixed-rate interest accrual;
- incentivized liquidation;
- bad-debt socialization across suppliers;
- virtual assets and shares with conservative rounding;
- bounded accounting types, safe token transfers, and reentrancy protection.

## Layout

```text
src/Smolpho.sol                  Main isolated lending market
src/interfaces/                  Token and oracle interfaces
src/libraries/SharesMath.sol     Asset/share conversions
src/libraries/SafeTransferLib.sol Safe ERC-20 calls
test/                            Deterministic, fuzz, and mock-token tests
```

## Commands

Run from this directory:

```sh
forge build --skip test
forge test
forge fmt --check
```

Or run the repository-level tasks:

```sh
task contracts:build
task contracts:test
task build
```

## TODO - tests

- [ ] Verify direct loan-token donations do not change accounting exchange rates.
- [ ] Demonstrate that fee-on-transfer token accounting is explicitly unsupported.
- [ ] Invariant: user supply shares sum to total supply shares.
- [ ] Invariant: user borrow shares sum to total borrow shares.
- [ ] Invariant: total borrow assets never exceed total supply assets.
- [ ] Invariant: successful borrowing and collateral withdrawal leave positions healthy.
- [ ] Fuzz that liquidation never seizes more than borrower collateral.
- [ ] Fuzz that bad-debt realization never increases supplier claims.
- [ ] Add a stateful handler covering supply, withdrawal, collateral, borrowing, repayment, price changes, liquidation, and bad debt.
