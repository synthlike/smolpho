// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {Smolpho} from "../src/Smolpho.sol";
import {IERC20} from "../src/interfaces/IERC20.sol";
import {IPriceOracle} from "../src/interfaces/IPriceOracle.sol";
import {SafeTransferLib} from "../src/libraries/SafeTransferLib.sol";
import {SharesMath} from "../src/libraries/SharesMath.sol";
import {FalseTransferERC20, MockERC20} from "./mocks/MockERC20.sol";

contract SmolphoWithdrawHarness is Smolpho {
    constructor(IERC20 loanToken_, uint256 ratePerSecond_)
        Smolpho(loanToken_, IERC20(address(0x2000)), IPriceOracle(address(0x3000)), 0.8e18, ratePerSecond_, 1.05e18)
    {}

    function setBorrowTotals(uint256 totalBorrowAssets, uint256 totalBorrowShares) external {
        market.totalBorrowAssets = totalBorrowAssets;
        market.totalBorrowShares = totalBorrowShares;
    }

    function setSupplyPosition(address user, uint256 totalSupplyAssets, uint256 totalSupplyShares, uint256 userShares)
        external
    {
        market.totalSupplyAssets = totalSupplyAssets;
        market.totalSupplyShares = totalSupplyShares;
        position[user].supplyShares = userShares;
    }
}

