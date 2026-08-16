// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {Smolpho} from "../src/Smolpho.sol";
import {IERC20} from "../src/interfaces/IERC20.sol";
import {SafeTransferLib} from "../src/libraries/SafeTransferLib.sol";
import {FalseTransferERC20, MockERC20} from "./mocks/MockERC20.sol";
import {MockOracle} from "./mocks/MockOracle.sol";

contract SmolphoWithdrawCollateralHarness is Smolpho {
    constructor(IERC20 loanToken_, IERC20 collateralToken_, MockOracle oracle_, uint256 ratePerSecond_)
        Smolpho(loanToken_, collateralToken_, oracle_, 0.8e18, ratePerSecond_, 1.05e18)
    {}

    function setDebt(address user, uint256 debtAssets) external {
        uint256 debtShares = debtAssets * 1e6;
        market.totalSupplyAssets = uint128(debtAssets);
        market.totalBorrowAssets = uint128(debtAssets);
        market.totalBorrowShares = uint128(debtShares);
        position[user].borrowShares = uint128(debtShares);
    }
}

contract ReentrantCollateralTransferERC20 is MockERC20 {
    Smolpho public target;
    bool public attemptedReentry;
    bool public reentrySucceeded;
    bytes4 public reentryError;

    constructor() MockERC20("Reentrant Collateral", "REENTER") {}

    function setTarget(Smolpho target_) external {
        target = target_;
    }

    function transfer(address to, uint256 value) external override returns (bool) {
        if (!attemptedReentry) {
            attemptedReentry = true;
            try target.withdrawCollateral(1) {
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

        _transfer(msg.sender, to, value);
        return true;
    }
}

contract WithdrawCollateralTest is Test {
    uint256 internal constant START_TIME = 1_000_000;
    uint256 internal constant PRICE = 2_000e18;

    address internal constant ALICE = address(0xA11CE);

    MockERC20 internal loanToken;
    MockERC20 internal collateralToken;
    MockOracle internal oracle;
    SmolphoWithdrawCollateralHarness internal smolpho;

    function setUp() public {
        vm.warp(START_TIME);
        loanToken = new MockERC20("Loan Token", "LOAN");
        collateralToken = new MockERC20("Collateral Token", "COLLATERAL");
        oracle = new MockOracle(PRICE);
        smolpho = _deploy(collateralToken, 0);

        collateralToken.mint(ALICE, 100e18);
        vm.prank(ALICE);
        collateralToken.approve(address(smolpho), type(uint256).max);
    }

    function test_DebtFreeUserCanWithdrawAllWithoutCallingOracle() public {
        _supplyCollateral(smolpho, collateralToken, 10e18);
        oracle.setShouldRevert(true);

        vm.prank(ALICE);
        smolpho.withdrawCollateral(10e18);

        (,, uint256 collateral) = smolpho.position(ALICE);
        assertEq(collateral, 0);
        assertEq(collateralToken.balanceOf(ALICE), 100e18);
        assertEq(collateralToken.balanceOf(address(smolpho)), 0);
    }

    function test_PartialWithdrawalTransfersCollateralAndEmitsEvent() public {
        _supplyCollateral(smolpho, collateralToken, 10e18);

        vm.expectEmit(address(smolpho));
        emit Smolpho.InterestAccrued(0, 0);
        vm.expectEmit(address(smolpho));
        emit Smolpho.CollateralWithdrawn(ALICE, 4e18);

        vm.prank(ALICE);
        smolpho.withdrawCollateral(4e18);

        (,, uint256 collateral) = smolpho.position(ALICE);
        assertEq(collateral, 6e18);
        assertEq(collateralToken.balanceOf(ALICE), 94e18);
        assertEq(collateralToken.balanceOf(address(smolpho)), 6e18);
    }

    function test_WithdrawalCanLeavePositionAtExactHealthBoundary() public {
        _supplyCollateral(smolpho, collateralToken, 3e18);
        smolpho.setDebt(ALICE, 3_200e18);

        vm.prank(ALICE);
        smolpho.withdrawCollateral(1e18);

        (,, uint256 collateral) = smolpho.position(ALICE);
        assertEq(collateral, 2e18);
        assertTrue(smolpho.isHealthy(ALICE));
    }

    function test_RevertsWhenWithdrawalMakesPositionUnhealthy() public {
        _supplyCollateral(smolpho, collateralToken, 3e18);
        smolpho.setDebt(ALICE, 3_200e18);

        vm.expectRevert(Smolpho.UnhealthyPosition.selector);
        vm.prank(ALICE);
        smolpho.withdrawCollateral(1e18 + 1);

        (,, uint256 collateral) = smolpho.position(ALICE);
        assertEq(collateral, 3e18);
        assertEq(collateralToken.balanceOf(address(smolpho)), 3e18);
    }

    function test_AccruedInterestCanMakeWithdrawalUnsafe() public {
        SmolphoWithdrawCollateralHarness interestMarket = _deploy(collateralToken, 1e15);
        vm.prank(ALICE);
        collateralToken.approve(address(interestMarket), type(uint256).max);
        _supplyCollateral(interestMarket, collateralToken, 2e18);
        interestMarket.setDebt(ALICE, 3_100e18);
        vm.warp(START_TIME + 10);

        vm.expectRevert(Smolpho.UnhealthyPosition.selector);
        vm.prank(ALICE);
        interestMarket.withdrawCollateral(0.05e18);

        (,, uint256 collateral) = interestMarket.position(ALICE);
        (,, uint256 borrowAssets,, uint256 lastUpdate) = interestMarket.market();
        assertEq(collateral, 2e18);
        assertEq(borrowAssets, 3_100e18);
        assertEq(lastUpdate, START_TIME);
    }

    function test_RevertsForZeroAssets() public {
        vm.expectRevert(Smolpho.ZeroAssets.selector);
        vm.prank(ALICE);
        smolpho.withdrawCollateral(0);
    }

    function test_RevertsForInsufficientCollateral() public {
        _supplyCollateral(smolpho, collateralToken, 1e18);

        vm.expectRevert(Smolpho.InsufficientCollateral.selector);
        vm.prank(ALICE);
        smolpho.withdrawCollateral(1e18 + 1);
    }

    function test_TransferFailureRollsBackAccounting() public {
        FalseTransferERC20 falseToken = new FalseTransferERC20();
        SmolphoWithdrawCollateralHarness falseMarket = _deploy(falseToken, 0);
        falseToken.mint(ALICE, 10e18);
        vm.prank(ALICE);
        falseToken.approve(address(falseMarket), type(uint256).max);
        _supplyCollateral(falseMarket, falseToken, 10e18);

        vm.expectRevert(SafeTransferLib.TransferFailed.selector);
        vm.prank(ALICE);
        falseMarket.withdrawCollateral(4e18);

        (,, uint256 collateral) = falseMarket.position(ALICE);
        assertEq(collateral, 10e18);
        assertEq(falseToken.balanceOf(address(falseMarket)), 10e18);
    }

    function test_ReentrantWithdrawalIsBlocked() public {
        ReentrantCollateralTransferERC20 reentrantToken = new ReentrantCollateralTransferERC20();
        SmolphoWithdrawCollateralHarness reentrantMarket = _deploy(reentrantToken, 0);
        reentrantToken.mint(ALICE, 10e18);
        vm.prank(ALICE);
        reentrantToken.approve(address(reentrantMarket), type(uint256).max);
        _supplyCollateral(reentrantMarket, reentrantToken, 10e18);
        reentrantToken.setTarget(reentrantMarket);

        vm.prank(ALICE);
        reentrantMarket.withdrawCollateral(10e18);

        assertTrue(reentrantToken.attemptedReentry());
        assertFalse(reentrantToken.reentrySucceeded());
        assertEq(reentrantToken.reentryError(), Smolpho.Reentrancy.selector);
        assertEq(reentrantToken.balanceOf(ALICE), 10e18);
    }

    function _deploy(IERC20 collateral, uint256 rate) internal returns (SmolphoWithdrawCollateralHarness) {
        return new SmolphoWithdrawCollateralHarness(loanToken, collateral, oracle, rate);
    }

    function _supplyCollateral(Smolpho market_, MockERC20 token, uint256 assets) internal {
        if (token.allowance(ALICE, address(market_)) < assets) {
            vm.prank(ALICE);
            token.approve(address(market_), type(uint256).max);
        }
        vm.prank(ALICE);
        market_.supplyCollateral(assets);
    }
}
