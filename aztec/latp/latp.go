// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package latp

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
	_ = abi.ConvertType
)

// LatpMetaData contains all meta data concerning the Latp contract.
var LatpMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getClaimed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"getClaimable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}]},{\"type\":\"function\",\"name\":\"getAllocation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRevokableAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakeableAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"}]",
}

// LatpABI is the input ABI used to generate the binding from.
// Deprecated: Use LatpMetaData.ABI instead.
var LatpABI = LatpMetaData.ABI

// Latp is an auto generated Go binding around an Ethereum contract.
type Latp struct {
	LatpCaller     // Read-only binding to the contract
	LatpTransactor // Write-only binding to the contract
	LatpFilterer   // Log filterer for contract events
}

// LatpCaller is an auto generated read-only Go binding around an Ethereum contract.
type LatpCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LatpTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LatpTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LatpFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LatpFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LatpSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LatpSession struct {
	Contract     *Latp             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LatpCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LatpCallerSession struct {
	Contract *LatpCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// LatpTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LatpTransactorSession struct {
	Contract     *LatpTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LatpRaw is an auto generated low-level Go binding around an Ethereum contract.
type LatpRaw struct {
	Contract *Latp // Generic contract binding to access the raw methods on
}

// LatpCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LatpCallerRaw struct {
	Contract *LatpCaller // Generic read-only contract binding to access the raw methods on
}

// LatpTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LatpTransactorRaw struct {
	Contract *LatpTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLatp creates a new instance of Latp, bound to a specific deployed contract.
func NewLatp(address common.Address, backend bind.ContractBackend) (*Latp, error) {
	contract, err := bindLatp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Latp{LatpCaller: LatpCaller{contract: contract}, LatpTransactor: LatpTransactor{contract: contract}, LatpFilterer: LatpFilterer{contract: contract}}, nil
}

// NewLatpCaller creates a new read-only instance of Latp, bound to a specific deployed contract.
func NewLatpCaller(address common.Address, caller bind.ContractCaller) (*LatpCaller, error) {
	contract, err := bindLatp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LatpCaller{contract: contract}, nil
}

// NewLatpTransactor creates a new write-only instance of Latp, bound to a specific deployed contract.
func NewLatpTransactor(address common.Address, transactor bind.ContractTransactor) (*LatpTransactor, error) {
	contract, err := bindLatp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LatpTransactor{contract: contract}, nil
}

// NewLatpFilterer creates a new log filterer instance of Latp, bound to a specific deployed contract.
func NewLatpFilterer(address common.Address, filterer bind.ContractFilterer) (*LatpFilterer, error) {
	contract, err := bindLatp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LatpFilterer{contract: contract}, nil
}

// bindLatp binds a generic wrapper to an already deployed contract.
func bindLatp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LatpMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Latp *LatpRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Latp.Contract.LatpCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Latp *LatpRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Latp.Contract.LatpTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Latp *LatpRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Latp.Contract.LatpTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Latp *LatpCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Latp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Latp *LatpTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Latp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Latp *LatpTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Latp.Contract.contract.Transact(opts, method, params...)
}

// GetAllocation is a free data retrieval call binding the contract method 0x0c9e1e8e.
//
// Solidity: function getAllocation() view returns(uint256)
func (_Latp *LatpCaller) GetAllocation(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Latp.contract.Call(opts, &out, "getAllocation")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAllocation is a free data retrieval call binding the contract method 0x0c9e1e8e.
//
// Solidity: function getAllocation() view returns(uint256)
func (_Latp *LatpSession) GetAllocation() (*big.Int, error) {
	return _Latp.Contract.GetAllocation(&_Latp.CallOpts)
}

// GetAllocation is a free data retrieval call binding the contract method 0x0c9e1e8e.
//
// Solidity: function getAllocation() view returns(uint256)
func (_Latp *LatpCallerSession) GetAllocation() (*big.Int, error) {
	return _Latp.Contract.GetAllocation(&_Latp.CallOpts)
}

// GetClaimable is a free data retrieval call binding the contract method 0xee28b744.
//
// Solidity: function getClaimable() view returns(uint256)
func (_Latp *LatpCaller) GetClaimable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Latp.contract.Call(opts, &out, "getClaimable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetClaimable is a free data retrieval call binding the contract method 0xee28b744.
//
// Solidity: function getClaimable() view returns(uint256)
func (_Latp *LatpSession) GetClaimable() (*big.Int, error) {
	return _Latp.Contract.GetClaimable(&_Latp.CallOpts)
}

// GetClaimable is a free data retrieval call binding the contract method 0xee28b744.
//
// Solidity: function getClaimable() view returns(uint256)
func (_Latp *LatpCallerSession) GetClaimable() (*big.Int, error) {
	return _Latp.Contract.GetClaimable(&_Latp.CallOpts)
}

// GetClaimed is a free data retrieval call binding the contract method 0xae3bb460.
//
// Solidity: function getClaimed() view returns(uint256)
func (_Latp *LatpCaller) GetClaimed(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Latp.contract.Call(opts, &out, "getClaimed")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetClaimed is a free data retrieval call binding the contract method 0xae3bb460.
//
// Solidity: function getClaimed() view returns(uint256)
func (_Latp *LatpSession) GetClaimed() (*big.Int, error) {
	return _Latp.Contract.GetClaimed(&_Latp.CallOpts)
}

// GetClaimed is a free data retrieval call binding the contract method 0xae3bb460.
//
// Solidity: function getClaimed() view returns(uint256)
func (_Latp *LatpCallerSession) GetClaimed() (*big.Int, error) {
	return _Latp.Contract.GetClaimed(&_Latp.CallOpts)
}

// GetRevokableAmount is a free data retrieval call binding the contract method 0xce828b06.
//
// Solidity: function getRevokableAmount() view returns(uint256)
func (_Latp *LatpCaller) GetRevokableAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Latp.contract.Call(opts, &out, "getRevokableAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRevokableAmount is a free data retrieval call binding the contract method 0xce828b06.
//
// Solidity: function getRevokableAmount() view returns(uint256)
func (_Latp *LatpSession) GetRevokableAmount() (*big.Int, error) {
	return _Latp.Contract.GetRevokableAmount(&_Latp.CallOpts)
}

// GetRevokableAmount is a free data retrieval call binding the contract method 0xce828b06.
//
// Solidity: function getRevokableAmount() view returns(uint256)
func (_Latp *LatpCallerSession) GetRevokableAmount() (*big.Int, error) {
	return _Latp.Contract.GetRevokableAmount(&_Latp.CallOpts)
}

// GetStakeableAmount is a free data retrieval call binding the contract method 0x592b07cd.
//
// Solidity: function getStakeableAmount() view returns(uint256)
func (_Latp *LatpCaller) GetStakeableAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Latp.contract.Call(opts, &out, "getStakeableAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakeableAmount is a free data retrieval call binding the contract method 0x592b07cd.
//
// Solidity: function getStakeableAmount() view returns(uint256)
func (_Latp *LatpSession) GetStakeableAmount() (*big.Int, error) {
	return _Latp.Contract.GetStakeableAmount(&_Latp.CallOpts)
}

// GetStakeableAmount is a free data retrieval call binding the contract method 0x592b07cd.
//
// Solidity: function getStakeableAmount() view returns(uint256)
func (_Latp *LatpCallerSession) GetStakeableAmount() (*big.Int, error) {
	return _Latp.Contract.GetStakeableAmount(&_Latp.CallOpts)
}
