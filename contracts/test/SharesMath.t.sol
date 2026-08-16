// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {Smolpho} from "../src/Smolpho.sol";
import {IERC20} from "../src/interfaces/IERC20.sol";
import {IPriceOracle} from "../src/interfaces/IPriceOracle.sol";
import {SharesMath} from "../src/libraries/SharesMath.sol";

contract SharesMathHarness {
    function toSharesDown(uint256 assets, uint256 totalAssets, uint256 totalShares) external pure returns (uint256) {
        return SharesMath.toSharesDown(assets, totalAssets, totalShares);
    }

    function toSharesUp(uint256 assets, uint256 totalAssets, uint256 totalShares) external pure returns (uint256) {
        return SharesMath.toSharesUp(assets, totalAssets, totalShares);
    }

    function toAssetsDown(uint256 shares, uint256 totalAssets, uint256 totalShares) external pure returns (uint256) {
        return SharesMath.toAssetsDown(shares, totalAssets, totalShares);
    }

    function toAssetsUp(uint256 shares, uint256 totalAssets, uint256 totalShares) external pure returns (uint256) {
        return SharesMath.toAssetsUp(shares, totalAssets, totalShares);
    }
}

contract SmolphoPreviewHarness is Smolpho {
    constructor()
        Smolpho(IERC20(address(0x1000)), IERC20(address(0x2000)), IPriceOracle(address(0x3000)), 0.8e18, 0, 1.05e18)
    {}

    function setTotals(
        uint256 totalSupplyAssets,
        uint256 totalSupplyShares,
        uint256 totalBorrowAssets,
        uint256 totalBorrowShares
    ) external {
        market.totalSupplyAssets = uint128(totalSupplyAssets);
        market.totalSupplyShares = uint128(totalSupplyShares);
        market.totalBorrowAssets = uint128(totalBorrowAssets);
        market.totalBorrowShares = uint128(totalBorrowShares);
    }
}

contract SharesMathTest is Test {
    uint256 internal constant MAX_FUZZ_VALUE = 1e24;

    SharesMathHarness internal math;

    function setUp() public {
        math = new SharesMathHarness();
    }

    function test_InitialRateIsOneMillionSharesPerAsset() public view {
        assertEq(math.toSharesDown(1, 0, 0), 1e6);
        assertEq(math.toSharesUp(1, 0, 0), 1e6);
        assertEq(math.toAssetsDown(1e6, 0, 0), 1);
        assertEq(math.toAssetsUp(1e6, 0, 0), 1);
    }

    function test_ConversionsRoundInTheRequiredDirection() public view {
        assertEq(math.toSharesDown(1, 2, 0), 333_333);
        assertEq(math.toSharesUp(1, 2, 0), 333_334);
        assertEq(math.toAssetsDown(1, 1, 0), 0);
        assertEq(math.toAssetsUp(1, 1, 0), 1);
    }

    function test_ZeroInputsReturnZero() public view {
        assertEq(math.toSharesDown(0, 0, 0), 0);
        assertEq(math.toSharesUp(0, 10, 20), 0);
        assertEq(math.toAssetsDown(0, 0, 0), 0);
        assertEq(math.toAssetsUp(0, 10, 20), 0);
    }

    function test_ChangedRateMintsFewerSharesPerAsset() public view {
        uint256 initialShares = math.toSharesDown(1, 0, 0);
        uint256 sharesAfterInterest = math.toSharesDown(1, 2, 1e6);

        assertEq(initialShares, 1e6);
        assertEq(sharesAfterInterest, 666_666);
        assertLt(sharesAfterInterest, initialShares);
    }

    function testFuzz_ConversionsAreMonotonic(uint256 smaller, uint256 larger, uint256 totalAssets, uint256 totalShares)
        public
        view
    {
        smaller = bound(smaller, 0, MAX_FUZZ_VALUE);
        larger = bound(larger, smaller, MAX_FUZZ_VALUE);
        totalAssets = bound(totalAssets, 0, MAX_FUZZ_VALUE);
        totalShares = bound(totalShares, 0, MAX_FUZZ_VALUE);

        assertLe(
            math.toSharesDown(smaller, totalAssets, totalShares), math.toSharesDown(larger, totalAssets, totalShares)
        );
        assertLe(math.toSharesUp(smaller, totalAssets, totalShares), math.toSharesUp(larger, totalAssets, totalShares));
        assertLe(
            math.toAssetsDown(smaller, totalAssets, totalShares), math.toAssetsDown(larger, totalAssets, totalShares)
        );
        assertLe(math.toAssetsUp(smaller, totalAssets, totalShares), math.toAssetsUp(larger, totalAssets, totalShares));
    }

    function testFuzz_DownRoundTripCannotProfit(uint256 assets, uint256 totalAssets, uint256 totalShares) public view {
        assets = bound(assets, 0, MAX_FUZZ_VALUE);
        totalAssets = bound(totalAssets, 0, MAX_FUZZ_VALUE);
        totalShares = bound(totalShares, 0, MAX_FUZZ_VALUE);

        uint256 shares = math.toSharesDown(assets, totalAssets, totalShares);
        uint256 assetsBack = math.toAssetsDown(shares, totalAssets, totalShares);

        assertLe(assetsBack, assets);
    }

    function testFuzz_UpRoundTripCannotUndercharge(uint256 assets, uint256 totalAssets, uint256 totalShares)
        public
        view
    {
        assets = bound(assets, 0, MAX_FUZZ_VALUE);
        totalAssets = bound(totalAssets, 0, MAX_FUZZ_VALUE);
        totalShares = bound(totalShares, 0, MAX_FUZZ_VALUE);

        uint256 shares = math.toSharesUp(assets, totalAssets, totalShares);
        uint256 assetsBack = math.toAssetsUp(shares, totalAssets, totalShares);

        assertGe(assetsBack, assets);
    }
}

contract SmolphoPreviewTest is Test {
    SmolphoPreviewHarness internal smolpho;

    function setUp() public {
        smolpho = new SmolphoPreviewHarness();
    }

    function test_PreviewsUseExpectedRoundingAndSeparateTotals() public {
        smolpho.setTotals({totalSupplyAssets: 2, totalSupplyShares: 0, totalBorrowAssets: 1, totalBorrowShares: 0});

        assertEq(smolpho.previewSupply(1), 333_333);
        assertEq(smolpho.previewWithdraw(1), 0);
        assertEq(smolpho.previewBorrow(1), 500_000);
        assertEq(smolpho.previewRepay(1), 1);
    }
}