contract ReentrantTransferERC20 is MockERC20 {
    Smolpho public target;
    bool public attemptedReentry;
    bool public reentrySucceeded;
    bytes4 public reentryError;

    constructor() MockERC20("Reentrant Transfer", "REENTER") {}

    function setTarget(Smolpho target_) external {
        target = target_;
    }

    function transfer(address to, uint256 value) external override returns (bool) {
        if (!attemptedReentry) {
            attemptedReentry = true;
            try target.withdraw(1) {
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

contract WithdrawTest is Test {
    uint256 internal constant START_TIME = 1_000_000;
    uint256 internal constant ASSETS = 10e18;
    uint256 internal constant SHARES = ASSETS * 1e6;

    address internal constant ALICE = address(0xA11CE);

    MockERC20 internal loanToken;
    SmolphoWithdrawHarness internal smolpho;

    function setUp() public {
        vm.warp(START_TIME);
        loanToken = new MockERC20("Loan Token", "LOAN");
        smolpho = new SmolphoWithdrawHarness(loanToken, 0);

        loanToken.mint(ALICE, 100e18);
        vm.prank(ALICE);
        loanToken.approve(address(smolpho), type(uint256).max);
        vm.prank(ALICE);
        smolpho.supply(ASSETS);
    }

    function test_FullWithdrawalBurnsSharesAndReturnsAssets() public {
        vm.expectEmit(address(smolpho));
        emit Smolpho.InterestAccrued(0, 0);
        vm.expectEmit(address(smolpho));
        emit Smolpho.Withdrawn(ALICE, ASSETS, SHARES);

        vm.prank(ALICE);
        uint256 assets = smolpho.withdraw(SHARES);

        assertEq(assets, ASSETS);
        assertEq(loanToken.balanceOf(ALICE), 100e18);
        assertEq(loanToken.balanceOf(address(smolpho)), 0);

        (uint256 supplyAssets, uint256 supplyShares, uint256 borrowAssets, uint256 borrowShares,) = smolpho.market();
        (uint256 userSupplyShares,,) = smolpho.position(ALICE);
        assertEq(supplyAssets, 0);
        assertEq(supplyShares, 0);
        assertEq(borrowAssets, 0);
        assertEq(borrowShares, 0);
        assertEq(userSupplyShares, 0);
        assertEq(smolpho.supplyAssets(ALICE), 0);
        assertEq(smolpho.availableLiquidity(), 0);
    }

    function test_WithdrawalIncludesAccruedInterestAndRoundsDown() public {
        SmolphoWithdrawHarness interestMarket = new SmolphoWithdrawHarness(loanToken, 1e15);
        loanToken.mint(ALICE, 100e18);
        vm.prank(ALICE);
        loanToken.approve(address(interestMarket), type(uint256).max);
        vm.prank(ALICE);
        uint256 suppliedShares = interestMarket.supply(100e18);
        interestMarket.setBorrowTotals(50e18, 50e24);
        vm.warp(START_TIME + 10);

        uint256 sharesToWithdraw = suppliedShares / 10;
        uint256 expectedInterest = 0.5e18;
        uint256 expectedAssets = SharesMath.toAssetsDown(sharesToWithdraw, 100e18 + expectedInterest, suppliedShares);

        vm.prank(ALICE);
        uint256 assets = interestMarket.withdraw(sharesToWithdraw);

        assertEq(assets, expectedAssets);
        (uint256 supplyAssets, uint256 supplyShares, uint256 borrowAssets,,) = interestMarket.market();
        assertEq(supplyAssets, 100e18 + expectedInterest - expectedAssets);
        assertEq(supplyShares, suppliedShares - sharesToWithdraw);
        assertEq(borrowAssets, 50e18 + expectedInterest);
    }

    function test_RevertsForZeroShares() public {
        vm.expectRevert(Smolpho.ZeroShares.selector);
        vm.prank(ALICE);
        smolpho.withdraw(0);
    }

    function test_RevertsForInsufficientUserShares() public {
        vm.expectRevert(Smolpho.InsufficientSupplyShares.selector);
        vm.prank(ALICE);
        smolpho.withdraw(SHARES + 1);
    }

    function test_RevertsWhenWithdrawalExceedsAvailableLiquidity() public {
        smolpho.setBorrowTotals(9e18, 9e24);
        assertEq(smolpho.availableLiquidity(), 1e18);

        vm.expectRevert(Smolpho.InsufficientLiquidity.selector);
        vm.prank(ALICE);
        smolpho.withdraw(SHARES);

        (uint256 supplyAssets, uint256 supplyShares,,,) = smolpho.market();
        (uint256 userSupplyShares,,) = smolpho.position(ALICE);
        assertEq(supplyAssets, ASSETS);
        assertEq(supplyShares, SHARES);
        assertEq(userSupplyShares, SHARES);
    }

    function test_TransferFailureRollsBackAccounting() public {
        FalseTransferERC20 falseToken = new FalseTransferERC20();
        SmolphoWithdrawHarness falseMarket = new SmolphoWithdrawHarness(falseToken, 0);
        falseToken.mint(ALICE, ASSETS);
        vm.prank(ALICE);
        falseToken.approve(address(falseMarket), ASSETS);
        vm.prank(ALICE);
        uint256 shares = falseMarket.supply(ASSETS);

        vm.expectRevert(SafeTransferLib.TransferFailed.selector);
        vm.prank(ALICE);
        falseMarket.withdraw(shares);

        (uint256 supplyAssets, uint256 supplyShares,,,) = falseMarket.market();
        (uint256 userSupplyShares,,) = falseMarket.position(ALICE);
        assertEq(supplyAssets, ASSETS);
        assertEq(supplyShares, shares);
        assertEq(userSupplyShares, shares);
        assertEq(falseToken.balanceOf(address(falseMarket)), ASSETS);
    }

    function test_ReentrantWithdrawalIsBlocked() public {
        ReentrantTransferERC20 reentrantToken = new ReentrantTransferERC20();
        SmolphoWithdrawHarness reentrantMarket = new SmolphoWithdrawHarness(reentrantToken, 0);
        reentrantToken.mint(ALICE, ASSETS);
        vm.prank(ALICE);
        reentrantToken.approve(address(reentrantMarket), ASSETS);
        vm.prank(ALICE);
        uint256 shares = reentrantMarket.supply(ASSETS);
        reentrantToken.setTarget(reentrantMarket);

        vm.prank(ALICE);
        reentrantMarket.withdraw(shares);

        assertTrue(reentrantToken.attemptedReentry());
        assertFalse(reentrantToken.reentrySucceeded());
        assertEq(reentrantToken.reentryError(), Smolpho.Reentrancy.selector);
        assertEq(reentrantToken.balanceOf(address(reentrantMarket)), 0);
        assertEq(reentrantToken.balanceOf(ALICE), ASSETS);
    }

    function test_OneShareCanRoundDownToZeroAssets() public {
        SmolphoWithdrawHarness roundingMarket = new SmolphoWithdrawHarness(loanToken, 0);
        roundingMarket.setSupplyPosition(ALICE, 1, 1, 1);

        vm.prank(ALICE);
        uint256 assets = roundingMarket.withdraw(1);

        assertEq(assets, 0);
        (uint256 supplyAssets, uint256 supplyShares,,,) = roundingMarket.market();
        assertEq(supplyAssets, 1);
        assertEq(supplyShares, 0);
    }
}
