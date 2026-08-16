// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {Smolpho} from "../src/Smolpho.sol";
import {IERC20} from "../src/interfaces/IERC20.sol";
import {IPriceOracle} from "../src/interfaces/IPriceOracle.sol";

contract SmolphoMarketCreationTest is Test {
    IERC20 internal constant LOAN_TOKEN = IERC20(address(0x1000));
    IERC20 internal constant COLLATERAL_TOKEN = IERC20(address(0x2000));
    IPriceOracle internal constant ORACLE = IPriceOracle(address(0x3000));

    uint256 internal constant LLTV = 0.8e18;
    uint256 internal constant RATE_PER_SECOND = 1e9;
    uint256 internal constant LIQUIDATION_INCENTIVE = 1.05e18;

    function test_CreatesMarketWithImmutableConfiguration() public {
        vm.warp(1_234_567);

        Smolpho smolpho = _deploy(LLTV, RATE_PER_SECOND, LIQUIDATION_INCENTIVE);

        assertEq(address(smolpho.loanToken()), address(LOAN_TOKEN));
        assertEq(address(smolpho.collateralToken()), address(COLLATERAL_TOKEN));
        assertEq(address(smolpho.oracle()), address(ORACLE));
        assertEq(smolpho.lltv(), LLTV);
        assertEq(smolpho.ratePerSecond(), RATE_PER_SECOND);
        assertEq(smolpho.liquidationIncentive(), LIQUIDATION_INCENTIVE);
        assertEq(smolpho.WAD(), 1e18);

        (
            uint256 totalSupplyAssets,
            uint256 totalSupplyShares,
            uint256 totalBorrowAssets,
            uint256 totalBorrowShares,
            uint256 lastUpdate
        ) = smolpho.market();
        assertEq(totalSupplyAssets, 0);
        assertEq(totalSupplyShares, 0);
        assertEq(totalBorrowAssets, 0);
        assertEq(totalBorrowShares, 0);
        assertEq(lastUpdate, block.timestamp);

        (uint256 supplyShares, uint256 borrowShares, uint256 collateral) = smolpho.position(address(this));
        assertEq(supplyShares, 0);
        assertEq(borrowShares, 0);
        assertEq(collateral, 0);
    }

    function test_AllowsZeroInterestRate() public {
        Smolpho smolpho = _deploy(LLTV, 0, LIQUIDATION_INCENTIVE);
        assertEq(smolpho.ratePerSecond(), 0);
    }

    function test_RevertsForZeroLoanToken() public {
        vm.expectRevert(Smolpho.ZeroAddress.selector);
        new Smolpho(IERC20(address(0)), COLLATERAL_TOKEN, ORACLE, LLTV, RATE_PER_SECOND, LIQUIDATION_INCENTIVE);
    }

    function test_RevertsForZeroCollateralToken() public {
        vm.expectRevert(Smolpho.ZeroAddress.selector);
        new Smolpho(LOAN_TOKEN, IERC20(address(0)), ORACLE, LLTV, RATE_PER_SECOND, LIQUIDATION_INCENTIVE);
    }

    function test_RevertsForZeroOracle() public {
        vm.expectRevert(Smolpho.ZeroAddress.selector);
        new Smolpho(
            LOAN_TOKEN, COLLATERAL_TOKEN, IPriceOracle(address(0)), LLTV, RATE_PER_SECOND, LIQUIDATION_INCENTIVE
        );
    }

    function test_RevertsForIdenticalTokens() public {
        vm.expectRevert(Smolpho.SameToken.selector);
        new Smolpho(LOAN_TOKEN, LOAN_TOKEN, ORACLE, LLTV, RATE_PER_SECOND, LIQUIDATION_INCENTIVE);
    }

    function test_RevertsForZeroLltv() public {
        vm.expectRevert(Smolpho.InvalidLltv.selector);
        _deploy(0, RATE_PER_SECOND, LIQUIDATION_INCENTIVE);
    }

    function test_RevertsForLltvEqualToWad() public {
        vm.expectRevert(Smolpho.InvalidLltv.selector);
        _deploy(1e18, RATE_PER_SECOND, LIQUIDATION_INCENTIVE);
    }

    function test_RevertsForLltvAboveWad() public {
        vm.expectRevert(Smolpho.InvalidLltv.selector);
        _deploy(1e18 + 1, RATE_PER_SECOND, LIQUIDATION_INCENTIVE);
    }

    function test_RevertsForLiquidationIncentiveBelowWad() public {
        vm.expectRevert(Smolpho.InvalidLiquidationIncentive.selector);
        _deploy(LLTV, RATE_PER_SECOND, 1e18 - 1);
    }

    function test_RevertsWhenLltvTimesIncentiveEqualsWadSquared() public {
        vm.expectRevert(Smolpho.InvalidLiquidationIncentive.selector);
        _deploy(0.8e18, RATE_PER_SECOND, 1.25e18);
    }

    function test_RevertsWhenLltvTimesIncentiveExceedsWadSquared() public {
        vm.expectRevert(Smolpho.InvalidLiquidationIncentive.selector);
        _deploy(0.8e18, RATE_PER_SECOND, 1.25e18 + 1);
    }

    function test_RevertsForExtremeIncentiveWithoutOverflowing() public {
        vm.expectRevert(Smolpho.InvalidLiquidationIncentive.selector);
        _deploy(LLTV, RATE_PER_SECOND, type(uint256).max);
    }

    function _deploy(uint256 lltv, uint256 ratePerSecond, uint256 liquidationIncentive) internal returns (Smolpho) {
        return new Smolpho(LOAN_TOKEN, COLLATERAL_TOKEN, ORACLE, lltv, ratePerSecond, liquidationIncentive);
    }
}
