// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {Smolpho} from "../src/Smolpho.sol";
import {IERC20} from "../src/interfaces/IERC20.sol";
import {SafeTransferLib} from "../src/libraries/SafeTransferLib.sol";
import {SharesMath} from "../src/libraries/SharesMath.sol";
import {FalseReturnERC20, MockERC20} from "./mocks/MockERC20.sol";
import {MockOracle} from "./mocks/MockOracle.sol";

contract SmolphoRepayHarness is Smolpho {
    constructor(IERC20 loanToken_, IERC20 collateralToken_, MockOracle oracle_, uint256 ratePerSecond_)
        Smolpho(loanToken_, collateralToken_, oracle_, 0.8e18, ratePerSecond_, 1.05e18)
    {}

    function setRawDebt(address user, uint256 totalBorrowAssets, uint256 totalBorrowShares, uint256 userBorrowShares)
        external
    {
        market.totalSupplyAssets = uint128(totalBorrowAssets);
        market.totalBorrowAssets = uint128(totalBorrowAssets);
        market.totalBorrowShares = uint128(totalBorrowShares);
        position[user].borrowShares = uint128(userBorrowShares);
    }
}

contract ReentrantRepayERC20 is MockERC20 {
    Smolpho public target;
    bool public attemptedReentry;
    bool public reentrySucceeded;
    bytes4 public reentryError;

    constructor() MockERC20("Reentrant Repay", "REENTER") {}

    function setTarget(Smolpho target_) external {
        target = target_;
    }

    function transferFrom(address from, address to, uint256 value) public override returns (bool) {
        if (address(target) != address(0) && !attemptedReentry) {
            attemptedReentry = true;
            try target.repay(1) {
                reentrySucceeded = true;
            } catch (bytes memory reason) {
                if (reason.length >= 4) {
                    bytes4 selector;
                    assembly {
                        selector := mload(add(reason, 32))
                    }
                    reentryError = selector;
                }
            }
        }

        return super.transferFrom(from, to, value);
    }
}

