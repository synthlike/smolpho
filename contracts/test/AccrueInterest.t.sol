// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {Smolpho} from "../src/Smolpho.sol";
import {IERC20} from "../src/interfaces/IERC20.sol";
import {IPriceOracle} from "../src/interfaces/IPriceOracle.sol";

contract SmolphoInterestHarness is Smolpho {
    constructor(uint256 ratePerSecond_)
        Smolpho(
            IERC20(address(0x1000)),
            IERC20(address(0x2000)),
            IPriceOracle(address(0x3000)),
            0.8e18,
            ratePerSecond_,
            1.05e18
        )
    {}

    function setMarket(
        uint256 totalSupplyAssets,
        uint256 totalSupplyShares,
        uint256 totalBorrowAssets,
        uint256 totalBorrowShares,
        uint256 lastUpdate
    ) external {
        market = Market({
            totalSupplyAssets: uint128(totalSupplyAssets),
            totalSupplyShares: uint128(totalSupplyShares),
            totalBorrowAssets: uint128(totalBorrowAssets),
            totalBorrowShares: uint128(totalBorrowShares),
            lastUpdate: uint64(lastUpdate)
        });
    }
}

contract AccrueInterestTest is Test {
    uint256 internal constant START_TIME = 1_000_000;
    uint256 internal constant RATE_PER_SECOND = 1e15;

    SmolphoInterestHarness internal smolpho;

    function setUp() public {
        vm.warp(START_TIME);
        smolpho = new SmolphoInterestHarness(RATE_PER_SECOND);
    }

    function test_NoElapsedTimeProducesZeroInterest() public {
        vm.expectEmit();
        emit Smolpho.InterestAccrued(0, 0);

        uint256 interest = smolpho.accrueInterest();

        assertEq(interest, 0);
        (,,,, uint256 lastUpdate) = smolpho.market();
        assertEq(lastUpdate, START_TIME);
    }

    function test_NoDebtProducesZeroInterestAndAdvancesTime() public {
        smolpho.setMarket(50e18, 12, 0, 0, START_TIME);
        vm.warp(START_TIME + 10);

        vm.expectEmit();
        emit Smolpho.InterestAccrued(10, 0);

        uint256 interest = smolpho.accrueInterest();

        (uint256 supplyAssets, uint256 supplyShares, uint256 borrowAssets, uint256 borrowShares, uint256 lastUpdate) =
            smolpho.market();
        assertEq(interest, 0);
        assertEq(supplyAssets, 50e18);
        assertEq(supplyShares, 12);
        assertEq(borrowAssets, 0);
        assertEq(borrowShares, 0);
        assertEq(lastUpdate, START_TIME + 10);
    }

    function test_AddsEqualInterestWithoutChangingShares() public {
        smolpho.setMarket(150e18, 77, 100e18, 33, START_TIME);
        vm.warp(START_TIME + 10);

        vm.expectEmit();
        emit Smolpho.InterestAccrued(10, 1e18);

        uint256 interest = smolpho.accrueInterest();

        (uint256 supplyAssets, uint256 supplyShares, uint256 borrowAssets, uint256 borrowShares, uint256 lastUpdate) =
            smolpho.market();
        assertEq(interest, 1e18);
        assertEq(supplyAssets, 151e18);
        assertEq(borrowAssets, 101e18);
        assertEq(supplyShares, 77);
        assertEq(borrowShares, 33);
        assertEq(lastUpdate, START_TIME + 10);
    }

    function test_SecondAccrualAtSameTimestampAddsNothing() public {
        smolpho.setMarket(150e18, 77, 100e18, 33, START_TIME);
        vm.warp(START_TIME + 10);
        smolpho.accrueInterest();

        vm.expectEmit();
        emit Smolpho.InterestAccrued(0, 0);
        uint256 secondInterest = smolpho.accrueInterest();

        (uint256 supplyAssets,, uint256 borrowAssets,, uint256 lastUpdate) = smolpho.market();
        assertEq(secondInterest, 0);
        assertEq(supplyAssets, 151e18);
        assertEq(borrowAssets, 101e18);
        assertEq(lastUpdate, START_TIME + 10);
    }

    function testFuzz_InterestNeverDecreasesTotals(uint256 borrowAssets, uint256 rate, uint256 elapsed) public {
        borrowAssets = bound(borrowAssets, 0, 1e30);
        rate = bound(rate, 0, type(uint64).max);
        elapsed = bound(elapsed, 0, 1e6);

        SmolphoInterestHarness fuzzMarket = new SmolphoInterestHarness(rate);
        fuzzMarket.setMarket(borrowAssets, 1, borrowAssets, 1, START_TIME);
        vm.warp(START_TIME + elapsed);

        uint256 expectedInterest = borrowAssets * rate * elapsed / 1e18;
        uint256 interest = fuzzMarket.accrueInterest();

        (uint256 supplyAssets,, uint256 accruedBorrowAssets,,) = fuzzMarket.market();
        assertEq(interest, expectedInterest);
        assertEq(supplyAssets, borrowAssets + expectedInterest);
        assertEq(accruedBorrowAssets, borrowAssets + expectedInterest);
        assertGe(supplyAssets, borrowAssets);
        assertGe(accruedBorrowAssets, borrowAssets);
    }
}
