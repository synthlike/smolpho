// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {IERC20} from "../interfaces/IERC20.sol";

library SafeTransferLib {
    error TransferFailed();

    function safeTransfer(IERC20 token, address to, uint256 value) internal {
        (bool success, bytes memory returnData) = address(token).call(abi.encodeCall(token.transfer, (to, value)));

        if (!success || (returnData.length != 0 && (returnData.length < 32 || !abi.decode(returnData, (bool))))) {
            revert TransferFailed();
        }
    }

    function safeTransferFrom(IERC20 token, address from, address to, uint256 value) internal {
        (bool success, bytes memory returnData) =
            address(token).call(abi.encodeCall(token.transferFrom, (from, to, value)));

        if (!success || (returnData.length != 0 && (returnData.length < 32 || !abi.decode(returnData, (bool))))) {
            revert TransferFailed();
        }
    }
}
