// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {IPriceOracle} from "../../src/interfaces/IPriceOracle.sol";

contract MockOracle is IPriceOracle {
    uint256 internal oraclePrice;
    bool public shouldRevert;

    error OracleReverted();

    constructor(uint256 price_) {
        oraclePrice = price_;
    }

    function setPrice(uint256 price_) external {
        oraclePrice = price_;
    }

    function setShouldRevert(bool shouldRevert_) external {
        shouldRevert = shouldRevert_;
    }

    function price() external view returns (uint256) {
        if (shouldRevert) revert OracleReverted();
        return oraclePrice;
    }
}
