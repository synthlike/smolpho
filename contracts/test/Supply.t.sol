// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {Smolpho} from "../src/Smolpho.sol";
import {IERC20} from "../src/interfaces/IERC20.sol";
import {IPriceOracle} from "../src/interfaces/IPriceOracle.sol";
import {SafeTransferLib} from "../src/libraries/SafeTransferLib.sol";
import {SharesMath} from "../src/libraries/SharesMath.sol";
import {FalseReturnERC20, MockERC20} from "./mocks/MockERC20.sol";

contract SmolphoSupplyHarness is Smolpho {
    constructor(IERC20 loanToken_, uint256 ratePerSecond_)
        Smolpho(loanToken_, IERC20(address(0x2000)), IPriceOracle(address(0x3000)), 0.8e18, ratePerSecond_, 1.05e18)
    {}

    function setMarket(
        uint256 totalSupplyAssets,
        uint256 totalSupplyShares,
        uint256 totalBorrowAssets,
        uint256 totalBorrowShares,
        uint256 lastUpdate
    ) external {
        market = Market({
            totalSupplyAssets: totalSupplyAssets,
            totalSupplyShares: totalSupplyShares,
            totalBorrowAssets: totalBorrowAssets,
            totalBorrowShares: totalBorrowShares,
            lastUpdate: lastUpdate
        });
    }
}

contract ReentrantERC20 is MockERC20 {
    Smolpho public target;
    bool public attemptedReentry;
    bool public reentrySucceeded;
    bytes4 public reentryError;

    constructor() MockERC20("Reentrant", "REENTER") {}

    function setTarget(Smolpho target_) external {
        target = target_;
    }

    function transferFrom(address from, address to, uint256 value) public override returns (bool) {
        if (!attemptedReentry) {
            attemptedReentry = true;
            try target.supply(1) {
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

contract SupplyTest is Test {
    uint256 internal constant START_TIME = 1_000_000;
    uint256 internal constant ASSETS = 10e18;

    address internal constant ALICE = address(0xA11CE);

    MockERC20 internal loanToken;
    SmolphoSupplyHarness internal smolpho;

    function setUp() public {
        vm.warp(START_TIME);
        loanToken = new MockERC20("Loan Token", "LOAN");
        smolpho = new SmolphoSupplyHarness(loanToken, 0);

        loanToken.mint(ALICE, 100e18);
        vm.prank(ALICE);
        loanToken.approve(address(smolpho), type(uint256).max);
    }

    function test_FirstSupplyMovesTokensAndMintsShares() public {
        uint256 expectedShares = ASSETS * 1e6;

        vm.expectEmit(address(smolpho));
        emit Smolpho.InterestAccrued(0, 0);
        vm.expectEmit(address(smolpho));
        emit Smolpho.Supplied(ALICE, ASSETS, expectedShares);

        vm.prank(ALICE);
        uint256 shares = smolpho.supply(ASSETS);

        assertEq(shares, expectedShares);
        assertEq(loanToken.balanceOf(ALICE), 90e18);
        assertEq(loanToken.balanceOf(address(smolpho)), ASSETS);

        (uint256 supplyAssets, uint256 supplyShares, uint256 borrowAssets, uint256 borrowShares,) = smolpho.market();
        (uint256 userSupplyShares, uint256 userBorrowShares, uint256 collateral) = smolpho.position(ALICE);
        assertEq(supplyAssets, ASSETS);
        assertEq(supplyShares, expectedShares);
        assertEq(borrowAssets, 0);
        assertEq(borrowShares, 0);
        assertEq(userSupplyShares, expectedShares);
        assertEq(userBorrowShares, 0);
        assertEq(collateral, 0);
        assertEq(smolpho.supplyAssets(ALICE), ASSETS);
    }

    function test_SupplyAccruesInterestBeforeCalculatingShares() public {
        SmolphoSupplyHarness interestMarket = new SmolphoSupplyHarness(loanToken, 1e15);
        interestMarket.setMarket(100e18, 100e24, 50e18, 50e24, START_TIME);
        vm.prank(ALICE);
        loanToken.approve(address(interestMarket), type(uint256).max);
        vm.warp(START_TIME + 10);

        uint256 expectedInterest = 0.5e18;
        uint256 expectedShares = SharesMath.toSharesDown(ASSETS, 100e18 + expectedInterest, 100e24);

        vm.prank(ALICE);
        uint256 shares = interestMarket.supply(ASSETS);

        assertEq(shares, expectedShares);
        (uint256 supplyAssets, uint256 supplyShares, uint256 borrowAssets,, uint256 lastUpdate) =
            interestMarket.market();
        assertEq(supplyAssets, 100e18 + expectedInterest + ASSETS);
        assertEq(supplyShares, 100e24 + expectedShares);
        assertEq(borrowAssets, 50e18 + expectedInterest);
        assertEq(lastUpdate, START_TIME + 10);
    }

    function test_RevertsForZeroAssets() public {
        vm.expectRevert(Smolpho.ZeroAssets.selector);
        vm.prank(ALICE);
        smolpho.supply(0);
    }

    function test_TransferFailureRollsBackAccounting() public {
        vm.prank(ALICE);
        loanToken.approve(address(smolpho), 0);

        vm.expectRevert(SafeTransferLib.TransferFailed.selector);
        vm.prank(ALICE);
        smolpho.supply(ASSETS);

        (uint256 supplyAssets, uint256 supplyShares,,,) = smolpho.market();
        (uint256 userSupplyShares,,) = smolpho.position(ALICE);
        assertEq(supplyAssets, 0);
        assertEq(supplyShares, 0);
        assertEq(userSupplyShares, 0);
        assertEq(loanToken.balanceOf(address(smolpho)), 0);
    }

    function test_FalseReturnTokenRevertsAndRollsBackAccounting() public {
        FalseReturnERC20 falseToken = new FalseReturnERC20();
        SmolphoSupplyHarness falseMarket = new SmolphoSupplyHarness(falseToken, 0);
        falseToken.mint(ALICE, ASSETS);
        vm.prank(ALICE);
        falseToken.approve(address(falseMarket), ASSETS);

        vm.expectRevert(SafeTransferLib.TransferFailed.selector);
        vm.prank(ALICE);
        falseMarket.supply(ASSETS);

        (uint256 supplyAssets, uint256 supplyShares,,,) = falseMarket.market();
        assertEq(supplyAssets, 0);
        assertEq(supplyShares, 0);
    }

    function test_ReentrantSupplyIsBlocked() public {
        ReentrantERC20 reentrantToken = new ReentrantERC20();
        SmolphoSupplyHarness reentrantMarket = new SmolphoSupplyHarness(reentrantToken, 0);
        reentrantToken.setTarget(reentrantMarket);
        reentrantToken.mint(ALICE, ASSETS);
        vm.prank(ALICE);
        reentrantToken.approve(address(reentrantMarket), ASSETS);

        vm.prank(ALICE);
        reentrantMarket.supply(ASSETS);

        assertTrue(reentrantToken.attemptedReentry());
        assertFalse(reentrantToken.reentrySucceeded());
        assertEq(reentrantToken.reentryError(), Smolpho.Reentrancy.selector);
        assertEq(reentrantToken.balanceOf(address(reentrantMarket)), ASSETS);
    }
}
