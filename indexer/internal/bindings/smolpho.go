// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SmolphoMetaData contains all meta data concerning the Smolpho contract.
var SmolphoMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"loanToken_\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"collateralToken_\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"oracle_\",\"type\":\"address\",\"internalType\":\"contractIPriceOracle\"},{\"name\":\"lltv_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ratePerSecond_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"liquidationIncentive_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"VIRTUAL_ASSETS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VIRTUAL_SHARES\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"WAD\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"accrueInterest\",\"inputs\":[],\"outputs\":[{\"name\":\"interest\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"availableLiquidity\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"borrow\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"borrowAssets\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"collateralToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isHealthy\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"liquidate\",\"inputs\":[{\"name\":\"borrower\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"repaidShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"repaidAssets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"seizedCollateral\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"badDebtAssets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"liquidationIncentive\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lltv\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"loanToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"market\",\"inputs\":[],\"outputs\":[{\"name\":\"totalSupplyAssets\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"totalSupplyShares\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"totalBorrowAssets\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"totalBorrowShares\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdate\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"oracle\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPriceOracle\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"position\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"supplyShares\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"borrowShares\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"collateral\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewBorrow\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewRepay\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewSupply\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewWithdraw\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ratePerSecond\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"repay\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supply\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supplyAssets\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"supplyCollateral\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawCollateral\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BadDebtRealized\",\"inputs\":[{\"name\":\"borrower\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"badDebtAssets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"badDebtShares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Borrowed\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollateralSupplied\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollateralWithdrawn\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InterestAccrued\",\"inputs\":[{\"name\":\"elapsed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"interest\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Liquidated\",\"inputs\":[{\"name\":\"liquidator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"borrower\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"repaidAssets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"repaidShares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"seizedCollateral\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Repaid\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Supplied\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawn\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AmountTooLarge\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"HealthyPosition\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientBorrowShares\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientCollateral\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientLiquidity\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientSupplyShares\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidLiquidationIncentive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidLltv\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OraclePriceTooLarge\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RateTooLarge\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Reentrancy\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SameToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TimestampTooLarge\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnhealthyPosition\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAssets\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroShares\",\"inputs\":[]}]",
}

// SmolphoABI is the input ABI used to generate the binding from.
// Deprecated: Use SmolphoMetaData.ABI instead.
var SmolphoABI = SmolphoMetaData.ABI

// Smolpho is an auto generated Go binding around an Ethereum contract.
type Smolpho struct {
	SmolphoCaller     // Read-only binding to the contract
	SmolphoTransactor // Write-only binding to the contract
	SmolphoFilterer   // Log filterer for contract events
}