contract RepayTest is Test {
    uint256 internal constant START_TIME = 1_000_000;
    uint256 internal constant SUPPLY_ASSETS = 500e18;
    uint256 internal constant BORROW_ASSETS = 100e18;
    uint256 internal constant COLLATERAL_ASSETS = 2e18;

    address internal constant ALICE = address(0xA11CE);
    address internal constant BOB = address(0xB0B);

    MockERC20 internal loanToken;
    MockERC20 internal collateralToken;
    MockOracle internal oracle;
    SmolphoRepayHarness internal smolpho;

    function setUp() public {
        vm.warp(START_TIME);
        loanToken = new MockERC20("Loan Token", "LOAN");
        collateralToken = new MockERC20("Collateral Token", "COLLATERAL");
        oracle = new MockOracle(2_000e18);
        smolpho = _deploy(loanToken, 0);
        _seedDebt(smolpho);
    }

    function test_FullShareRepaymentClearsDebt() public {
        (, uint256 shares,) = smolpho.position(ALICE);

        vm.expectEmit(address(smolpho));
        emit Smolpho.InterestAccrued(0, 0);
        vm.expectEmit(address(smolpho));
        emit Smolpho.Repaid(ALICE, BORROW_ASSETS, shares);

        vm.prank(ALICE);
        uint256 assets = smolpho.repay(shares);

        assertEq(assets, BORROW_ASSETS);
        assertEq(smolpho.borrowAssets(ALICE), 0);
        assertEq(loanToken.balanceOf(ALICE), 0);
        assertEq(loanToken.balanceOf(address(smolpho)), SUPPLY_ASSETS);

        (uint256 supplyAssets,, uint256 borrowAssets, uint256 borrowShares,) = smolpho.market();
        (, uint256 userBorrowShares, uint256 collateral) = smolpho.position(ALICE);
        assertEq(supplyAssets, SUPPLY_ASSETS);
        assertEq(borrowAssets, 0);
        assertEq(borrowShares, 0);
        assertEq(userBorrowShares, 0);
        assertEq(collateral, COLLATERAL_ASSETS);
    }

    function test_PartialRepaymentReducesDebtAndRestoresLiquidity() public {
        (, uint256 totalShares,) = smolpho.position(ALICE);
        uint256 repaidShares = totalShares / 2;

        vm.prank(ALICE);
        uint256 assets = smolpho.repay(repaidShares);

        assertEq(assets, 50e18);
        assertEq(smolpho.borrowAssets(ALICE), 50e18);
        assertEq(smolpho.availableLiquidity(), 450e18);
        (,, uint256 borrowAssets, uint256 borrowShares,) = smolpho.market();
        assertEq(borrowAssets, 50e18);
        assertEq(borrowShares, totalShares - repaidShares);
    }

    function test_AccruesInterestBeforeCalculatingRepaymentAssets() public {
        SmolphoRepayHarness interestMarket = _deploy(loanToken, 1e15);
        _seedDebt(interestMarket);
        (, uint256 shares,) = interestMarket.position(ALICE);
        vm.warp(START_TIME + 10);
        uint256 expectedAssets = SharesMath.toAssetsUp(shares, 101e18, shares);
        loanToken.mint(ALICE, expectedAssets - BORROW_ASSETS);

        vm.prank(ALICE);
        uint256 assets = interestMarket.repay(shares);

        assertEq(assets, expectedAssets);
        assertEq(interestMarket.borrowAssets(ALICE), 0);
        (,, uint256 borrowAssets, uint256 borrowShares, uint256 lastUpdate) = interestMarket.market();
        assertEq(borrowAssets, 0);
        assertEq(borrowShares, 0);
        assertEq(lastUpdate, START_TIME + 10);
    }

    function test_RepaymentRoundsAssetsUp() public {
        SmolphoRepayHarness roundingMarket = _deploy(loanToken, 0);
        roundingMarket.setRawDebt(ALICE, 1, 1, 1);
        loanToken.mint(ALICE, 1);
        vm.prank(ALICE);
        loanToken.approve(address(roundingMarket), 1);

        vm.prank(ALICE);
        uint256 assets = roundingMarket.repay(1);

        assertEq(assets, 1);
        (,, uint256 borrowAssets, uint256 borrowShares,) = roundingMarket.market();
        assertEq(borrowAssets, 0);
        assertEq(borrowShares, 0);
    }

    function test_DoesNotCallOracle() public {
        oracle.setShouldRevert(true);
        (, uint256 shares,) = smolpho.position(ALICE);

        vm.prank(ALICE);
        smolpho.repay(shares);

        assertEq(smolpho.borrowAssets(ALICE), 0);
    }

    function test_RevertsForZeroShares() public {
        vm.expectRevert(Smolpho.ZeroShares.selector);
        vm.prank(ALICE);
        smolpho.repay(0);
    }

    function test_RevertsForInsufficientBorrowShares() public {
        (, uint256 shares,) = smolpho.position(ALICE);

        vm.expectRevert(Smolpho.InsufficientBorrowShares.selector);
        vm.prank(ALICE);
        smolpho.repay(shares + 1);
    }

    function test_TransferFailureRollsBackDebtAccounting() public {
        FalseReturnERC20 falseToken = new FalseReturnERC20();
        SmolphoRepayHarness falseMarket = _deploy(falseToken, 0);
        falseMarket.setRawDebt(ALICE, 10e18, 10e24, 10e24);
        falseToken.mint(ALICE, 10e18);
        vm.prank(ALICE);
        falseToken.approve(address(falseMarket), 10e18);

        vm.expectRevert(SafeTransferLib.TransferFailed.selector);
        vm.prank(ALICE);
        falseMarket.repay(10e24);

        (,, uint256 borrowAssets, uint256 borrowShares,) = falseMarket.market();
        (, uint256 userBorrowShares,) = falseMarket.position(ALICE);
        assertEq(borrowAssets, 10e18);
        assertEq(borrowShares, 10e24);
        assertEq(userBorrowShares, 10e24);
        assertEq(falseToken.balanceOf(ALICE), 10e18);
    }

    function test_ReentrantRepaymentIsBlocked() public {
        ReentrantRepayERC20 reentrantToken = new ReentrantRepayERC20();
        SmolphoRepayHarness reentrantMarket = _deploy(reentrantToken, 0);
        _seedDebtWithToken(reentrantMarket, reentrantToken);
        reentrantToken.setTarget(reentrantMarket);
        (, uint256 shares,) = reentrantMarket.position(ALICE);

        vm.prank(ALICE);
        reentrantMarket.repay(shares);

        assertTrue(reentrantToken.attemptedReentry());
        assertFalse(reentrantToken.reentrySucceeded());
        assertEq(reentrantToken.reentryError(), Smolpho.Reentrancy.selector);
        assertEq(reentrantMarket.borrowAssets(ALICE), 0);
    }

    function _deploy(IERC20 loan, uint256 rate) internal returns (SmolphoRepayHarness) {
        return new SmolphoRepayHarness(loan, collateralToken, oracle, rate);
    }

    function _seedDebt(SmolphoRepayHarness market_) internal {
        _seedDebtWithToken(market_, loanToken);
    }

    function _seedDebtWithToken(Smolpho market_, MockERC20 loan) internal {
        loan.mint(BOB, SUPPLY_ASSETS);
        vm.prank(BOB);
        loan.approve(address(market_), SUPPLY_ASSETS);
        vm.prank(BOB);
        market_.supply(SUPPLY_ASSETS);

        collateralToken.mint(ALICE, COLLATERAL_ASSETS);
        vm.prank(ALICE);
        collateralToken.approve(address(market_), COLLATERAL_ASSETS);
        vm.prank(ALICE);
        market_.supplyCollateral(COLLATERAL_ASSETS);

        vm.prank(ALICE);
        market_.borrow(BORROW_ASSETS);
        vm.prank(ALICE);
        loan.approve(address(market_), type(uint256).max);
    }
}
