// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package matp

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

// MatpMetaData contains all meta data concerning the Matp contract.
var MatpMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getClaimed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"getClaimable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}]},{\"type\":\"function\",\"name\":\"getAllocation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRevokableAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakeableAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"getBeneficiary\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}]},{\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"getIsRevoked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}]}]",
}

// MatpABI is the input ABI used to generate the binding from.
// Deprecated: Use MatpMetaData.ABI instead.
var MatpABI = MatpMetaData.ABI

// Matp is an auto generated Go binding around an Ethereum contract.
type Matp struct {
	MatpCaller     // Read-only binding to the contract
	MatpTransactor // Write-only binding to the contract
	MatpFilterer   // Log filterer for contract events
}

// MatpCaller is an auto generated read-only Go binding around an Ethereum contract.
type MatpCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MatpTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MatpTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MatpFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MatpFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MatpSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MatpSession struct {
	Contract     *Matp             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MatpCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MatpCallerSession struct {
	Contract *MatpCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// MatpTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MatpTransactorSession struct {
	Contract     *MatpTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MatpRaw is an auto generated low-level Go binding around an Ethereum contract.
type MatpRaw struct {
	Contract *Matp // Generic contract binding to access the raw methods on
}

// MatpCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MatpCallerRaw struct {
	Contract *MatpCaller // Generic read-only contract binding to access the raw methods on
}

// MatpTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MatpTransactorRaw struct {
	Contract *MatpTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMatp creates a new instance of Matp, bound to a specific deployed contract.
func NewMatp(address common.Address, backend bind.ContractBackend) (*Matp, error) {
	contract, err := bindMatp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Matp{MatpCaller: MatpCaller{contract: contract}, MatpTransactor: MatpTransactor{contract: contract}, MatpFilterer: MatpFilterer{contract: contract}}, nil
}

// NewMatpCaller creates a new read-only instance of Matp, bound to a specific deployed contract.
func NewMatpCaller(address common.Address, caller bind.ContractCaller) (*MatpCaller, error) {
	contract, err := bindMatp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MatpCaller{contract: contract}, nil
}

// NewMatpTransactor creates a new write-only instance of Matp, bound to a specific deployed contract.
func NewMatpTransactor(address common.Address, transactor bind.ContractTransactor) (*MatpTransactor, error) {
	contract, err := bindMatp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MatpTransactor{contract: contract}, nil
}

// NewMatpFilterer creates a new log filterer instance of Matp, bound to a specific deployed contract.
func NewMatpFilterer(address common.Address, filterer bind.ContractFilterer) (*MatpFilterer, error) {
	contract, err := bindMatp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MatpFilterer{contract: contract}, nil
}

// bindMatp binds a generic wrapper to an already deployed contract.
func bindMatp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MatpMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Matp *MatpRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Matp.Contract.MatpCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Matp *MatpRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Matp.Contract.MatpTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Matp *MatpRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Matp.Contract.MatpTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Matp *MatpCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Matp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Matp *MatpTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Matp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Matp *MatpTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Matp.Contract.contract.Transact(opts, method, params...)
}