// SmolphoCaller is an auto generated read-only Go binding around an Ethereum contract.
type SmolphoCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SmolphoTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SmolphoTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SmolphoFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SmolphoFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SmolphoSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SmolphoSession struct {
	Contract     *Smolpho          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SmolphoCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SmolphoCallerSession struct {
	Contract *SmolphoCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// SmolphoTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SmolphoTransactorSession struct {
	Contract     *SmolphoTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// SmolphoRaw is an auto generated low-level Go binding around an Ethereum contract.
type SmolphoRaw struct {
	Contract *Smolpho // Generic contract binding to access the raw methods on
}

// SmolphoCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SmolphoCallerRaw struct {
	Contract *SmolphoCaller // Generic read-only contract binding to access the raw methods on
}

// SmolphoTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SmolphoTransactorRaw struct {
	Contract *SmolphoTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSmolpho creates a new instance of Smolpho, bound to a specific deployed contract.
func NewSmolpho(address common.Address, backend bind.ContractBackend) (*Smolpho, error) {
	contract, err := bindSmolpho(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Smolpho{SmolphoCaller: SmolphoCaller{contract: contract}, SmolphoTransactor: SmolphoTransactor{contract: contract}, SmolphoFilterer: SmolphoFilterer{contract: contract}}, nil
}

// NewSmolphoCaller creates a new read-only instance of Smolpho, bound to a specific deployed contract.
func NewSmolphoCaller(address common.Address, caller bind.ContractCaller) (*SmolphoCaller, error) {
	contract, err := bindSmolpho(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SmolphoCaller{contract: contract}, nil
}

// NewSmolphoTransactor creates a new write-only instance of Smolpho, bound to a specific deployed contract.
func NewSmolphoTransactor(address common.Address, transactor bind.ContractTransactor) (*SmolphoTransactor, error) {
	contract, err := bindSmolpho(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SmolphoTransactor{contract: contract}, nil
}

// NewSmolphoFilterer creates a new log filterer instance of Smolpho, bound to a specific deployed contract.
func NewSmolphoFilterer(address common.Address, filterer bind.ContractFilterer) (*SmolphoFilterer, error) {
	contract, err := bindSmolpho(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SmolphoFilterer{contract: contract}, nil
}

// bindSmolpho binds a generic wrapper to an already deployed contract.
func bindSmolpho(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SmolphoABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Smolpho *SmolphoRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Smolpho.Contract.SmolphoCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Smolpho *SmolphoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Smolpho.Contract.SmolphoTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Smolpho *SmolphoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Smolpho.Contract.SmolphoTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Smolpho *SmolphoCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Smolpho.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Smolpho *SmolphoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Smolpho.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Smolpho *SmolphoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Smolpho.Contract.contract.Transact(opts, method, params...)
}

// VIRTUALASSETS is a free data retrieval call binding the contract method 0xb6608409.
//
// Solidity: function VIRTUAL_ASSETS() view returns(uint256)
func (_Smolpho *SmolphoCaller) VIRTUALASSETS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "VIRTUAL_ASSETS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VIRTUALASSETS is a free data retrieval call binding the contract method 0xb6608409.
//
// Solidity: function VIRTUAL_ASSETS() view returns(uint256)
func (_Smolpho *SmolphoSession) VIRTUALASSETS() (*big.Int, error) {
	return _Smolpho.Contract.VIRTUALASSETS(&_Smolpho.CallOpts)
}

// VIRTUALASSETS is a free data retrieval call binding the contract method 0xb6608409.
//
// Solidity: function VIRTUAL_ASSETS() view returns(uint256)
func (_Smolpho *SmolphoCallerSession) VIRTUALASSETS() (*big.Int, error) {
	return _Smolpho.Contract.VIRTUALASSETS(&_Smolpho.CallOpts)
}

// VIRTUALSHARES is a free data retrieval call binding the contract method 0x88c47f68.
//
// Solidity: function VIRTUAL_SHARES() view returns(uint256)
func (_Smolpho *SmolphoCaller) VIRTUALSHARES(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "VIRTUAL_SHARES")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VIRTUALSHARES is a free data retrieval call binding the contract method 0x88c47f68.
//
// Solidity: function VIRTUAL_SHARES() view returns(uint256)
func (_Smolpho *SmolphoSession) VIRTUALSHARES() (*big.Int, error) {
	return _Smolpho.Contract.VIRTUALSHARES(&_Smolpho.CallOpts)
}

// VIRTUALSHARES is a free data retrieval call binding the contract method 0x88c47f68.
//
// Solidity: function VIRTUAL_SHARES() view returns(uint256)
func (_Smolpho *SmolphoCallerSession) VIRTUALSHARES() (*big.Int, error) {
	return _Smolpho.Contract.VIRTUALSHARES(&_Smolpho.CallOpts)
}

// WAD is a free data retrieval call binding the contract method 0x6a146024.
//
// Solidity: function WAD() view returns(uint256)
func (_Smolpho *SmolphoCaller) WAD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "WAD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WAD is a free data retrieval call binding the contract method 0x6a146024.
//
// Solidity: function WAD() view returns(uint256)
func (_Smolpho *SmolphoSession) WAD() (*big.Int, error) {
	return _Smolpho.Contract.WAD(&_Smolpho.CallOpts)
}

// WAD is a free data retrieval call binding the contract method 0x6a146024.
//
// Solidity: function WAD() view returns(uint256)
func (_Smolpho *SmolphoCallerSession) WAD() (*big.Int, error) {
	return _Smolpho.Contract.WAD(&_Smolpho.CallOpts)
}

// AvailableLiquidity is a free data retrieval call binding the contract method 0x74375359.
//
// Solidity: function availableLiquidity() view returns(uint256)
func (_Smolpho *SmolphoCaller) AvailableLiquidity(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "availableLiquidity")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AvailableLiquidity is a free data retrieval call binding the contract method 0x74375359.
//
// Solidity: function availableLiquidity() view returns(uint256)
func (_Smolpho *SmolphoSession) AvailableLiquidity() (*big.Int, error) {
	return _Smolpho.Contract.AvailableLiquidity(&_Smolpho.CallOpts)
}

// AvailableLiquidity is a free data retrieval call binding the contract method 0x74375359.
//
// Solidity: function availableLiquidity() view returns(uint256)
func (_Smolpho *SmolphoCallerSession) AvailableLiquidity() (*big.Int, error) {
	return _Smolpho.Contract.AvailableLiquidity(&_Smolpho.CallOpts)
}

// BorrowAssets is a free data retrieval call binding the contract method 0xfa7594f4.
//
// Solidity: function borrowAssets(address user) view returns(uint256)
func (_Smolpho *SmolphoCaller) BorrowAssets(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "borrowAssets", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BorrowAssets is a free data retrieval call binding the contract method 0xfa7594f4.
//
// Solidity: function borrowAssets(address user) view returns(uint256)
func (_Smolpho *SmolphoSession) BorrowAssets(user common.Address) (*big.Int, error) {
	return _Smolpho.Contract.BorrowAssets(&_Smolpho.CallOpts, user)
}

// BorrowAssets is a free data retrieval call binding the contract method 0xfa7594f4.
//
// Solidity: function borrowAssets(address user) view returns(uint256)
func (_Smolpho *SmolphoCallerSession) BorrowAssets(user common.Address) (*big.Int, error) {
	return _Smolpho.Contract.BorrowAssets(&_Smolpho.CallOpts, user)
}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_Smolpho *SmolphoCaller) CollateralToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "collateralToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_Smolpho *SmolphoSession) CollateralToken() (common.Address, error) {
	return _Smolpho.Contract.CollateralToken(&_Smolpho.CallOpts)
}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_Smolpho *SmolphoCallerSession) CollateralToken() (common.Address, error) {
	return _Smolpho.Contract.CollateralToken(&_Smolpho.CallOpts)
}

// IsHealthy is a free data retrieval call binding the contract method 0xd4d9bc68.
//
// Solidity: function isHealthy(address user) view returns(bool)
func (_Smolpho *SmolphoCaller) IsHealthy(opts *bind.CallOpts, user common.Address) (bool, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "isHealthy", user)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsHealthy is a free data retrieval call binding the contract method 0xd4d9bc68.
//
// Solidity: function isHealthy(address user) view returns(bool)
func (_Smolpho *SmolphoSession) IsHealthy(user common.Address) (bool, error) {
	return _Smolpho.Contract.IsHealthy(&_Smolpho.CallOpts, user)
}

// IsHealthy is a free data retrieval call binding the contract method 0xd4d9bc68.
//
// Solidity: function isHealthy(address user) view returns(bool)
func (_Smolpho *SmolphoCallerSession) IsHealthy(user common.Address) (bool, error) {
	return _Smolpho.Contract.IsHealthy(&_Smolpho.CallOpts, user)
}

// LiquidationIncentive is a free data retrieval call binding the contract method 0x8c765e94.
//
// Solidity: function liquidationIncentive() view returns(uint256)
func (_Smolpho *SmolphoCaller) LiquidationIncentive(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "liquidationIncentive")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LiquidationIncentive is a free data retrieval call binding the contract method 0x8c765e94.
//
// Solidity: function liquidationIncentive() view returns(uint256)
func (_Smolpho *SmolphoSession) LiquidationIncentive() (*big.Int, error) {
	return _Smolpho.Contract.LiquidationIncentive(&_Smolpho.CallOpts)
}

// LiquidationIncentive is a free data retrieval call binding the contract method 0x8c765e94.
//
// Solidity: function liquidationIncentive() view returns(uint256)
func (_Smolpho *SmolphoCallerSession) LiquidationIncentive() (*big.Int, error) {
	return _Smolpho.Contract.LiquidationIncentive(&_Smolpho.CallOpts)
}

// Lltv is a free data retrieval call binding the contract method 0x217b7ffe.
//
// Solidity: function lltv() view returns(uint256)
func (_Smolpho *SmolphoCaller) Lltv(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "lltv")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Lltv is a free data retrieval call binding the contract method 0x217b7ffe.
//
// Solidity: function lltv() view returns(uint256)
func (_Smolpho *SmolphoSession) Lltv() (*big.Int, error) {
	return _Smolpho.Contract.Lltv(&_Smolpho.CallOpts)
}

// Lltv is a free data retrieval call binding the contract method 0x217b7ffe.
//
// Solidity: function lltv() view returns(uint256)
func (_Smolpho *SmolphoCallerSession) Lltv() (*big.Int, error) {
	return _Smolpho.Contract.Lltv(&_Smolpho.CallOpts)
}

// LoanToken is a free data retrieval call binding the contract method 0x06d37817.
//
// Solidity: function loanToken() view returns(address)
func (_Smolpho *SmolphoCaller) LoanToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "loanToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LoanToken is a free data retrieval call binding the contract method 0x06d37817.
//
// Solidity: function loanToken() view returns(address)
func (_Smolpho *SmolphoSession) LoanToken() (common.Address, error) {
	return _Smolpho.Contract.LoanToken(&_Smolpho.CallOpts)
}

// LoanToken is a free data retrieval call binding the contract method 0x06d37817.
//
// Solidity: function loanToken() view returns(address)
func (_Smolpho *SmolphoCallerSession) LoanToken() (common.Address, error) {
	return _Smolpho.Contract.LoanToken(&_Smolpho.CallOpts)
}

// Market is a free data retrieval call binding the contract method 0x80f55605.
//
// Solidity: function market() view returns(uint128 totalSupplyAssets, uint128 totalSupplyShares, uint128 totalBorrowAssets, uint128 totalBorrowShares, uint64 lastUpdate)
func (_Smolpho *SmolphoCaller) Market(opts *bind.CallOpts) (struct {
	TotalSupplyAssets *big.Int
	TotalSupplyShares *big.Int
	TotalBorrowAssets *big.Int
	TotalBorrowShares *big.Int
	LastUpdate        uint64
}, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "market")

	outstruct := new(struct {
		TotalSupplyAssets *big.Int
		TotalSupplyShares *big.Int
		TotalBorrowAssets *big.Int
		TotalBorrowShares *big.Int
		LastUpdate        uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalSupplyAssets = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalSupplyShares = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TotalBorrowAssets = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.TotalBorrowShares = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.LastUpdate = *abi.ConvertType(out[4], new(uint64)).(*uint64)

	return *outstruct, err

}

// Market is a free data retrieval call binding the contract method 0x80f55605.
//
// Solidity: function market() view returns(uint128 totalSupplyAssets, uint128 totalSupplyShares, uint128 totalBorrowAssets, uint128 totalBorrowShares, uint64 lastUpdate)
func (_Smolpho *SmolphoSession) Market() (struct {
	TotalSupplyAssets *big.Int
	TotalSupplyShares *big.Int
	TotalBorrowAssets *big.Int
	TotalBorrowShares *big.Int
	LastUpdate        uint64
}, error) {
	return _Smolpho.Contract.Market(&_Smolpho.CallOpts)
}

// Market is a free data retrieval call binding the contract method 0x80f55605.
//
// Solidity: function market() view returns(uint128 totalSupplyAssets, uint128 totalSupplyShares, uint128 totalBorrowAssets, uint128 totalBorrowShares, uint64 lastUpdate)
func (_Smolpho *SmolphoCallerSession) Market() (struct {
	TotalSupplyAssets *big.Int
	TotalSupplyShares *big.Int
	TotalBorrowAssets *big.Int
	TotalBorrowShares *big.Int
	LastUpdate        uint64
}, error) {
	return _Smolpho.Contract.Market(&_Smolpho.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_Smolpho *SmolphoCaller) Oracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "oracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_Smolpho *SmolphoSession) Oracle() (common.Address, error) {
	return _Smolpho.Contract.Oracle(&_Smolpho.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_Smolpho *SmolphoCallerSession) Oracle() (common.Address, error) {
	return _Smolpho.Contract.Oracle(&_Smolpho.CallOpts)
}

// Position is a free data retrieval call binding the contract method 0xb7648fb9.
//
// Solidity: function position(address ) view returns(uint128 supplyShares, uint128 borrowShares, uint128 collateral)
func (_Smolpho *SmolphoCaller) Position(opts *bind.CallOpts, arg0 common.Address) (struct {
	SupplyShares *big.Int
	BorrowShares *big.Int
	Collateral   *big.Int
}, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "position", arg0)

	outstruct := new(struct {
		SupplyShares *big.Int
		BorrowShares *big.Int
		Collateral   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SupplyShares = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.BorrowShares = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Collateral = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Position is a free data retrieval call binding the contract method 0xb7648fb9.
//
// Solidity: function position(address ) view returns(uint128 supplyShares, uint128 borrowShares, uint128 collateral)
func (_Smolpho *SmolphoSession) Position(arg0 common.Address) (struct {
	SupplyShares *big.Int
	BorrowShares *big.Int
	Collateral   *big.Int
}, error) {
	return _Smolpho.Contract.Position(&_Smolpho.CallOpts, arg0)
}

// Position is a free data retrieval call binding the contract method 0xb7648fb9.
//
// Solidity: function position(address ) view returns(uint128 supplyShares, uint128 borrowShares, uint128 collateral)
func (_Smolpho *SmolphoCallerSession) Position(arg0 common.Address) (struct {
	SupplyShares *big.Int
	BorrowShares *big.Int
	Collateral   *big.Int
}, error) {
	return _Smolpho.Contract.Position(&_Smolpho.CallOpts, arg0)
}

// PreviewBorrow is a free data retrieval call binding the contract method 0x78007e23.
//
// Solidity: function previewBorrow(uint256 assets) view returns(uint256)
func (_Smolpho *SmolphoCaller) PreviewBorrow(opts *bind.CallOpts, assets *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "previewBorrow", assets)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewBorrow is a free data retrieval call binding the contract method 0x78007e23.
//
// Solidity: function previewBorrow(uint256 assets) view returns(uint256)
func (_Smolpho *SmolphoSession) PreviewBorrow(assets *big.Int) (*big.Int, error) {
	return _Smolpho.Contract.PreviewBorrow(&_Smolpho.CallOpts, assets)
}

// PreviewBorrow is a free data retrieval call binding the contract method 0x78007e23.
//
// Solidity: function previewBorrow(uint256 assets) view returns(uint256)
func (_Smolpho *SmolphoCallerSession) PreviewBorrow(assets *big.Int) (*big.Int, error) {
	return _Smolpho.Contract.PreviewBorrow(&_Smolpho.CallOpts, assets)
}

// PreviewRepay is a free data retrieval call binding the contract method 0xbf722ca2.
//
// Solidity: function previewRepay(uint256 shares) view returns(uint256)
func (_Smolpho *SmolphoCaller) PreviewRepay(opts *bind.CallOpts, shares *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "previewRepay", shares)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewRepay is a free data retrieval call binding the contract method 0xbf722ca2.
//
// Solidity: function previewRepay(uint256 shares) view returns(uint256)
func (_Smolpho *SmolphoSession) PreviewRepay(shares *big.Int) (*big.Int, error) {
	return _Smolpho.Contract.PreviewRepay(&_Smolpho.CallOpts, shares)
}

// PreviewRepay is a free data retrieval call binding the contract method 0xbf722ca2.
//
// Solidity: function previewRepay(uint256 shares) view returns(uint256)
func (_Smolpho *SmolphoCallerSession) PreviewRepay(shares *big.Int) (*big.Int, error) {
	return _Smolpho.Contract.PreviewRepay(&_Smolpho.CallOpts, shares)
}

// PreviewSupply is a free data retrieval call binding the contract method 0xc999d906.
//
// Solidity: function previewSupply(uint256 assets) view returns(uint256)
func (_Smolpho *SmolphoCaller) PreviewSupply(opts *bind.CallOpts, assets *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "previewSupply", assets)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewSupply is a free data retrieval call binding the contract method 0xc999d906.
//
// Solidity: function previewSupply(uint256 assets) view returns(uint256)
func (_Smolpho *SmolphoSession) PreviewSupply(assets *big.Int) (*big.Int, error) {
	return _Smolpho.Contract.PreviewSupply(&_Smolpho.CallOpts, assets)
}

// PreviewSupply is a free data retrieval call binding the contract method 0xc999d906.
//
// Solidity: function previewSupply(uint256 assets) view returns(uint256)
func (_Smolpho *SmolphoCallerSession) PreviewSupply(assets *big.Int) (*big.Int, error) {
	return _Smolpho.Contract.PreviewSupply(&_Smolpho.CallOpts, assets)
}

// PreviewWithdraw is a free data retrieval call binding the contract method 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 shares) view returns(uint256)
func (_Smolpho *SmolphoCaller) PreviewWithdraw(opts *bind.CallOpts, shares *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "previewWithdraw", shares)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewWithdraw is a free data retrieval call binding the contract method 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 shares) view returns(uint256)
func (_Smolpho *SmolphoSession) PreviewWithdraw(shares *big.Int) (*big.Int, error) {
	return _Smolpho.Contract.PreviewWithdraw(&_Smolpho.CallOpts, shares)
}

// PreviewWithdraw is a free data retrieval call binding the contract method 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 shares) view returns(uint256)
func (_Smolpho *SmolphoCallerSession) PreviewWithdraw(shares *big.Int) (*big.Int, error) {
	return _Smolpho.Contract.PreviewWithdraw(&_Smolpho.CallOpts, shares)
}

// RatePerSecond is a free data retrieval call binding the contract method 0x8eff1a98.
//
// Solidity: function ratePerSecond() view returns(uint64)
func (_Smolpho *SmolphoCaller) RatePerSecond(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "ratePerSecond")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// RatePerSecond is a free data retrieval call binding the contract method 0x8eff1a98.
//
// Solidity: function ratePerSecond() view returns(uint64)
func (_Smolpho *SmolphoSession) RatePerSecond() (uint64, error) {
	return _Smolpho.Contract.RatePerSecond(&_Smolpho.CallOpts)
}

// RatePerSecond is a free data retrieval call binding the contract method 0x8eff1a98.
//
// Solidity: function ratePerSecond() view returns(uint64)
func (_Smolpho *SmolphoCallerSession) RatePerSecond() (uint64, error) {
	return _Smolpho.Contract.RatePerSecond(&_Smolpho.CallOpts)
}

// SupplyAssets is a free data retrieval call binding the contract method 0x262a0a0e.
//
// Solidity: function supplyAssets(address user) view returns(uint256)
func (_Smolpho *SmolphoCaller) SupplyAssets(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Smolpho.contract.Call(opts, &out, "supplyAssets", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SupplyAssets is a free data retrieval call binding the contract method 0x262a0a0e.
//
// Solidity: function supplyAssets(address user) view returns(uint256)
func (_Smolpho *SmolphoSession) SupplyAssets(user common.Address) (*big.Int, error) {
	return _Smolpho.Contract.SupplyAssets(&_Smolpho.CallOpts, user)
}

// SupplyAssets is a free data retrieval call binding the contract method 0x262a0a0e.
//
// Solidity: function supplyAssets(address user) view returns(uint256)
func (_Smolpho *SmolphoCallerSession) SupplyAssets(user common.Address) (*big.Int, error) {
	return _Smolpho.Contract.SupplyAssets(&_Smolpho.CallOpts, user)
}

// AccrueInterest is a paid mutator transaction binding the contract method 0xa6afed95.
//
// Solidity: function accrueInterest() returns(uint256 interest)
func (_Smolpho *SmolphoTransactor) AccrueInterest(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Smolpho.contract.Transact(opts, "accrueInterest")
}

// AccrueInterest is a paid mutator transaction binding the contract method 0xa6afed95.
//
// Solidity: function accrueInterest() returns(uint256 interest)
func (_Smolpho *SmolphoSession) AccrueInterest() (*types.Transaction, error) {
	return _Smolpho.Contract.AccrueInterest(&_Smolpho.TransactOpts)
}

// AccrueInterest is a paid mutator transaction binding the contract method 0xa6afed95.
//
// Solidity: function accrueInterest() returns(uint256 interest)
func (_Smolpho *SmolphoTransactorSession) AccrueInterest() (*types.Transaction, error) {
	return _Smolpho.Contract.AccrueInterest(&_Smolpho.TransactOpts)
}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 assets) returns(uint256 shares)
func (_Smolpho *SmolphoTransactor) Borrow(opts *bind.TransactOpts, assets *big.Int) (*types.Transaction, error) {
	return _Smolpho.contract.Transact(opts, "borrow", assets)
}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 assets) returns(uint256 shares)
func (_Smolpho *SmolphoSession) Borrow(assets *big.Int) (*types.Transaction, error) {
	return _Smolpho.Contract.Borrow(&_Smolpho.TransactOpts, assets)
}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 assets) returns(uint256 shares)
func (_Smolpho *SmolphoTransactorSession) Borrow(assets *big.Int) (*types.Transaction, error) {
	return _Smolpho.Contract.Borrow(&_Smolpho.TransactOpts, assets)
}

// Liquidate is a paid mutator transaction binding the contract method 0xbcbaf487.
//
// Solidity: function liquidate(address borrower, uint256 repaidShares) returns(uint256 repaidAssets, uint256 seizedCollateral, uint256 badDebtAssets)
func (_Smolpho *SmolphoTransactor) Liquidate(opts *bind.TransactOpts, borrower common.Address, repaidShares *big.Int) (*types.Transaction, error) {
	return _Smolpho.contract.Transact(opts, "liquidate", borrower, repaidShares)
}

// Liquidate is a paid mutator transaction binding the contract method 0xbcbaf487.
//
// Solidity: function liquidate(address borrower, uint256 repaidShares) returns(uint256 repaidAssets, uint256 seizedCollateral, uint256 badDebtAssets)
func (_Smolpho *SmolphoSession) Liquidate(borrower common.Address, repaidShares *big.Int) (*types.Transaction, error) {
	return _Smolpho.Contract.Liquidate(&_Smolpho.TransactOpts, borrower, repaidShares)
}

// Liquidate is a paid mutator transaction binding the contract method 0xbcbaf487.
//
// Solidity: function liquidate(address borrower, uint256 repaidShares) returns(uint256 repaidAssets, uint256 seizedCollateral, uint256 badDebtAssets)
func (_Smolpho *SmolphoTransactorSession) Liquidate(borrower common.Address, repaidShares *big.Int) (*types.Transaction, error) {
	return _Smolpho.Contract.Liquidate(&_Smolpho.TransactOpts, borrower, repaidShares)
}

// Repay is a paid mutator transaction binding the contract method 0x371fd8e6.
//
// Solidity: function repay(uint256 shares) returns(uint256 assets)
func (_Smolpho *SmolphoTransactor) Repay(opts *bind.TransactOpts, shares *big.Int) (*types.Transaction, error) {
	return _Smolpho.contract.Transact(opts, "repay", shares)
}

// Repay is a paid mutator transaction binding the contract method 0x371fd8e6.
//
// Solidity: function repay(uint256 shares) returns(uint256 assets)
func (_Smolpho *SmolphoSession) Repay(shares *big.Int) (*types.Transaction, error) {
	return _Smolpho.Contract.Repay(&_Smolpho.TransactOpts, shares)
}

// Repay is a paid mutator transaction binding the contract method 0x371fd8e6.
//
// Solidity: function repay(uint256 shares) returns(uint256 assets)
func (_Smolpho *SmolphoTransactorSession) Repay(shares *big.Int) (*types.Transaction, error) {
	return _Smolpho.Contract.Repay(&_Smolpho.TransactOpts, shares)
}

// Supply is a paid mutator transaction binding the contract method 0x35403023.
//
// Solidity: function supply(uint256 assets) returns(uint256 shares)
func (_Smolpho *SmolphoTransactor) Supply(opts *bind.TransactOpts, assets *big.Int) (*types.Transaction, error) {
	return _Smolpho.contract.Transact(opts, "supply", assets)
}

// Supply is a paid mutator transaction binding the contract method 0x35403023.
//
// Solidity: function supply(uint256 assets) returns(uint256 shares)
func (_Smolpho *SmolphoSession) Supply(assets *big.Int) (*types.Transaction, error) {
	return _Smolpho.Contract.Supply(&_Smolpho.TransactOpts, assets)
}

// Supply is a paid mutator transaction binding the contract method 0x35403023.
//
// Solidity: function supply(uint256 assets) returns(uint256 shares)
func (_Smolpho *SmolphoTransactorSession) Supply(assets *big.Int) (*types.Transaction, error) {
	return _Smolpho.Contract.Supply(&_Smolpho.TransactOpts, assets)
}

// SupplyCollateral is a paid mutator transaction binding the contract method 0x367febea.
//
// Solidity: function supplyCollateral(uint256 assets) returns()
func (_Smolpho *SmolphoTransactor) SupplyCollateral(opts *bind.TransactOpts, assets *big.Int) (*types.Transaction, error) {
	return _Smolpho.contract.Transact(opts, "supplyCollateral", assets)
}

// SupplyCollateral is a paid mutator transaction binding the contract method 0x367febea.
//
// Solidity: function supplyCollateral(uint256 assets) returns()
func (_Smolpho *SmolphoSession) SupplyCollateral(assets *big.Int) (*types.Transaction, error) {
	return _Smolpho.Contract.SupplyCollateral(&_Smolpho.TransactOpts, assets)
}

// SupplyCollateral is a paid mutator transaction binding the contract method 0x367febea.
//
// Solidity: function supplyCollateral(uint256 assets) returns()
func (_Smolpho *SmolphoTransactorSession) SupplyCollateral(assets *big.Int) (*types.Transaction, error) {
	return _Smolpho.Contract.SupplyCollateral(&_Smolpho.TransactOpts, assets)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 shares) returns(uint256 assets)
func (_Smolpho *SmolphoTransactor) Withdraw(opts *bind.TransactOpts, shares *big.Int) (*types.Transaction, error) {
	return _Smolpho.contract.Transact(opts, "withdraw", shares)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 shares) returns(uint256 assets)
func (_Smolpho *SmolphoSession) Withdraw(shares *big.Int) (*types.Transaction, error) {
	return _Smolpho.Contract.Withdraw(&_Smolpho.TransactOpts, shares)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 shares) returns(uint256 assets)
func (_Smolpho *SmolphoTransactorSession) Withdraw(shares *big.Int) (*types.Transaction, error) {
	return _Smolpho.Contract.Withdraw(&_Smolpho.TransactOpts, shares)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x6112fe2e.
//
// Solidity: function withdrawCollateral(uint256 assets) returns()
func (_Smolpho *SmolphoTransactor) WithdrawCollateral(opts *bind.TransactOpts, assets *big.Int) (*types.Transaction, error) {
	return _Smolpho.contract.Transact(opts, "withdrawCollateral", assets)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x6112fe2e.
//
// Solidity: function withdrawCollateral(uint256 assets) returns()
func (_Smolpho *SmolphoSession) WithdrawCollateral(assets *big.Int) (*types.Transaction, error) {
	return _Smolpho.Contract.WithdrawCollateral(&_Smolpho.TransactOpts, assets)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x6112fe2e.
//
// Solidity: function withdrawCollateral(uint256 assets) returns()
func (_Smolpho *SmolphoTransactorSession) WithdrawCollateral(assets *big.Int) (*types.Transaction, error) {
	return _Smolpho.Contract.WithdrawCollateral(&_Smolpho.TransactOpts, assets)
}

// SmolphoBadDebtRealizedIterator is returned from FilterBadDebtRealized and is used to iterate over the raw logs and unpacked data for BadDebtRealized events raised by the Smolpho contract.
type SmolphoBadDebtRealizedIterator struct {
	Event *SmolphoBadDebtRealized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SmolphoBadDebtRealizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SmolphoBadDebtRealized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SmolphoBadDebtRealized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SmolphoBadDebtRealizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SmolphoBadDebtRealizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SmolphoBadDebtRealized represents a BadDebtRealized event raised by the Smolpho contract.
type SmolphoBadDebtRealized struct {
	Borrower      common.Address
	BadDebtAssets *big.Int
	BadDebtShares *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterBadDebtRealized is a free log retrieval operation binding the contract event 0xd14612560c7661a79aa4084251af048c123c29032771c1c57b30f592173bbef8.
//
// Solidity: event BadDebtRealized(address indexed borrower, uint256 badDebtAssets, uint256 badDebtShares)
func (_Smolpho *SmolphoFilterer) FilterBadDebtRealized(opts *bind.FilterOpts, borrower []common.Address) (*SmolphoBadDebtRealizedIterator, error) {

	var borrowerRule []interface{}
	for _, borrowerItem := range borrower {
		borrowerRule = append(borrowerRule, borrowerItem)
	}

	logs, sub, err := _Smolpho.contract.FilterLogs(opts, "BadDebtRealized", borrowerRule)
	if err != nil {
		return nil, err
	}
	return &SmolphoBadDebtRealizedIterator{contract: _Smolpho.contract, event: "BadDebtRealized", logs: logs, sub: sub}, nil
}

// WatchBadDebtRealized is a free log subscription operation binding the contract event 0xd14612560c7661a79aa4084251af048c123c29032771c1c57b30f592173bbef8.
//
// Solidity: event BadDebtRealized(address indexed borrower, uint256 badDebtAssets, uint256 badDebtShares)
func (_Smolpho *SmolphoFilterer) WatchBadDebtRealized(opts *bind.WatchOpts, sink chan<- *SmolphoBadDebtRealized, borrower []common.Address) (event.Subscription, error) {

	var borrowerRule []interface{}
	for _, borrowerItem := range borrower {
		borrowerRule = append(borrowerRule, borrowerItem)
	}

	logs, sub, err := _Smolpho.contract.WatchLogs(opts, "BadDebtRealized", borrowerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SmolphoBadDebtRealized)
				if err := _Smolpho.contract.UnpackLog(event, "BadDebtRealized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBadDebtRealized is a log parse operation binding the contract event 0xd14612560c7661a79aa4084251af048c123c29032771c1c57b30f592173bbef8.
//
// Solidity: event BadDebtRealized(address indexed borrower, uint256 badDebtAssets, uint256 badDebtShares)
func (_Smolpho *SmolphoFilterer) ParseBadDebtRealized(log types.Log) (*SmolphoBadDebtRealized, error) {
	event := new(SmolphoBadDebtRealized)
	if err := _Smolpho.contract.UnpackLog(event, "BadDebtRealized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SmolphoBorrowedIterator is returned from FilterBorrowed and is used to iterate over the raw logs and unpacked data for Borrowed events raised by the Smolpho contract.
type SmolphoBorrowedIterator struct {
	Event *SmolphoBorrowed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SmolphoBorrowedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SmolphoBorrowed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SmolphoBorrowed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SmolphoBorrowedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SmolphoBorrowedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SmolphoBorrowed represents a Borrowed event raised by the Smolpho contract.
type SmolphoBorrowed struct {
	User   common.Address
	Assets *big.Int
	Shares *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBorrowed is a free log retrieval operation binding the contract event 0xeae9cfbc77fdd40ca899f36b608256063b2bc9d8178b0220f7ad513e178d6730.
//
// Solidity: event Borrowed(address indexed user, uint256 assets, uint256 shares)
func (_Smolpho *SmolphoFilterer) FilterBorrowed(opts *bind.FilterOpts, user []common.Address) (*SmolphoBorrowedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Smolpho.contract.FilterLogs(opts, "Borrowed", userRule)
	if err != nil {
		return nil, err
	}
	return &SmolphoBorrowedIterator{contract: _Smolpho.contract, event: "Borrowed", logs: logs, sub: sub}, nil
}

// WatchBorrowed is a free log subscription operation binding the contract event 0xeae9cfbc77fdd40ca899f36b608256063b2bc9d8178b0220f7ad513e178d6730.
//
// Solidity: event Borrowed(address indexed user, uint256 assets, uint256 shares)
func (_Smolpho *SmolphoFilterer) WatchBorrowed(opts *bind.WatchOpts, sink chan<- *SmolphoBorrowed, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Smolpho.contract.WatchLogs(opts, "Borrowed", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SmolphoBorrowed)
				if err := _Smolpho.contract.UnpackLog(event, "Borrowed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBorrowed is a log parse operation binding the contract event 0xeae9cfbc77fdd40ca899f36b608256063b2bc9d8178b0220f7ad513e178d6730.
//
// Solidity: event Borrowed(address indexed user, uint256 assets, uint256 shares)
func (_Smolpho *SmolphoFilterer) ParseBorrowed(log types.Log) (*SmolphoBorrowed, error) {
	event := new(SmolphoBorrowed)
	if err := _Smolpho.contract.UnpackLog(event, "Borrowed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SmolphoCollateralSuppliedIterator is returned from FilterCollateralSupplied and is used to iterate over the raw logs and unpacked data for CollateralSupplied events raised by the Smolpho contract.
type SmolphoCollateralSuppliedIterator struct {
	Event *SmolphoCollateralSupplied // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SmolphoCollateralSuppliedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SmolphoCollateralSupplied)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SmolphoCollateralSupplied)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SmolphoCollateralSuppliedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SmolphoCollateralSuppliedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SmolphoCollateralSupplied represents a CollateralSupplied event raised by the Smolpho contract.
type SmolphoCollateralSupplied struct {
	User   common.Address
	Assets *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCollateralSupplied is a free log retrieval operation binding the contract event 0xb0c1a992a318d3f9e5ee4ef9bce6d9310f55f81d40dd18429c1b4ad5aca3d0d1.
//
// Solidity: event CollateralSupplied(address indexed user, uint256 assets)
func (_Smolpho *SmolphoFilterer) FilterCollateralSupplied(opts *bind.FilterOpts, user []common.Address) (*SmolphoCollateralSuppliedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Smolpho.contract.FilterLogs(opts, "CollateralSupplied", userRule)
	if err != nil {
		return nil, err
	}
	return &SmolphoCollateralSuppliedIterator{contract: _Smolpho.contract, event: "CollateralSupplied", logs: logs, sub: sub}, nil
}

// WatchCollateralSupplied is a free log subscription operation binding the contract event 0xb0c1a992a318d3f9e5ee4ef9bce6d9310f55f81d40dd18429c1b4ad5aca3d0d1.
//
// Solidity: event CollateralSupplied(address indexed user, uint256 assets)
func (_Smolpho *SmolphoFilterer) WatchCollateralSupplied(opts *bind.WatchOpts, sink chan<- *SmolphoCollateralSupplied, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Smolpho.contract.WatchLogs(opts, "CollateralSupplied", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SmolphoCollateralSupplied)
				if err := _Smolpho.contract.UnpackLog(event, "CollateralSupplied", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCollateralSupplied is a log parse operation binding the contract event 0xb0c1a992a318d3f9e5ee4ef9bce6d9310f55f81d40dd18429c1b4ad5aca3d0d1.
//
// Solidity: event CollateralSupplied(address indexed user, uint256 assets)
func (_Smolpho *SmolphoFilterer) ParseCollateralSupplied(log types.Log) (*SmolphoCollateralSupplied, error) {
	event := new(SmolphoCollateralSupplied)
	if err := _Smolpho.contract.UnpackLog(event, "CollateralSupplied", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SmolphoCollateralWithdrawnIterator is returned from FilterCollateralWithdrawn and is used to iterate over the raw logs and unpacked data for CollateralWithdrawn events raised by the Smolpho contract.
type SmolphoCollateralWithdrawnIterator struct {
	Event *SmolphoCollateralWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SmolphoCollateralWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SmolphoCollateralWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SmolphoCollateralWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SmolphoCollateralWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SmolphoCollateralWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SmolphoCollateralWithdrawn represents a CollateralWithdrawn event raised by the Smolpho contract.
type SmolphoCollateralWithdrawn struct {
	User   common.Address
	Assets *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCollateralWithdrawn is a free log retrieval operation binding the contract event 0xc30fcfbcaac9e0deffa719714eaa82396ff506a0d0d0eebe170830177288715d.
//
// Solidity: event CollateralWithdrawn(address indexed user, uint256 assets)
func (_Smolpho *SmolphoFilterer) FilterCollateralWithdrawn(opts *bind.FilterOpts, user []common.Address) (*SmolphoCollateralWithdrawnIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Smolpho.contract.FilterLogs(opts, "CollateralWithdrawn", userRule)
	if err != nil {
		return nil, err
	}
	return &SmolphoCollateralWithdrawnIterator{contract: _Smolpho.contract, event: "CollateralWithdrawn", logs: logs, sub: sub}, nil
}

// WatchCollateralWithdrawn is a free log subscription operation binding the contract event 0xc30fcfbcaac9e0deffa719714eaa82396ff506a0d0d0eebe170830177288715d.
//
// Solidity: event CollateralWithdrawn(address indexed user, uint256 assets)
func (_Smolpho *SmolphoFilterer) WatchCollateralWithdrawn(opts *bind.WatchOpts, sink chan<- *SmolphoCollateralWithdrawn, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Smolpho.contract.WatchLogs(opts, "CollateralWithdrawn", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SmolphoCollateralWithdrawn)
				if err := _Smolpho.contract.UnpackLog(event, "CollateralWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCollateralWithdrawn is a log parse operation binding the contract event 0xc30fcfbcaac9e0deffa719714eaa82396ff506a0d0d0eebe170830177288715d.
//
// Solidity: event CollateralWithdrawn(address indexed user, uint256 assets)
func (_Smolpho *SmolphoFilterer) ParseCollateralWithdrawn(log types.Log) (*SmolphoCollateralWithdrawn, error) {
	event := new(SmolphoCollateralWithdrawn)
	if err := _Smolpho.contract.UnpackLog(event, "CollateralWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SmolphoInterestAccruedIterator is returned from FilterInterestAccrued and is used to iterate over the raw logs and unpacked data for InterestAccrued events raised by the Smolpho contract.
type SmolphoInterestAccruedIterator struct {
	Event *SmolphoInterestAccrued // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SmolphoInterestAccruedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SmolphoInterestAccrued)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SmolphoInterestAccrued)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SmolphoInterestAccruedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SmolphoInterestAccruedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SmolphoInterestAccrued represents a InterestAccrued event raised by the Smolpho contract.
type SmolphoInterestAccrued struct {
	Elapsed  *big.Int
	Interest *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterInterestAccrued is a free log retrieval operation binding the contract event 0x3d8708181c6452c5711a2b4c6ed3f12149f6ea35f78549aba00a89d041e766ca.
//
// Solidity: event InterestAccrued(uint256 elapsed, uint256 interest)
func (_Smolpho *SmolphoFilterer) FilterInterestAccrued(opts *bind.FilterOpts) (*SmolphoInterestAccruedIterator, error) {

	logs, sub, err := _Smolpho.contract.FilterLogs(opts, "InterestAccrued")
	if err != nil {
		return nil, err
	}
	return &SmolphoInterestAccruedIterator{contract: _Smolpho.contract, event: "InterestAccrued", logs: logs, sub: sub}, nil
}

// WatchInterestAccrued is a free log subscription operation binding the contract event 0x3d8708181c6452c5711a2b4c6ed3f12149f6ea35f78549aba00a89d041e766ca.
//
// Solidity: event InterestAccrued(uint256 elapsed, uint256 interest)
func (_Smolpho *SmolphoFilterer) WatchInterestAccrued(opts *bind.WatchOpts, sink chan<- *SmolphoInterestAccrued) (event.Subscription, error) {

	logs, sub, err := _Smolpho.contract.WatchLogs(opts, "InterestAccrued")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SmolphoInterestAccrued)
				if err := _Smolpho.contract.UnpackLog(event, "InterestAccrued", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInterestAccrued is a log parse operation binding the contract event 0x3d8708181c6452c5711a2b4c6ed3f12149f6ea35f78549aba00a89d041e766ca.
//
// Solidity: event InterestAccrued(uint256 elapsed, uint256 interest)
func (_Smolpho *SmolphoFilterer) ParseInterestAccrued(log types.Log) (*SmolphoInterestAccrued, error) {
	event := new(SmolphoInterestAccrued)
	if err := _Smolpho.contract.UnpackLog(event, "InterestAccrued", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SmolphoLiquidatedIterator is returned from FilterLiquidated and is used to iterate over the raw logs and unpacked data for Liquidated events raised by the Smolpho contract.
type SmolphoLiquidatedIterator struct {
	Event *SmolphoLiquidated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SmolphoLiquidatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SmolphoLiquidated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SmolphoLiquidated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SmolphoLiquidatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SmolphoLiquidatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SmolphoLiquidated represents a Liquidated event raised by the Smolpho contract.
type SmolphoLiquidated struct {
	Liquidator       common.Address
	Borrower         common.Address
	RepaidAssets     *big.Int
	RepaidShares     *big.Int
	SeizedCollateral *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLiquidated is a free log retrieval operation binding the contract event 0xfcbc974bf3a532baf2bb229db3c37fd58299b62d2d1db6a855dac5b693bb6ff3.
//
// Solidity: event Liquidated(address indexed liquidator, address indexed borrower, uint256 repaidAssets, uint256 repaidShares, uint256 seizedCollateral)
func (_Smolpho *SmolphoFilterer) FilterLiquidated(opts *bind.FilterOpts, liquidator []common.Address, borrower []common.Address) (*SmolphoLiquidatedIterator, error) {

	var liquidatorRule []interface{}
	for _, liquidatorItem := range liquidator {
		liquidatorRule = append(liquidatorRule, liquidatorItem)
	}
	var borrowerRule []interface{}
	for _, borrowerItem := range borrower {
		borrowerRule = append(borrowerRule, borrowerItem)
	}

	logs, sub, err := _Smolpho.contract.FilterLogs(opts, "Liquidated", liquidatorRule, borrowerRule)
	if err != nil {
		return nil, err
	}
	return &SmolphoLiquidatedIterator{contract: _Smolpho.contract, event: "Liquidated", logs: logs, sub: sub}, nil
}

// WatchLiquidated is a free log subscription operation binding the contract event 0xfcbc974bf3a532baf2bb229db3c37fd58299b62d2d1db6a855dac5b693bb6ff3.
//
// Solidity: event Liquidated(address indexed liquidator, address indexed borrower, uint256 repaidAssets, uint256 repaidShares, uint256 seizedCollateral)
func (_Smolpho *SmolphoFilterer) WatchLiquidated(opts *bind.WatchOpts, sink chan<- *SmolphoLiquidated, liquidator []common.Address, borrower []common.Address) (event.Subscription, error) {

	var liquidatorRule []interface{}
	for _, liquidatorItem := range liquidator {
		liquidatorRule = append(liquidatorRule, liquidatorItem)
	}
	var borrowerRule []interface{}
	for _, borrowerItem := range borrower {
		borrowerRule = append(borrowerRule, borrowerItem)
	}

	logs, sub, err := _Smolpho.contract.WatchLogs(opts, "Liquidated", liquidatorRule, borrowerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SmolphoLiquidated)
				if err := _Smolpho.contract.UnpackLog(event, "Liquidated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLiquidated is a log parse operation binding the contract event 0xfcbc974bf3a532baf2bb229db3c37fd58299b62d2d1db6a855dac5b693bb6ff3.
//
// Solidity: event Liquidated(address indexed liquidator, address indexed borrower, uint256 repaidAssets, uint256 repaidShares, uint256 seizedCollateral)
func (_Smolpho *SmolphoFilterer) ParseLiquidated(log types.Log) (*SmolphoLiquidated, error) {
	event := new(SmolphoLiquidated)
	if err := _Smolpho.contract.UnpackLog(event, "Liquidated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SmolphoRepaidIterator is returned from FilterRepaid and is used to iterate over the raw logs and unpacked data for Repaid events raised by the Smolpho contract.
type SmolphoRepaidIterator struct {
	Event *SmolphoRepaid // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SmolphoRepaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SmolphoRepaid)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SmolphoRepaid)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SmolphoRepaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SmolphoRepaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SmolphoRepaid represents a Repaid event raised by the Smolpho contract.
type SmolphoRepaid struct {
	User   common.Address
	Assets *big.Int
	Shares *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRepaid is a free log retrieval operation binding the contract event 0x1b8cd61ed43bec7c6bdad3a18ffee613f99c853d16c50678d248d879e1b43438.
//
// Solidity: event Repaid(address indexed user, uint256 assets, uint256 shares)
func (_Smolpho *SmolphoFilterer) FilterRepaid(opts *bind.FilterOpts, user []common.Address) (*SmolphoRepaidIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Smolpho.contract.FilterLogs(opts, "Repaid", userRule)
	if err != nil {
		return nil, err
	}
	return &SmolphoRepaidIterator{contract: _Smolpho.contract, event: "Repaid", logs: logs, sub: sub}, nil
}

// WatchRepaid is a free log subscription operation binding the contract event 0x1b8cd61ed43bec7c6bdad3a18ffee613f99c853d16c50678d248d879e1b43438.
//
// Solidity: event Repaid(address indexed user, uint256 assets, uint256 shares)
func (_Smolpho *SmolphoFilterer) WatchRepaid(opts *bind.WatchOpts, sink chan<- *SmolphoRepaid, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Smolpho.contract.WatchLogs(opts, "Repaid", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SmolphoRepaid)
				if err := _Smolpho.contract.UnpackLog(event, "Repaid", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRepaid is a log parse operation binding the contract event 0x1b8cd61ed43bec7c6bdad3a18ffee613f99c853d16c50678d248d879e1b43438.
//
// Solidity: event Repaid(address indexed user, uint256 assets, uint256 shares)
func (_Smolpho *SmolphoFilterer) ParseRepaid(log types.Log) (*SmolphoRepaid, error) {
	event := new(SmolphoRepaid)
	if err := _Smolpho.contract.UnpackLog(event, "Repaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SmolphoSuppliedIterator is returned from FilterSupplied and is used to iterate over the raw logs and unpacked data for Supplied events raised by the Smolpho contract.
type SmolphoSuppliedIterator struct {
	Event *SmolphoSupplied // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SmolphoSuppliedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SmolphoSupplied)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SmolphoSupplied)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SmolphoSuppliedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SmolphoSuppliedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SmolphoSupplied represents a Supplied event raised by the Smolpho contract.
type SmolphoSupplied struct {
	User   common.Address
	Assets *big.Int
	Shares *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterSupplied is a free log retrieval operation binding the contract event 0x5c2c0d2616a06b35bb159b4d7e227972b59bb33f8d5229ca0e5e438259bfd5d3.
//
// Solidity: event Supplied(address indexed user, uint256 assets, uint256 shares)
func (_Smolpho *SmolphoFilterer) FilterSupplied(opts *bind.FilterOpts, user []common.Address) (*SmolphoSuppliedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Smolpho.contract.FilterLogs(opts, "Supplied", userRule)
	if err != nil {
		return nil, err
	}
	return &SmolphoSuppliedIterator{contract: _Smolpho.contract, event: "Supplied", logs: logs, sub: sub}, nil
}

// WatchSupplied is a free log subscription operation binding the contract event 0x5c2c0d2616a06b35bb159b4d7e227972b59bb33f8d5229ca0e5e438259bfd5d3.
//
// Solidity: event Supplied(address indexed user, uint256 assets, uint256 shares)
func (_Smolpho *SmolphoFilterer) WatchSupplied(opts *bind.WatchOpts, sink chan<- *SmolphoSupplied, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Smolpho.contract.WatchLogs(opts, "Supplied", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SmolphoSupplied)
				if err := _Smolpho.contract.UnpackLog(event, "Supplied", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSupplied is a log parse operation binding the contract event 0x5c2c0d2616a06b35bb159b4d7e227972b59bb33f8d5229ca0e5e438259bfd5d3.
//
// Solidity: event Supplied(address indexed user, uint256 assets, uint256 shares)
func (_Smolpho *SmolphoFilterer) ParseSupplied(log types.Log) (*SmolphoSupplied, error) {
	event := new(SmolphoSupplied)
	if err := _Smolpho.contract.UnpackLog(event, "Supplied", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SmolphoWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the Smolpho contract.
type SmolphoWithdrawnIterator struct {
	Event *SmolphoWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SmolphoWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SmolphoWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SmolphoWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SmolphoWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SmolphoWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SmolphoWithdrawn represents a Withdrawn event raised by the Smolpho contract.
type SmolphoWithdrawn struct {
	User   common.Address
	Assets *big.Int
	Shares *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x92ccf450a286a957af52509bc1c9939d1a6a481783e142e41e2499f0bb66ebc6.
//
// Solidity: event Withdrawn(address indexed user, uint256 assets, uint256 shares)
func (_Smolpho *SmolphoFilterer) FilterWithdrawn(opts *bind.FilterOpts, user []common.Address) (*SmolphoWithdrawnIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Smolpho.contract.FilterLogs(opts, "Withdrawn", userRule)
	if err != nil {
		return nil, err
	}
	return &SmolphoWithdrawnIterator{contract: _Smolpho.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0x92ccf450a286a957af52509bc1c9939d1a6a481783e142e41e2499f0bb66ebc6.
//
// Solidity: event Withdrawn(address indexed user, uint256 assets, uint256 shares)
func (_Smolpho *SmolphoFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *SmolphoWithdrawn, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Smolpho.contract.WatchLogs(opts, "Withdrawn", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SmolphoWithdrawn)
				if err := _Smolpho.contract.UnpackLog(event, "Withdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawn is a log parse operation binding the contract event 0x92ccf450a286a957af52509bc1c9939d1a6a481783e142e41e2499f0bb66ebc6.
//
// Solidity: event Withdrawn(address indexed user, uint256 assets, uint256 shares)
func (_Smolpho *SmolphoFilterer) ParseWithdrawn(log types.Log) (*SmolphoWithdrawn, error) {
	event := new(SmolphoWithdrawn)
	if err := _Smolpho.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
