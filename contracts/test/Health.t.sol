// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {Smolpho} from "../src/Smolpho.sol";
import {IERC20} from "../src/interfaces/IERC20.sol";
import {MockOracle} from "./mocks/MockOracle.sol";

contract SmolphoHealthHarness is Smolpho {
    constructor(MockOracle oracle_, uint256 ratePerSecond_)
        Smolpho(IERC20(address(0x1000)), IERC20(address(0x2000)), oracle_, 0.8e18, ratePerSecond_, 1.05e18)
    {}

    function setDebt(address user, uint256 debtAssets, uint256 collateral) external {
        uint256 debtShares = debtAssets * 1e6;
        market.totalSupplyAssets = uint128(debtAssets);
        market.totalBorrowAssets = uint128(debtAssets);
        market.totalBorrowShares = uint128(debtShares);
        position[user].borrowShares = uint128(debtShares);
        position[user].collateral = uint128(collateral);
    }

    function setRawDebt(
        address user,
        uint256 totalBorrowAssets,
        uint256 totalBorrowShares,
        uint256 userBorrowShares,
        uint256 collateral
    ) external {
        market.totalSupplyAssets = uint128(totalBorrowAssets);
        market.totalBorrowAssets = uint128(totalBorrowAssets);
        market.totalBorrowShares = uint128(totalBorrowShares);
        position[user].borrowShares = uint128(userBorrowShares);
        position[user].collateral = uint128(collateral);
    }
}

contract HealthTest is Test {
    uint256 internal constant START_TIME = 1_000_000;
    uint256 internal constant INITIAL_PRICE = 2_000e18;

    address internal constant ALICE = address(0xA11CE);

    MockOracle internal oracle;
    SmolphoHealthHarness internal smolpho;

    function setUp() public {
        vm.warp(START_TIME);
        oracle = new MockOracle(INITIAL_PRICE);
        smolpho = new SmolphoHealthHarness(oracle, 0);
    }

    function test_DebtFreePositionIsHealthyWithoutCallingOracle() public {
        oracle.setShouldRevert(true);

        assertTrue(smolpho.isHealthy(ALICE));
        assertEq(smolpho.borrowAssets(ALICE), 0);
    }

    function test_PositionAtExactBorrowLimitIsHealthy() public {
        smolpho.setDebt(ALICE, 3_200e18, 2e18);

        assertEq(smolpho.borrowAssets(ALICE), 3_200e18);
        assertTrue(smolpho.isHealthy(ALICE));
    }

    function test_PositionOneUnitAboveBorrowLimitIsUnhealthy() public {
        smolpho.setDebt(ALICE, 3_200e18 + 1, 2e18);

        assertFalse(smolpho.isHealthy(ALICE));
    }

    function test_PositionWithDebtAndNoCollateralIsUnhealthy() public {
        smolpho.setDebt(ALICE, 1, 0);

        assertFalse(smolpho.isHealthy(ALICE));
    }

    function test_ZeroOraclePriceMakesDebtPositionUnhealthy() public {
        smolpho.setDebt(ALICE, 1, 1e18);
        oracle.setPrice(0);

        assertFalse(smolpho.isHealthy(ALICE));
    }

    function test_PriceDeclineMakesPositionUnhealthy() public {
        smolpho.setDebt(ALICE, 3_100e18, 2e18);
        assertTrue(smolpho.isHealthy(ALICE));

        oracle.setPrice(1_900e18);

        assertFalse(smolpho.isHealthy(ALICE));
    }

    function test_BorrowAssetsRoundsDebtUp() public {
        smolpho.setRawDebt(ALICE, 1, 0, 1, 1e18);

        assertEq(smolpho.borrowAssets(ALICE), 1);
    }

    function test_RevertsWhenOraclePriceExceedsUint128() public {
        smolpho.setDebt(ALICE, 1, 1e18);
        oracle.setPrice(uint256(type(uint128).max) + 1);

        vm.expectRevert(Smolpho.OraclePriceTooLarge.selector);
        smolpho.isHealthy(ALICE);
    }

    function test_HealthSupportsMaximumOraclePrice() public {
        smolpho.setDebt(ALICE, 1, 1e18);
        oracle.setPrice(type(uint128).max);

        assertTrue(smolpho.isHealthy(ALICE));
    }

    function test_ViewsUseStoredDebtUntilInterestIsAccrued() public {
        SmolphoHealthHarness interestMarket = new SmolphoHealthHarness(oracle, 1e15);
        interestMarket.setDebt(ALICE, 100e18, 1e18);
        vm.warp(START_TIME + 10);

        assertEq(interestMarket.borrowAssets(ALICE), 100e18);

        interestMarket.accrueInterest();

        assertGt(interestMarket.borrowAssets(ALICE), 100e18);
        assertEq(interestMarket.borrowAssets(ALICE), interestMarket.previewRepay(100e24));
    }
}
