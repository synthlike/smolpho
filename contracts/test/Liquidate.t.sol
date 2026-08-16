// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {Smolpho} from "../src/Smolpho.sol";
import {IERC20} from "../src/interfaces/IERC20.sol";
import {SafeTransferLib} from "../src/libraries/SafeTransferLib.sol";
import {MockERC20} from "./mocks/MockERC20.sol";
import {MockOracle} from "./mocks/MockOracle.sol";

contract SmolphoLiquidationHarness is Smolpho {
    constructor(IERC20 loanToken_, IERC20 collateralToken_, MockOracle oracle_, uint256 ratePerSecond_)
        Smolpho(loanToken_, collateralToken_, oracle_, 0.8e18, ratePerSecond_, 1.05e18)
    {}
}

contract ToggleFalseTransferERC20 is MockERC20 {
    bool public transferShouldFail;

    constructor() MockERC20("Toggle Transfer", "TOGGLE") {}

    function setTransferShouldFail(bool shouldFail) external {
        transferShouldFail = shouldFail;
    }

    function transfer(address to, uint256 value) external override returns (bool) {
        if (transferShouldFail) return false;
        _transfer(msg.sender, to, value);
        return true;
    }
}

contract ReentrantLiquidationERC20 is MockERC20 {
    Smolpho public target;
    address public borrower;
    bool public attemptedReentry;
    bool public reentrySucceeded;
    bytes4 public reentryError;

    constructor() MockERC20("Reentrant Liquidation", "REENTER") {}

    function setReentry(Smolpho target_, address borrower_) external {
        target = target_;
        borrower = borrower_;
    }

    function transferFrom(address from, address to, uint256 value) public override returns (bool) {
        if (address(target) != address(0) && !attemptedReentry) {
            attemptedReentry = true;
            try target.liquidate(borrower, 1) {
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

contract LiquidateTest is Test {
    uint256 internal constant START_TIME = 1_000_000;
    uint256 internal constant INITIAL_PRICE = 2_000e18;
    uint256 internal constant SUPPLY_ASSETS = 5_000e18;
    uint256 internal constant BORROW_ASSETS = 3_000e18;
    uint256 internal constant COLLATERAL_ASSETS = 2e18;

    address internal constant ALICE = address(0xA11CE);
    address internal constant BOB = address(0xB0B);
    address internal constant CAROL = address(0xCA401);

    MockERC20 internal loanToken;
    MockERC20 internal collateralToken;
    MockOracle internal oracle;
    SmolphoLiquidationHarness internal smolpho;

    function setUp() public {
        vm.warp(START_TIME);
        loanToken = new MockERC20("Loan Token", "LOAN");
        collateralToken = new MockERC20("Collateral Token", "COLLATERAL");
        oracle = new MockOracle(INITIAL_PRICE);
        smolpho = _deploy(loanToken, collateralToken, 0);
        _seedPosition(smolpho, loanToken, collateralToken, BORROW_ASSETS);
    }

    function test_PartiallyLiquidatesUnhealthyPosition() public {
        oracle.setPrice(1_800e18);
        uint256 repaidShares = 100e24;
        uint256 expectedSeized = uint256(100e18) * 1.05e18 / 1_800e18;

        vm.expectEmit(address(smolpho));
        emit Smolpho.InterestAccrued(0, 0);
        vm.expectEmit(address(smolpho));
        emit Smolpho.Liquidated(CAROL, ALICE, 100e18, repaidShares, expectedSeized);

        vm.prank(CAROL);
        (uint256 repaidAssets, uint256 seizedCollateral, uint256 badDebtAssets) = smolpho.liquidate(ALICE, repaidShares);

        assertEq(repaidAssets, 100e18);
        assertEq(seizedCollateral, expectedSeized);
        assertEq(badDebtAssets, 0);
        assertEq(loanToken.balanceOf(CAROL), SUPPLY_ASSETS - 100e18);
        assertEq(collateralToken.balanceOf(CAROL), expectedSeized);
        assertEq(smolpho.borrowAssets(ALICE), 2_900e18);

        (,, uint256 borrowAssets, uint256 borrowShares,) = smolpho.market();
        (, uint256 userBorrowShares, uint256 collateral) = smolpho.position(ALICE);
        assertEq(borrowAssets, 2_900e18);
        assertEq(borrowShares, 2_900e24);
        assertEq(userBorrowShares, 2_900e24);
        assertEq(collateral, COLLATERAL_ASSETS - expectedSeized);
    }

    function test_RevertsWhenPositionIsHealthy() public {
        vm.expectRevert(Smolpho.HealthyPosition.selector);
        vm.prank(CAROL);
        smolpho.liquidate(ALICE, 100e24);
    }

    function test_RevertsForZeroRepaymentShares() public {
        oracle.setPrice(1_800e18);

        vm.expectRevert(Smolpho.ZeroShares.selector);
        vm.prank(CAROL);
        smolpho.liquidate(ALICE, 0);
    }

    function test_RevertsForExcessRepaymentShares() public {
        oracle.setPrice(1_800e18);

        vm.expectRevert(Smolpho.InsufficientBorrowShares.selector);
        vm.prank(CAROL);
        smolpho.liquidate(ALICE, 3_000e24 + 1);
    }

    function test_AccruedInterestCanMakePositionLiquidatable() public {
        SmolphoLiquidationHarness interestMarket = _deploy(loanToken, collateralToken, 1e15);
        _seedPosition(interestMarket, loanToken, collateralToken, 3_190e18);
        vm.warp(START_TIME + 10);

        vm.prank(CAROL);
        (uint256 repaidAssets,,) = interestMarket.liquidate(ALICE, 100e24);

        assertGt(repaidAssets, 100e18);
        (,,,, uint256 lastUpdate) = interestMarket.market();
        assertEq(lastUpdate, START_TIME + 10);
    }

    function test_ExhaustedCollateralRealizesRemainingBadDebt() public {
        oracle.setPrice(100e18);
        uint256 repaidShares = 200e24;
        uint256 badDebtShares = 2_800e24;

        vm.expectEmit(address(smolpho));
        emit Smolpho.Liquidated(CAROL, ALICE, 200e18, repaidShares, COLLATERAL_ASSETS);
        vm.expectEmit(address(smolpho));
        emit Smolpho.BadDebtRealized(ALICE, 2_800e18, badDebtShares);

        vm.prank(CAROL);
        (uint256 repaidAssets, uint256 seizedCollateral, uint256 badDebtAssets) = smolpho.liquidate(ALICE, repaidShares);

        assertEq(repaidAssets, 200e18);
        assertEq(seizedCollateral, COLLATERAL_ASSETS);
        assertEq(badDebtAssets, 2_800e18);
        (uint256 supplyAssets,, uint256 borrowAssets, uint256 borrowShares,) = smolpho.market();
        (, uint256 userBorrowShares, uint256 collateral) = smolpho.position(ALICE);
        assertEq(supplyAssets, 2_200e18);
        assertEq(smolpho.supplyAssets(BOB), 2_200e18);
        assertEq(borrowAssets, 0);
        assertEq(borrowShares, 0);
        assertEq(userBorrowShares, 0);
        assertEq(collateral, 0);
    }

    function test_ZeroPriceSeizesAllCollateralAndRealizesBadDebt() public {
        oracle.setPrice(0);

        vm.prank(CAROL);
        (, uint256 seizedCollateral, uint256 badDebtAssets) = smolpho.liquidate(ALICE, 1e24);

        assertEq(seizedCollateral, COLLATERAL_ASSETS);
        assertEq(badDebtAssets, 2_999e18);
        assertEq(smolpho.borrowAssets(ALICE), 0);
    }

    function test_CollateralTransferFailureRollsBackEverything() public {
        ToggleFalseTransferERC20 falseCollateral = new ToggleFalseTransferERC20();
        SmolphoLiquidationHarness falseMarket = _deploy(loanToken, falseCollateral, 0);
        _seedPosition(falseMarket, loanToken, falseCollateral, BORROW_ASSETS);
        oracle.setPrice(1_800e18);
        falseCollateral.setTransferShouldFail(true);
        uint256 liquidatorBalanceBefore = loanToken.balanceOf(CAROL);

        vm.expectRevert(SafeTransferLib.TransferFailed.selector);
        vm.prank(CAROL);
        falseMarket.liquidate(ALICE, 100e24);

        assertEq(loanToken.balanceOf(CAROL), liquidatorBalanceBefore);
        assertEq(falseMarket.borrowAssets(ALICE), BORROW_ASSETS);
        (,, uint256 collateral) = falseMarket.position(ALICE);
        assertEq(collateral, COLLATERAL_ASSETS);
    }

    function test_ReentrantLiquidationIsBlocked() public {
        ReentrantLiquidationERC20 reentrantLoan = new ReentrantLiquidationERC20();
        SmolphoLiquidationHarness reentrantMarket = _deploy(reentrantLoan, collateralToken, 0);
        _seedPosition(reentrantMarket, reentrantLoan, collateralToken, BORROW_ASSETS);
        oracle.setPrice(1_800e18);
        reentrantLoan.setReentry(reentrantMarket, ALICE);

        vm.prank(CAROL);
        reentrantMarket.liquidate(ALICE, 100e24);

        assertTrue(reentrantLoan.attemptedReentry());
        assertFalse(reentrantLoan.reentrySucceeded());
        assertEq(reentrantLoan.reentryError(), Smolpho.Reentrancy.selector);
    }

    function _deploy(IERC20 loan, IERC20 collateral, uint256 rate) internal returns (SmolphoLiquidationHarness) {
        return new SmolphoLiquidationHarness(loan, collateral, oracle, rate);
    }

    function _seedPosition(Smolpho market_, MockERC20 loan, MockERC20 collateral, uint256 borrowedAssets) internal {
        loan.mint(BOB, SUPPLY_ASSETS);
        vm.prank(BOB);
        loan.approve(address(market_), SUPPLY_ASSETS);
        vm.prank(BOB);
        market_.supply(SUPPLY_ASSETS);

        collateral.mint(ALICE, COLLATERAL_ASSETS);
        vm.prank(ALICE);
        collateral.approve(address(market_), COLLATERAL_ASSETS);
        vm.prank(ALICE);
        market_.supplyCollateral(COLLATERAL_ASSETS);

        vm.prank(ALICE);
        market_.borrow(borrowedAssets);

        loan.mint(CAROL, SUPPLY_ASSETS);
        vm.prank(CAROL);
        loan.approve(address(market_), type(uint256).max);
    }
}
