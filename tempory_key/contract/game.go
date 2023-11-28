// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// GameMetaData contains all meta data concerning the Game contract.
var GameMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"issuer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lineOfCredit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"postionMove\",\"type\":\"uint256\"}],\"name\":\"movePlayer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"position\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"issuerAddress\",\"type\":\"address\"}],\"name\":\"setIssuer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lineOfCreditAmount\",\"type\":\"uint256\"}],\"name\":\"setLineOfCredit\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// GameABI is the input ABI used to generate the binding from.
// Deprecated: Use GameMetaData.ABI instead.
var GameABI = GameMetaData.ABI

// Game is an auto generated Go binding around an Ethereum contract.
type Game struct {
	GameCaller     // Read-only binding to the contract
	GameTransactor // Write-only binding to the contract
	GameFilterer   // Log filterer for contract events
}

// GameCaller is an auto generated read-only Go binding around an Ethereum contract.
type GameCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GameTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GameTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GameFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GameFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GameSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GameSession struct {
	Contract     *Game             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GameCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GameCallerSession struct {
	Contract *GameCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// GameTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GameTransactorSession struct {
	Contract     *GameTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GameRaw is an auto generated low-level Go binding around an Ethereum contract.
type GameRaw struct {
	Contract *Game // Generic contract binding to access the raw methods on
}

// GameCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GameCallerRaw struct {
	Contract *GameCaller // Generic read-only contract binding to access the raw methods on
}

// GameTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GameTransactorRaw struct {
	Contract *GameTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGame creates a new instance of Game, bound to a specific deployed contract.
func NewGame(address common.Address, backend bind.ContractBackend) (*Game, error) {
	contract, err := bindGame(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Game{GameCaller: GameCaller{contract: contract}, GameTransactor: GameTransactor{contract: contract}, GameFilterer: GameFilterer{contract: contract}}, nil
}

// NewGameCaller creates a new read-only instance of Game, bound to a specific deployed contract.
func NewGameCaller(address common.Address, caller bind.ContractCaller) (*GameCaller, error) {
	contract, err := bindGame(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GameCaller{contract: contract}, nil
}

// NewGameTransactor creates a new write-only instance of Game, bound to a specific deployed contract.
func NewGameTransactor(address common.Address, transactor bind.ContractTransactor) (*GameTransactor, error) {
	contract, err := bindGame(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GameTransactor{contract: contract}, nil
}

// NewGameFilterer creates a new log filterer instance of Game, bound to a specific deployed contract.
func NewGameFilterer(address common.Address, filterer bind.ContractFilterer) (*GameFilterer, error) {
	contract, err := bindGame(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GameFilterer{contract: contract}, nil
}

// bindGame binds a generic wrapper to an already deployed contract.
func bindGame(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GameMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Game *GameRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Game.Contract.GameCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Game *GameRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Game.Contract.GameTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Game *GameRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Game.Contract.GameTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Game *GameCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Game.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Game *GameTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Game.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Game *GameTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Game.Contract.contract.Transact(opts, method, params...)
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() view returns(address)
func (_Game *GameCaller) Issuer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Game.contract.Call(opts, &out, "issuer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() view returns(address)
func (_Game *GameSession) Issuer() (common.Address, error) {
	return _Game.Contract.Issuer(&_Game.CallOpts)
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() view returns(address)
func (_Game *GameCallerSession) Issuer() (common.Address, error) {
	return _Game.Contract.Issuer(&_Game.CallOpts)
}

// LineOfCredit is a free data retrieval call binding the contract method 0x99a3c623.
//
// Solidity: function lineOfCredit() view returns(uint256)
func (_Game *GameCaller) LineOfCredit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Game.contract.Call(opts, &out, "lineOfCredit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LineOfCredit is a free data retrieval call binding the contract method 0x99a3c623.
//
// Solidity: function lineOfCredit() view returns(uint256)
func (_Game *GameSession) LineOfCredit() (*big.Int, error) {
	return _Game.Contract.LineOfCredit(&_Game.CallOpts)
}

// LineOfCredit is a free data retrieval call binding the contract method 0x99a3c623.
//
// Solidity: function lineOfCredit() view returns(uint256)
func (_Game *GameCallerSession) LineOfCredit() (*big.Int, error) {
	return _Game.Contract.LineOfCredit(&_Game.CallOpts)
}

// Position is a free data retrieval call binding the contract method 0x09218e91.
//
// Solidity: function position() view returns(uint256)
func (_Game *GameCaller) Position(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Game.contract.Call(opts, &out, "position")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Position is a free data retrieval call binding the contract method 0x09218e91.
//
// Solidity: function position() view returns(uint256)
func (_Game *GameSession) Position() (*big.Int, error) {
	return _Game.Contract.Position(&_Game.CallOpts)
}

// Position is a free data retrieval call binding the contract method 0x09218e91.
//
// Solidity: function position() view returns(uint256)
func (_Game *GameCallerSession) Position() (*big.Int, error) {
	return _Game.Contract.Position(&_Game.CallOpts)
}

// MovePlayer is a paid mutator transaction binding the contract method 0x6a3fc9d2.
//
// Solidity: function movePlayer(uint256 postionMove) returns(bool success)
func (_Game *GameTransactor) MovePlayer(opts *bind.TransactOpts, postionMove *big.Int) (*types.Transaction, error) {
	return _Game.contract.Transact(opts, "movePlayer", postionMove)
}

// MovePlayer is a paid mutator transaction binding the contract method 0x6a3fc9d2.
//
// Solidity: function movePlayer(uint256 postionMove) returns(bool success)
func (_Game *GameSession) MovePlayer(postionMove *big.Int) (*types.Transaction, error) {
	return _Game.Contract.MovePlayer(&_Game.TransactOpts, postionMove)
}

// MovePlayer is a paid mutator transaction binding the contract method 0x6a3fc9d2.
//
// Solidity: function movePlayer(uint256 postionMove) returns(bool success)
func (_Game *GameTransactorSession) MovePlayer(postionMove *big.Int) (*types.Transaction, error) {
	return _Game.Contract.MovePlayer(&_Game.TransactOpts, postionMove)
}

// SetIssuer is a paid mutator transaction binding the contract method 0x55cc4e57.
//
// Solidity: function setIssuer(address issuerAddress) returns(bool success)
func (_Game *GameTransactor) SetIssuer(opts *bind.TransactOpts, issuerAddress common.Address) (*types.Transaction, error) {
	return _Game.contract.Transact(opts, "setIssuer", issuerAddress)
}

// SetIssuer is a paid mutator transaction binding the contract method 0x55cc4e57.
//
// Solidity: function setIssuer(address issuerAddress) returns(bool success)
func (_Game *GameSession) SetIssuer(issuerAddress common.Address) (*types.Transaction, error) {
	return _Game.Contract.SetIssuer(&_Game.TransactOpts, issuerAddress)
}

// SetIssuer is a paid mutator transaction binding the contract method 0x55cc4e57.
//
// Solidity: function setIssuer(address issuerAddress) returns(bool success)
func (_Game *GameTransactorSession) SetIssuer(issuerAddress common.Address) (*types.Transaction, error) {
	return _Game.Contract.SetIssuer(&_Game.TransactOpts, issuerAddress)
}

// SetLineOfCredit is a paid mutator transaction binding the contract method 0x8cc64805.
//
// Solidity: function setLineOfCredit(uint256 lineOfCreditAmount) returns(bool success)
func (_Game *GameTransactor) SetLineOfCredit(opts *bind.TransactOpts, lineOfCreditAmount *big.Int) (*types.Transaction, error) {
	return _Game.contract.Transact(opts, "setLineOfCredit", lineOfCreditAmount)
}

// SetLineOfCredit is a paid mutator transaction binding the contract method 0x8cc64805.
//
// Solidity: function setLineOfCredit(uint256 lineOfCreditAmount) returns(bool success)
func (_Game *GameSession) SetLineOfCredit(lineOfCreditAmount *big.Int) (*types.Transaction, error) {
	return _Game.Contract.SetLineOfCredit(&_Game.TransactOpts, lineOfCreditAmount)
}

// SetLineOfCredit is a paid mutator transaction binding the contract method 0x8cc64805.
//
// Solidity: function setLineOfCredit(uint256 lineOfCreditAmount) returns(bool success)
func (_Game *GameTransactorSession) SetLineOfCredit(lineOfCreditAmount *big.Int) (*types.Transaction, error) {
	return _Game.Contract.SetLineOfCredit(&_Game.TransactOpts, lineOfCreditAmount)
}
