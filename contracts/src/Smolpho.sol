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

    // Reserve virtual-value headroom so bounded operands can be safely multiplied as uint256.
    uint256 internal constant MAX_ASSETS = type(uint128).max - VIRTUAL_ASSETS;
    uint256 internal constant MAX_SHARES = type(uint128).max - VIRTUAL_SHARES;

    IERC20 public immutable loanToken;
    IERC20 public immutable collateralToken;
    IPriceOracle public immutable oracle;
    uint256 public immutable lltv;
    uint64 public immutable ratePerSecond;
    uint256 public immutable liquidationIncentive;

    struct Market {
        uint128 totalSupplyAssets;
        uint128 totalSupplyShares;
        uint128 totalBorrowAssets;
        uint128 totalBorrowShares;
        uint64 lastUpdate;
    }

    struct Position {
        uint128 supplyShares;
        uint128 borrowShares;
        uint128 collateral;
    }

    Market public market;
    mapping(address => Position) public position;

    uint256 private locked = 1;

    event InterestAccrued(uint256 elapsed, uint256 interest);
    event Supplied(address indexed user, uint256 assets, uint256 shares);
    event Withdrawn(address indexed user, uint256 assets, uint256 shares);
    event CollateralSupplied(address indexed user, uint256 assets);
    event CollateralWithdrawn(address indexed user, uint256 assets);
    event Borrowed(address indexed user, uint256 assets, uint256 shares);

    error ZeroAddress();
    error ZeroAssets();
    error ZeroShares();
    error InsufficientSupplyShares();
    error InsufficientCollateral();
    error InsufficientLiquidity();
    error UnhealthyPosition();
    error Reentrancy();
    error SameToken();
    error InvalidLltv();
    error InvalidLiquidationIncentive();
    error AmountTooLarge();
    error RateTooLarge();
    error TimestampTooLarge();
    error OraclePriceTooLarge();

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

        if (ratePerSecond_ > type(uint64).max) revert RateTooLarge();
        // This checks representability, not a security-sensitive time threshold.
        // forge-lint: disable-next-line(block-timestamp)
        if (block.timestamp > type(uint64).max) revert TimestampTooLarge();

        loanToken = loanToken_;
        collateralToken = collateralToken_;
        oracle = oracle_;
        lltv = lltv_;
        // Safe because ratePerSecond_ was bounded above.
        // forge-lint: disable-next-line(unsafe-typecast)
        ratePerSecond = uint64(ratePerSecond_);
        liquidationIncentive = liquidationIncentive_;
        market.lastUpdate = uint64(block.timestamp);
    }

    function accrueInterest() public returns (uint256 interest) {
        // This checks representability, not a security-sensitive time threshold.
        // forge-lint: disable-next-line(block-timestamp)
        if (block.timestamp > type(uint64).max) revert TimestampTooLarge();

        uint256 elapsed = block.timestamp - market.lastUpdate;
        interest = uint256(market.totalBorrowAssets) * uint256(ratePerSecond) * elapsed / WAD;

        market.totalBorrowAssets = _toAssets(uint256(market.totalBorrowAssets) + interest);
        market.totalSupplyAssets = _toAssets(uint256(market.totalSupplyAssets) + interest);
        market.lastUpdate = uint64(block.timestamp);

        emit InterestAccrued(elapsed, interest);
    }

    function supply(uint256 assets) external nonReentrant returns (uint256 shares) {
        if (assets == 0) revert ZeroAssets();

        if (assets > MAX_ASSETS) revert AmountTooLarge();

        accrueInterest();
        shares = SharesMath.toSharesDown(assets, market.totalSupplyAssets, market.totalSupplyShares);

        position[msg.sender].supplyShares = _toShares(uint256(position[msg.sender].supplyShares) + shares);
        market.totalSupplyShares = _toShares(uint256(market.totalSupplyShares) + shares);
        market.totalSupplyAssets = _toAssets(uint256(market.totalSupplyAssets) + assets);

        SafeTransferLib.safeTransferFrom(loanToken, msg.sender, address(this), assets);

        emit Supplied(msg.sender, assets, shares);
    }

    function supplyCollateral(uint256 assets) external nonReentrant {
        if (assets == 0) revert ZeroAssets();

        position[msg.sender].collateral = _toUint128(uint256(position[msg.sender].collateral) + assets);

        SafeTransferLib.safeTransferFrom(collateralToken, msg.sender, address(this), assets);

        emit CollateralSupplied(msg.sender, assets);
    }

    function withdrawCollateral(uint256 assets) external nonReentrant {
        if (assets == 0) revert ZeroAssets();
        if (assets > position[msg.sender].collateral) revert InsufficientCollateral();

        accrueInterest();
        position[msg.sender].collateral -= _toUint128(assets);

        if (!_isHealthy(msg.sender)) revert UnhealthyPosition();

        SafeTransferLib.safeTransfer(collateralToken, msg.sender, assets);

        emit CollateralWithdrawn(msg.sender, assets);
    }

    function borrow(uint256 assets) external nonReentrant returns (uint256 shares) {
        if (assets == 0) revert ZeroAssets();
        if (assets > MAX_ASSETS) revert AmountTooLarge();

        accrueInterest();
        shares = SharesMath.toSharesUp(assets, market.totalBorrowAssets, market.totalBorrowShares);

        position[msg.sender].borrowShares = _toShares(uint256(position[msg.sender].borrowShares) + shares);
        market.totalBorrowShares = _toShares(uint256(market.totalBorrowShares) + shares);
        market.totalBorrowAssets = _toAssets(uint256(market.totalBorrowAssets) + assets);

        if (!_isHealthy(msg.sender)) revert UnhealthyPosition();
        if (market.totalBorrowAssets > market.totalSupplyAssets) revert InsufficientLiquidity();

        SafeTransferLib.safeTransfer(loanToken, msg.sender, assets);

        emit Borrowed(msg.sender, assets, shares);
    }

    function withdraw(uint256 shares) external nonReentrant returns (uint256 assets) {
        if (shares == 0) revert ZeroShares();
        if (shares > position[msg.sender].supplyShares) revert InsufficientSupplyShares();

        accrueInterest();
        assets = SharesMath.toAssetsDown(shares, market.totalSupplyAssets, market.totalSupplyShares);

        if (assets > market.totalSupplyAssets - market.totalBorrowAssets) revert InsufficientLiquidity();

        uint128 shares128 = _toUint128(shares);
        uint128 assets128 = _toUint128(assets);
        position[msg.sender].supplyShares -= shares128;
        market.totalSupplyShares -= shares128;
        market.totalSupplyAssets -= assets128;

        SafeTransferLib.safeTransfer(loanToken, msg.sender, assets);

        emit Withdrawn(msg.sender, assets, shares);
    }

    function supplyAssets(address user) external view returns (uint256) {
        return SharesMath.toAssetsDown(position[user].supplyShares, market.totalSupplyAssets, market.totalSupplyShares);
    }

    function borrowAssets(address user) external view returns (uint256) {
        return _borrowAssets(user);
    }

    function availableLiquidity() external view returns (uint256) {
        return market.totalSupplyAssets - market.totalBorrowAssets;
    }

    function isHealthy(address user) external view returns (bool) {
        return _isHealthy(user);
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

    function _borrowAssets(address user) internal view returns (uint256) {
        return SharesMath.toAssetsUp(position[user].borrowShares, market.totalBorrowAssets, market.totalBorrowShares);
    }

    function _isHealthy(address user) internal view returns (bool) {
        Position storage userPosition = position[user];
        if (userPosition.borrowShares == 0) return true;

        uint256 price = oracle.price();
        if (price > type(uint128).max) revert OraclePriceTooLarge();

        uint256 collateralValue = uint256(userPosition.collateral) * price / WAD;
        uint256 maxBorrow = collateralValue * lltv / WAD;

        return _borrowAssets(user) <= maxBorrow;
    }

    function _toAssets(uint256 value) internal pure returns (uint128) {
        if (value > MAX_ASSETS) revert AmountTooLarge();
        return _toUint128(value);
    }

    function _toShares(uint256 value) internal pure returns (uint128) {
        if (value > MAX_SHARES) revert AmountTooLarge();
        return _toUint128(value);
    }

    function _toUint128(uint256 value) internal pure returns (uint128) {
        if (value > type(uint128).max) revert AmountTooLarge();
        // Safe because value was bounded above.
        // forge-lint: disable-next-line(unsafe-typecast)
        return uint128(value);
    }
}
