// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {Smolpho} from "../src/Smolpho.sol";
import {IERC20} from "../src/interfaces/IERC20.sol";
import {SafeTransferLib} from "../src/libraries/SafeTransferLib.sol";
import {SharesMath} from "../src/libraries/SharesMath.sol";
import {FalseTransferERC20, MockERC20} from "./mocks/MockERC20.sol";
import {MockOracle} from "./mocks/MockOracle.sol";

contract SmolphoBorrowHarness is Smolpho {
    constructor(IERC20 loanToken_, IERC20 collateralToken_, MockOracle oracle_, uint256 ratePerSecond_)
        Smolpho(loanToken_, collateralToken_, oracle_, 0.8e18, ratePerSecond_, 1.05e18)
    {}

    function setBorrowTotals(uint256 totalBorrowAssets, uint256 totalBorrowShares) external {
        market.totalBorrowAssets = uint128(totalBorrowAssets);
        market.totalBorrowShares = uint128(totalBorrowShares);
    }
}

contract ReentrantBorrowERC20 is MockERC20 {
    Smolpho public target;
    bool public attemptedReentry;
    bool public reentrySucceeded;
    bytes4 public reentryError;

    constructor() MockERC20("Reentrant Loan", "REENTER") {}

    function setTarget(Smolpho target_) external {
        target = target_;
    }

    function transfer(address to, uint256 value) external override returns (bool) {
        if (!attemptedReentry) {
            attemptedReentry = true;
            try target.borrow(1) {
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

contract BorrowTest is Test {
    uint256 internal constant START_TIME = 1_000_000;
    uint256 internal constant PRICE = 2_000e18;
    uint256 internal constant SUPPLY_ASSETS = 5_000e18;
    uint256 internal constant COLLATERAL_ASSETS = 2e18;

    address internal constant ALICE = address(0xA11CE);
    address internal constant BOB = address(0xB0B);
    address internal constant CAROL = address(0xCA401);

    MockERC20 internal loanToken;
    MockERC20 internal collateralToken;
    MockOracle internal oracle;
    SmolphoBorrowHarness internal smolpho;

    function setUp() public {
        vm.warp(START_TIME);
        loanToken = new MockERC20("Loan Token", "LOAN");
        collateralToken = new MockERC20("Collateral Token", "COLLATERAL");
        oracle = new MockOracle(PRICE);
        smolpho = _deploy(loanToken, collateralToken, 0);

        loanToken.mint(BOB, 10_000e18);
        collateralToken.mint(ALICE, 100e18);
        _seedMarket(smolpho, loanToken, collateralToken, SUPPLY_ASSETS, COLLATERAL_ASSETS);
    }

    function test_BorrowsAtExactHealthBoundary() public {
        uint256 assets = 3_200e18;
        uint256 expectedShares = assets * 1e6;

        vm.expectEmit(address(smolpho));
        emit Smolpho.InterestAccrued(0, 0);
        vm.expectEmit(address(smolpho));
        emit Smolpho.Borrowed(ALICE, assets, expectedShares);

        vm.prank(ALICE);
        uint256 shares = smolpho.borrow(assets);

        assertEq(shares, expectedShares);
        assertEq(loanToken.balanceOf(ALICE), assets);
        assertEq(loanToken.balanceOf(address(smolpho)), SUPPLY_ASSETS - assets);
        assertEq(smolpho.borrowAssets(ALICE), assets);
        assertEq(smolpho.availableLiquidity(), SUPPLY_ASSETS - assets);
        assertTrue(smolpho.isHealthy(ALICE));

        (uint256 supplyAssets,, uint256 borrowAssets, uint256 borrowShares,) = smolpho.market();
        (, uint256 userBorrowShares, uint256 collateral) = smolpho.position(ALICE);
        assertEq(supplyAssets, SUPPLY_ASSETS);
        assertEq(borrowAssets, assets);
        assertEq(borrowShares, expectedShares);
        assertEq(userBorrowShares, expectedShares);
        assertEq(collateral, COLLATERAL_ASSETS);
    }

    function test_RevertsWithoutCollateral() public {
        vm.expectRevert(Smolpho.UnhealthyPosition.selector);
        vm.prank(CAROL);
        smolpho.borrow(1e18);

        assertEq(smolpho.borrowAssets(CAROL), 0);
        assertEq(loanToken.balanceOf(CAROL), 0);
    }

    function test_RevertsWhenBorrowExceedsHealthLimitByOneUnit() public {
        vm.expectRevert(Smolpho.UnhealthyPosition.selector);
        vm.prank(ALICE);
        smolpho.borrow(3_200e18 + 1);

        assertEq(smolpho.borrowAssets(ALICE), 0);
    }

    function test_PriceDeclineCanPreventBorrowing() public {
        oracle.setPrice(100e18);

        vm.expectRevert(Smolpho.UnhealthyPosition.selector);
        vm.prank(ALICE);
        smolpho.borrow(200e18);
    }

    function test_RevertsWhenMarketHasInsufficientLiquidity() public {
        SmolphoBorrowHarness smallMarket = _deploy(loanToken, collateralToken, 0);
        loanToken.mint(BOB, 100e18);
        collateralToken.mint(ALICE, COLLATERAL_ASSETS);
        _seedMarket(smallMarket, loanToken, collateralToken, 100e18, COLLATERAL_ASSETS);

        vm.expectRevert(Smolpho.InsufficientLiquidity.selector);
        vm.prank(ALICE);
        smallMarket.borrow(101e18);

        assertEq(smallMarket.borrowAssets(ALICE), 0);
        assertEq(loanToken.balanceOf(address(smallMarket)), 100e18);
    }

    function test_DebtSharesRoundUp() public {
        smolpho.setBorrowTotals(2, 0);

        vm.prank(ALICE);
        uint256 shares = smolpho.borrow(1);

        assertEq(shares, 333_334);
    }

    function test_AccruesInterestBeforeCalculatingDebtShares() public {
        SmolphoBorrowHarness interestMarket = _deploy(loanToken, collateralToken, 1e15);
        loanToken.mint(BOB, SUPPLY_ASSETS);
        collateralToken.mint(ALICE, COLLATERAL_ASSETS);
        _seedMarket(interestMarket, loanToken, collateralToken, SUPPLY_ASSETS, COLLATERAL_ASSETS);
        vm.prank(ALICE);
        interestMarket.borrow(100e18);
        vm.warp(START_TIME + 10);

        uint256 expectedShares = SharesMath.toSharesUp(1e18, 101e18, 100e24);

        vm.prank(ALICE);
        uint256 shares = interestMarket.borrow(1e18);

        assertEq(shares, expectedShares);
        (,, uint256 borrowAssets,, uint256 lastUpdate) = interestMarket.market();
        assertEq(borrowAssets, 102e18);
        assertEq(lastUpdate, START_TIME + 10);
    }

    function test_RevertsForZeroAssets() public {
        vm.expectRevert(Smolpho.ZeroAssets.selector);
        vm.prank(ALICE);
        smolpho.borrow(0);
    }

    function test_RevertsWhenAssetsExceedStorageBound() public {
        vm.expectRevert(Smolpho.AmountTooLarge.selector);
        vm.prank(ALICE);
        smolpho.borrow(type(uint128).max);
    }

    function test_TransferFailureRollsBackDebtAccounting() public {
        FalseTransferERC20 falseToken = new FalseTransferERC20();
        SmolphoBorrowHarness falseMarket = _deploy(falseToken, collateralToken, 0);
        falseToken.mint(BOB, 100e18);
        collateralToken.mint(ALICE, COLLATERAL_ASSETS);
        _seedMarket(falseMarket, falseToken, collateralToken, 100e18, COLLATERAL_ASSETS);

        vm.expectRevert(SafeTransferLib.TransferFailed.selector);
        vm.prank(ALICE);
        falseMarket.borrow(10e18);

        (,, uint256 borrowAssets, uint256 borrowShares,) = falseMarket.market();
        (, uint256 userBorrowShares,) = falseMarket.position(ALICE);
        assertEq(borrowAssets, 0);
        assertEq(borrowShares, 0);
        assertEq(userBorrowShares, 0);
        assertEq(falseToken.balanceOf(address(falseMarket)), 100e18);
    }

    function test_ReentrantBorrowIsBlocked() public {
        ReentrantBorrowERC20 reentrantToken = new ReentrantBorrowERC20();
        SmolphoBorrowHarness reentrantMarket = _deploy(reentrantToken, collateralToken, 0);
        reentrantToken.mint(BOB, 100e18);
        collateralToken.mint(ALICE, COLLATERAL_ASSETS);
        _seedMarket(reentrantMarket, reentrantToken, collateralToken, 100e18, COLLATERAL_ASSETS);
        reentrantToken.setTarget(reentrantMarket);

        vm.prank(ALICE);
        reentrantMarket.borrow(10e18);

        assertTrue(reentrantToken.attemptedReentry());
        assertFalse(reentrantToken.reentrySucceeded());
        assertEq(reentrantToken.reentryError(), Smolpho.Reentrancy.selector);
        assertEq(reentrantToken.balanceOf(ALICE), 10e18);
    }

    function _deploy(IERC20 loan, IERC20 collateral, uint256 rate) internal returns (SmolphoBorrowHarness) {
        return new SmolphoBorrowHarness(loan, collateral, oracle, rate);
    }

    function _seedMarket(
        Smolpho market_,
        MockERC20 loan,
        MockERC20 collateral,
        uint256 suppliedAssets,
        uint256 collateralAssets
    ) internal {
        vm.prank(BOB);
        loan.approve(address(market_), suppliedAssets);
        vm.prank(BOB);
        market_.supply(suppliedAssets);

        vm.prank(ALICE);
        collateral.approve(address(market_), collateralAssets);
        vm.prank(ALICE);
        market_.supplyCollateral(collateralAssets);
    }
}
