// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {Smolpho} from "../src/Smolpho.sol";
import {IPriceOracle} from "../src/interfaces/IPriceOracle.sol";
import {SafeTransferLib} from "../src/libraries/SafeTransferLib.sol";
import {FalseReturnERC20, MockERC20} from "./mocks/MockERC20.sol";

contract ReentrantCollateralERC20 is MockERC20 {
    Smolpho public target;
    bool public attemptedReentry;
    bool public reentrySucceeded;
    bytes4 public reentryError;

    constructor() MockERC20("Reentrant Collateral", "REENTER") {}

    function setTarget(Smolpho target_) external {
        target = target_;
    }

    function transferFrom(address from, address to, uint256 value) public override returns (bool) {
        if (!attemptedReentry) {
            attemptedReentry = true;
            try target.supplyCollateral(1) {
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

contract SupplyCollateralTest is Test {
    uint256 internal constant START_TIME = 1_000_000;
    uint256 internal constant ASSETS = 10e18;

    address internal constant ALICE = address(0xA11CE);

    MockERC20 internal loanToken;
    MockERC20 internal collateralToken;
    Smolpho internal smolpho;

    function setUp() public {
        vm.warp(START_TIME);
        loanToken = new MockERC20("Loan Token", "LOAN");
        collateralToken = new MockERC20("Collateral Token", "COLLATERAL");
        smolpho = _deploy(collateralToken);

        collateralToken.mint(ALICE, 100e18);
        vm.prank(ALICE);
        collateralToken.approve(address(smolpho), type(uint256).max);
    }

    function test_SuppliesCollateralWithoutChangingMarketAccounting() public {
        vm.expectEmit(address(smolpho));
        emit Smolpho.CollateralSupplied(ALICE, ASSETS);

        vm.prank(ALICE);
        smolpho.supplyCollateral(ASSETS);

        assertEq(collateralToken.balanceOf(ALICE), 90e18);
        assertEq(collateralToken.balanceOf(address(smolpho)), ASSETS);
        assertEq(loanToken.balanceOf(address(smolpho)), 0);

        (uint256 supplyShares, uint256 borrowShares, uint256 collateral) = smolpho.position(ALICE);
        assertEq(supplyShares, 0);
        assertEq(borrowShares, 0);
        assertEq(collateral, ASSETS);

        (uint256 supplyAssets, uint256 totalSupplyShares, uint256 borrowAssets, uint256 totalBorrowShares,) =
            smolpho.market();
        assertEq(supplyAssets, 0);
        assertEq(totalSupplyShares, 0);
        assertEq(borrowAssets, 0);
        assertEq(totalBorrowShares, 0);
    }

    function test_MultipleCollateralSuppliesAccumulate() public {
        vm.startPrank(ALICE);
        smolpho.supplyCollateral(4e18);
        smolpho.supplyCollateral(6e18);
        vm.stopPrank();

        (,, uint256 collateral) = smolpho.position(ALICE);
        assertEq(collateral, ASSETS);
        assertEq(collateralToken.balanceOf(address(smolpho)), ASSETS);
    }

    function test_DoesNotAccrueInterest() public {
        vm.warp(START_TIME + 30 days);

        vm.prank(ALICE);
        smolpho.supplyCollateral(ASSETS);

        (,,,, uint256 lastUpdate) = smolpho.market();
        assertEq(lastUpdate, START_TIME);
    }

    function test_RevertsForZeroAssets() public {
        vm.expectRevert(Smolpho.ZeroAssets.selector);
        vm.prank(ALICE);
        smolpho.supplyCollateral(0);
    }

    function test_TransferFailureRollsBackCollateralAccounting() public {
        vm.prank(ALICE);
        collateralToken.approve(address(smolpho), 0);

        vm.expectRevert(SafeTransferLib.TransferFailed.selector);
        vm.prank(ALICE);
        smolpho.supplyCollateral(ASSETS);

        (,, uint256 collateral) = smolpho.position(ALICE);
        assertEq(collateral, 0);
        assertEq(collateralToken.balanceOf(address(smolpho)), 0);
    }

    function test_FalseReturnTokenRevertsAndRollsBackAccounting() public {
        FalseReturnERC20 falseToken = new FalseReturnERC20();
        Smolpho falseMarket = _deploy(falseToken);
        falseToken.mint(ALICE, ASSETS);
        vm.prank(ALICE);
        falseToken.approve(address(falseMarket), ASSETS);

        vm.expectRevert(SafeTransferLib.TransferFailed.selector);
        vm.prank(ALICE);
        falseMarket.supplyCollateral(ASSETS);

        (,, uint256 collateral) = falseMarket.position(ALICE);
        assertEq(collateral, 0);
        assertEq(falseToken.balanceOf(address(falseMarket)), 0);
    }

    function test_ReentrantCollateralSupplyIsBlocked() public {
        ReentrantCollateralERC20 reentrantToken = new ReentrantCollateralERC20();
        Smolpho reentrantMarket = _deploy(reentrantToken);
        reentrantToken.setTarget(reentrantMarket);
        reentrantToken.mint(ALICE, ASSETS);
        vm.prank(ALICE);
        reentrantToken.approve(address(reentrantMarket), ASSETS);

        vm.prank(ALICE);
        reentrantMarket.supplyCollateral(ASSETS);

        assertTrue(reentrantToken.attemptedReentry());
        assertFalse(reentrantToken.reentrySucceeded());
        assertEq(reentrantToken.reentryError(), Smolpho.Reentrancy.selector);
        assertEq(reentrantToken.balanceOf(address(reentrantMarket)), ASSETS);
        (,, uint256 collateral) = reentrantMarket.position(ALICE);
        assertEq(collateral, ASSETS);
    }

    function _deploy(MockERC20 collateral) internal returns (Smolpho) {
        return new Smolpho(loanToken, collateral, IPriceOracle(address(0x3000)), 0.8e18, 1e15, 1.05e18);
    }
}
