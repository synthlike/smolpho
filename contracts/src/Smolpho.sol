// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {IERC20} from "./interfaces/IERC20.sol";
import {IPriceOracle} from "./interfaces/IPriceOracle.sol";
import {SafeTransferLib} from "./libraries/SafeTransferLib.sol";
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

    uint256 private locked = 1;

    event InterestAccrued(uint256 elapsed, uint256 interest);
    event Supplied(address indexed user, uint256 assets, uint256 shares);
    event Withdrawn(address indexed user, uint256 assets, uint256 shares);
    event CollateralSupplied(address indexed user, uint256 assets);

    error ZeroAddress();
    error ZeroAssets();
    error ZeroShares();
    error InsufficientSupplyShares();
    error InsufficientLiquidity();
    error Reentrancy();
    error SameToken();
    error InvalidLltv();
    error InvalidLiquidationIncentive();

    modifier nonReentrant() {
        if (locked != 1) revert Reentrancy();
        locked = 2;
        _;
        locked = 1;
    }

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

    function supply(uint256 assets) external nonReentrant returns (uint256 shares) {
        if (assets == 0) revert ZeroAssets();

        accrueInterest();
        shares = SharesMath.toSharesDown(assets, market.totalSupplyAssets, market.totalSupplyShares);

        position[msg.sender].supplyShares += shares;
        market.totalSupplyShares += shares;
        market.totalSupplyAssets += assets;

        SafeTransferLib.safeTransferFrom(loanToken, msg.sender, address(this), assets);

        emit Supplied(msg.sender, assets, shares);
    }

    function supplyCollateral(uint256 assets) external nonReentrant {
        if (assets == 0) revert ZeroAssets();

        position[msg.sender].collateral += assets;

        SafeTransferLib.safeTransferFrom(collateralToken, msg.sender, address(this), assets);

        emit CollateralSupplied(msg.sender, assets);
    }

    function withdraw(uint256 shares) external nonReentrant returns (uint256 assets) {
        if (shares == 0) revert ZeroShares();
        if (shares > position[msg.sender].supplyShares) revert InsufficientSupplyShares();

        accrueInterest();
        assets = SharesMath.toAssetsDown(shares, market.totalSupplyAssets, market.totalSupplyShares);

        if (assets > market.totalSupplyAssets - market.totalBorrowAssets) revert InsufficientLiquidity();

        position[msg.sender].supplyShares -= shares;
        market.totalSupplyShares -= shares;
        market.totalSupplyAssets -= assets;

        SafeTransferLib.safeTransfer(loanToken, msg.sender, assets);

        emit Withdrawn(msg.sender, assets, shares);
    }

    function supplyAssets(address user) external view returns (uint256) {
        return SharesMath.toAssetsDown(position[user].supplyShares, market.totalSupplyAssets, market.totalSupplyShares);
    }

    function availableLiquidity() external view returns (uint256) {
        return market.totalSupplyAssets - market.totalBorrowAssets;
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
