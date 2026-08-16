// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {IERC20} from "./interfaces/IERC20.sol";
import {IPriceOracle} from "./interfaces/IPriceOracle.sol";
import {SharesMath} from "./libraries/SharesMath.sol";

contract Smolpho {
    uint256 public constant WAD = 1e18;
    uint256 public constant VIRTUAL_ASSETS = 1;
    uint256 public constant VIRTUAL_SHARES = 1e6;

    IERC20 public immutable loanToken;
    IERC20 public immutable collateralToken;
    IPriceOracle public immutable oracle;
    uint256 public immutable lltv;
    uint256 public immutable ratePerSecond;
    uint256 public immutable liquidationIncentive;

    struct Market {
        uint256 totalSupplyAssets;
        uint256 totalSupplyShares;
        uint256 totalBorrowAssets;
        uint256 totalBorrowShares;
        uint256 lastUpdate;
    }

    struct Position {
        uint256 supplyShares;
        uint256 borrowShares;
        uint256 collateral;
    }

    Market public market;
    mapping(address => Position) public position;

    event InterestAccrued(uint256 elapsed, uint256 interest);

    error ZeroAddress();
    error SameToken();
    error InvalidLltv();
    error InvalidLiquidationIncentive();

    constructor(
        IERC20 loanToken_,
        IERC20 collateralToken_,
        IPriceOracle oracle_,
        uint256 lltv_,
        uint256 ratePerSecond_,
        uint256 liquidationIncentive_
    ) {
        if (
            address(loanToken_) == address(0) || address(collateralToken_) == address(0)
                || address(oracle_) == address(0)
        ) {
            revert ZeroAddress();
        }
        if (address(loanToken_) == address(collateralToken_)) revert SameToken();
        if (lltv_ == 0 || lltv_ >= WAD) revert InvalidLltv();
        if (liquidationIncentive_ < WAD || lltv_ > (WAD * WAD - 1) / liquidationIncentive_) {
            revert InvalidLiquidationIncentive();
        }

        loanToken = loanToken_;
        collateralToken = collateralToken_;
        oracle = oracle_;
        lltv = lltv_;
        ratePerSecond = ratePerSecond_;
        liquidationIncentive = liquidationIncentive_;
        market.lastUpdate = block.timestamp;
    }

    function accrueInterest() public returns (uint256 interest) {
        uint256 elapsed = block.timestamp - market.lastUpdate;
        interest = market.totalBorrowAssets * ratePerSecond * elapsed / WAD;

        market.totalBorrowAssets += interest;
        market.totalSupplyAssets += interest;
        market.lastUpdate = block.timestamp;

        emit InterestAccrued(elapsed, interest);
    }

    function previewSupply(uint256 assets) external view returns (uint256) {
        return SharesMath.toSharesDown(assets, market.totalSupplyAssets, market.totalSupplyShares);
    }

    function previewWithdraw(uint256 shares) external view returns (uint256) {
        return SharesMath.toAssetsDown(shares, market.totalSupplyAssets, market.totalSupplyShares);
    }

    function previewBorrow(uint256 assets) external view returns (uint256) {
        return SharesMath.toSharesUp(assets, market.totalBorrowAssets, market.totalBorrowShares);
    }

    function previewRepay(uint256 shares) external view returns (uint256) {
        return SharesMath.toAssetsUp(shares, market.totalBorrowAssets, market.totalBorrowShares);
    }
}