// GetAllocation is a free data retrieval call binding the contract method 0x0c9e1e8e.
//
// Solidity: function getAllocation() view returns(uint256)
func (_Matp *MatpCaller) GetAllocation(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Matp.contract.Call(opts, &out, "getAllocation")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAllocation is a free data retrieval call binding the contract method 0x0c9e1e8e.
//
// Solidity: function getAllocation() view returns(uint256)
func (_Matp *MatpSession) GetAllocation() (*big.Int, error) {
	return _Matp.Contract.GetAllocation(&_Matp.CallOpts)
}

// GetAllocation is a free data retrieval call binding the contract method 0x0c9e1e8e.
//
// Solidity: function getAllocation() view returns(uint256)
func (_Matp *MatpCallerSession) GetAllocation() (*big.Int, error) {
	return _Matp.Contract.GetAllocation(&_Matp.CallOpts)
}

// GetBeneficiary is a free data retrieval call binding the contract method 0x565a2e2c.
//
// Solidity: function getBeneficiary() view returns(address)
func (_Matp *MatpCaller) GetBeneficiary(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Matp.contract.Call(opts, &out, "getBeneficiary")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBeneficiary is a free data retrieval call binding the contract method 0x565a2e2c.
//
// Solidity: function getBeneficiary() view returns(address)
func (_Matp *MatpSession) GetBeneficiary() (common.Address, error) {
	return _Matp.Contract.GetBeneficiary(&_Matp.CallOpts)
}

// GetBeneficiary is a free data retrieval call binding the contract method 0x565a2e2c.
//
// Solidity: function getBeneficiary() view returns(address)
func (_Matp *MatpCallerSession) GetBeneficiary() (common.Address, error) {
	return _Matp.Contract.GetBeneficiary(&_Matp.CallOpts)
}

// GetClaimable is a free data retrieval call binding the contract method 0xee28b744.
//
// Solidity: function getClaimable() view returns(uint256)
func (_Matp *MatpCaller) GetClaimable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Matp.contract.Call(opts, &out, "getClaimable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetClaimable is a free data retrieval call binding the contract method 0xee28b744.
//
// Solidity: function getClaimable() view returns(uint256)
func (_Matp *MatpSession) GetClaimable() (*big.Int, error) {
	return _Matp.Contract.GetClaimable(&_Matp.CallOpts)
}

// GetClaimable is a free data retrieval call binding the contract method 0xee28b744.
//
// Solidity: function getClaimable() view returns(uint256)
func (_Matp *MatpCallerSession) GetClaimable() (*big.Int, error) {
	return _Matp.Contract.GetClaimable(&_Matp.CallOpts)
}

// GetClaimed is a free data retrieval call binding the contract method 0xae3bb460.
//
// Solidity: function getClaimed() view returns(uint256)
func (_Matp *MatpCaller) GetClaimed(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Matp.contract.Call(opts, &out, "getClaimed")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetClaimed is a free data retrieval call binding the contract method 0xae3bb460.
//
// Solidity: function getClaimed() view returns(uint256)
func (_Matp *MatpSession) GetClaimed() (*big.Int, error) {
	return _Matp.Contract.GetClaimed(&_Matp.CallOpts)
}

// GetClaimed is a free data retrieval call binding the contract method 0xae3bb460.
//
// Solidity: function getClaimed() view returns(uint256)
func (_Matp *MatpCallerSession) GetClaimed() (*big.Int, error) {
	return _Matp.Contract.GetClaimed(&_Matp.CallOpts)
}

// GetIsRevoked is a free data retrieval call binding the contract method 0x09b058aa.
//
// Solidity: function getIsRevoked() view returns(bool)
func (_Matp *MatpCaller) GetIsRevoked(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Matp.contract.Call(opts, &out, "getIsRevoked")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsRevoked is a free data retrieval call binding the contract method 0x09b058aa.
//
// Solidity: function getIsRevoked() view returns(bool)
func (_Matp *MatpSession) GetIsRevoked() (bool, error) {
	return _Matp.Contract.GetIsRevoked(&_Matp.CallOpts)
}

// GetIsRevoked is a free data retrieval call binding the contract method 0x09b058aa.
//
// Solidity: function getIsRevoked() view returns(bool)
func (_Matp *MatpCallerSession) GetIsRevoked() (bool, error) {
	return _Matp.Contract.GetIsRevoked(&_Matp.CallOpts)
}

// GetRevokableAmount is a free data retrieval call binding the contract method 0xce828b06.
//
// Solidity: function getRevokableAmount() view returns(uint256)
func (_Matp *MatpCaller) GetRevokableAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Matp.contract.Call(opts, &out, "getRevokableAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRevokableAmount is a free data retrieval call binding the contract method 0xce828b06.
//
// Solidity: function getRevokableAmount() view returns(uint256)
func (_Matp *MatpSession) GetRevokableAmount() (*big.Int, error) {
	return _Matp.Contract.GetRevokableAmount(&_Matp.CallOpts)
}

// GetRevokableAmount is a free data retrieval call binding the contract method 0xce828b06.
//
// Solidity: function getRevokableAmount() view returns(uint256)
func (_Matp *MatpCallerSession) GetRevokableAmount() (*big.Int, error) {
	return _Matp.Contract.GetRevokableAmount(&_Matp.CallOpts)
}

// GetStakeableAmount is a free data retrieval call binding the contract method 0x592b07cd.
//
// Solidity: function getStakeableAmount() view returns(uint256)
func (_Matp *MatpCaller) GetStakeableAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Matp.contract.Call(opts, &out, "getStakeableAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakeableAmount is a free data retrieval call binding the contract method 0x592b07cd.
//
// Solidity: function getStakeableAmount() view returns(uint256)
func (_Matp *MatpSession) GetStakeableAmount() (*big.Int, error) {
	return _Matp.Contract.GetStakeableAmount(&_Matp.CallOpts)
}

// GetStakeableAmount is a free data retrieval call binding the contract method 0x592b07cd.
//
// Solidity: function getStakeableAmount() view returns(uint256)
func (_Matp *MatpCallerSession) GetStakeableAmount() (*big.Int, error) {
	return _Matp.Contract.GetStakeableAmount(&_Matp.CallOpts)
}

// GetToken is a free data retrieval call binding the contract method 0x21df0da7.
//
// Solidity: function getToken() view returns(address)
func (_Matp *MatpCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Matp.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetToken is a free data retrieval call binding the contract method 0x21df0da7.
//
// Solidity: function getToken() view returns(address)
func (_Matp *MatpSession) GetToken() (common.Address, error) {
	return _Matp.Contract.GetToken(&_Matp.CallOpts)
}

// GetToken is a free data retrieval call binding the contract method 0x21df0da7.
//
// Solidity: function getToken() view returns(address)
func (_Matp *MatpCallerSession) GetToken() (common.Address, error) {
	return _Matp.Contract.GetToken(&_Matp.CallOpts)
}
