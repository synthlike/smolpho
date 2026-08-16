// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

library SharesMath {
    uint256 internal constant VIRTUAL_ASSETS = 1;
    uint256 internal constant VIRTUAL_SHARES = 1e6;

    function toSharesDown(uint256 assets, uint256 totalAssets, uint256 totalShares) internal pure returns (uint256) {
        return assets * (totalShares + VIRTUAL_SHARES) / (totalAssets + VIRTUAL_ASSETS);
    }

    function toSharesUp(uint256 assets, uint256 totalAssets, uint256 totalShares) internal pure returns (uint256) {
        return mulDivUp(assets, totalShares + VIRTUAL_SHARES, totalAssets + VIRTUAL_ASSETS);
    }

    function toAssetsDown(uint256 shares, uint256 totalAssets, uint256 totalShares) internal pure returns (uint256) {
        return shares * (totalAssets + VIRTUAL_ASSETS) / (totalShares + VIRTUAL_SHARES);
    }

    function toAssetsUp(uint256 shares, uint256 totalAssets, uint256 totalShares) internal pure returns (uint256) {
        return mulDivUp(shares, totalAssets + VIRTUAL_ASSETS, totalShares + VIRTUAL_SHARES);
    }

    function mulDivUp(uint256 x, uint256 y, uint256 denominator) private pure returns (uint256 result) {
        uint256 product = x * y;
        result = product / denominator;
        if (product % denominator != 0) result += 1;
    }
}
