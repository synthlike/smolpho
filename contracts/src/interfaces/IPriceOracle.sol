// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

interface IPriceOracle {
    /// @return The loan-token value of one collateral token, scaled by 1e18.
    function price() external view returns (uint256);
}
