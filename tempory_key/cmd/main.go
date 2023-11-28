package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/bubbleDevelop/tempory_key/contract"
	"github.com/bubbleDevelop/tempory_key/tempPk"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	workPrivateKey                *ecdsa.PrivateKey
	workAddress                   common.Address // 0x4307ffd08477668dC6d9f49f90b084B1f1CCC82b
	tempPrivateKey                *ecdsa.PrivateKey
	tempAddress                   common.Address
	operatorPrivateKey            *ecdsa.PrivateKey
	operatorAddress               common.Address // 0x9FD5bD701Fc8105E46399104AC4B8c1B391df760
	tempPrivateKeyContractAddress common.Address // 0x1000000000000000000000000000000000000021
	gameContractAddress           common.Address
)

func init() {
	var err error

	// work address
	workPrivateKey, err = crypto.HexToECDSA("47d790a96ca73b23fbb65a6b911b8b57a1d915d364f12e2bc7fae83c196c9c97")
	if nil != err {
		panic(err)
	}
	publicKey := workPrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	workAddress = crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("workaddress: ", workAddress.Hex())

	// temporary address
	tempPrivateKey, err = crypto.HexToECDSA("e8e14120bb5c085622253540e886527d24746cd42d764a5974be47090d3cbc42")
	if err != nil {
		panic(err)
	}
	tempPublicKey := tempPrivateKey.Public()
	tempPublicKeyECDSA, ok := tempPublicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	tempAddress = crypto.PubkeyToAddress(*tempPublicKeyECDSA)
	fmt.Println("tempAddress: ", tempAddress.Hex())

	// operator address
	operatorPrivateKey, err = crypto.HexToECDSA("e3166f9f62f109d19fb1b73f8d9c647530153cd03822d3091951081bac7f7c5e")
	if err != nil {
		panic(err)
	}
	operatorPublicKey := operatorPrivateKey.Public()
	operatorPublicKeyECDSA, ok := operatorPublicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	operatorAddress = crypto.PubkeyToAddress(*operatorPublicKeyECDSA)
	fmt.Println("operatorAddress: ", operatorAddress.Hex())

	// temporary contract address
	tempPrivateKeyContractAddress = common.HexToAddress("0x1000000000000000000000000000000000000021")

	// game contract address
	gameContractAddress = common.HexToAddress("0x3a9d4C411F8A37be2f34B208A03719a2cCf4Aee0")
}

func getChainInfo(client *ethclient.Client, fromAddress common.Address) (nonce uint64, chainId, gasPrice *big.Int, err error) {
	nonce, err = client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return
	}
	chainId, err = client.ChainID(context.Background())
	if err != nil {
		return
	}

	gasPrice, err = client.SuggestGasPrice(context.Background())
	return
}

func sendTempPrivateKeyContractTx(client *ethclient.Client, privateKey *ecdsa.PrivateKey, fromAddress common.Address, input []byte) error {
	nonce, chainId, gasPrice, err := getChainInfo(client, fromAddress)
	if err != nil {
		return err
	}
	value := big.NewInt(0)
	gasLimit := uint64(3000000)
	rawTx := types.NewTransaction(nonce, tempPrivateKeyContractAddress, value, gasLimit, gasPrice, input)

	// sign transaction
	signer := types.NewEIP155Signer(chainId)
	sigTransaction, err := types.SignTx(rawTx, signer, privateKey)
	if err != nil {
		return err
	}

	// send transaction
	err = client.SendTransaction(context.Background(), sigTransaction)
	if err != nil {

		return err
	}
	fmt.Println("send transaction success,tx: ", sigTransaction.Hash().Hex())
	return nil
}

func bindTempPrivateKeyCall(client *ethclient.Client, gameContractAddress, tempAddress common.Address, period []byte) error {
	input := tempPk.BindTempPrivateKey(gameContractAddress, tempAddress, period)
	return sendTempPrivateKeyContractTx(client, workPrivateKey, workAddress, input)
}

func invalidateTempPrivateKeyCall(client *ethclient.Client, gameContractAddress, tempAddress common.Address) error {
	input := tempPk.InvalidateTempPrivateKey(gameContractAddress, tempAddress)
	return sendTempPrivateKeyContractTx(client, workPrivateKey, workAddress, input)
}

func behalfSignatureCall(client *ethclient.Client, workAddress, gameContractAddress common.Address, periodArg, input []byte) error {
	paras := tempPk.BehalfSignature(workAddress, gameContractAddress, periodArg, input)
	return sendTempPrivateKeyContractTx(client, tempPrivateKey, tempAddress, paras)
}

func addLineOfCreditCall(client *ethclient.Client, gameContractAddress, workAddress common.Address, addValue *big.Int) error {
	input := tempPk.AddLineOfCredit(gameContractAddress, workAddress, addValue)
	return sendTempPrivateKeyContractTx(client, operatorPrivateKey, operatorAddress, input)
}

func gameInfo(client *ethclient.Client) {
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		fmt.Println("get chain id error!!!")
		panic(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("get suggest gas price error!!!")
		panic(err)
	}

	// 创建合约对象
	gameContract, err := contract.NewGame(gameContractAddress, client)
	if err != nil {
		fmt.Println("new game contract instance error!!!")
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(workPrivateKey, chainId)
	if err != nil {
		fmt.Println("new keyed transaction with chain id error!!!")
		panic(err)
	}

	// 设置 issuer 地址
	tx, err := gameContract.SetIssuer(&bind.TransactOpts{
		From: auth.From,
		//Nonce:     nil,
		Signer: auth.Signer,
		//Value:     nil,
		GasPrice: gasPrice,
		//GasFeeCap: nil,
		//GasTipCap: nil,
		GasLimit: uint64(3000000),
		//Context:   nil,
		//NoSend:    false,
	}, workAddress)
	if err != nil {
		fmt.Println("set issuer error!!!")
		panic(err)
	}
	fmt.Println("set issuer transaction: ", tx.Hash())

	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		fmt.Println("wait set issuer transaction error!!!")
		panic(err)
	}
	setIssuerReceipt, err := json.Marshal(receipt)
	if err != nil {
		fmt.Println("get set issuer receipt error!!!")
		panic(err)
	}
	fmt.Println("set issuer transaction receipt: ", string(setIssuerReceipt))

	// 查询 issuer 地址
	issuerRes, err := gameContract.Issuer(&bind.CallOpts{
		Pending:     false,
		From:        common.Address{},
		BlockNumber: nil,
		Context:     nil,
	})
	if err != nil {
		fmt.Println("get issuer error!!!")
		panic(err)
	}
	fmt.Println("issuer", issuerRes)

	// 设置授信额度
	tx, err = gameContract.SetLineOfCredit(&bind.TransactOpts{
		From: auth.From,
		//Nonce:     nil,
		Signer: auth.Signer,
		//Value:     nil,
		GasPrice: gasPrice,
		//GasFeeCap: nil,
		//GasTipCap: nil,
		GasLimit: uint64(3000000),
		//Context:   nil,
		//NoSend:    false,
	}, big.NewInt(1234567890))
	if err != nil {
		fmt.Println("set line of credit error!!!")
		panic(err)
	}
	fmt.Println("set line of credit transaction: ", tx.Hash())

	receipt, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		fmt.Println("wait set line of credit transaction error!!!")
		panic(err)
	}
	setLineOfCreditReceipt, err := json.Marshal(receipt)
	if err != nil {
		fmt.Println("get set line of credit receipt error!!!")
		panic(err)
	}
	fmt.Println("set line of credit transaction receipt: ", string(setLineOfCreditReceipt))

	// 查询授信额度
	lineOfCreditRes, err := gameContract.LineOfCredit(&bind.CallOpts{
		Pending:     false,
		From:        common.Address{},
		BlockNumber: nil,
		Context:     nil,
	})
	if err != nil {
		fmt.Println("get line of credit error!!!")
		panic(err)
	}
	fmt.Println("line of credit: ", lineOfCreditRes)

	// 移动位置
	tx, err = gameContract.MovePlayer(&bind.TransactOpts{
		From: auth.From,
		//Nonce:     nil,
		Signer: auth.Signer,
		//Value:     nil,
		GasPrice: gasPrice,
		//GasFeeCap: nil,
		//GasTipCap: nil,
		GasLimit: uint64(3000000),
		//Context:   nil,
		//NoSend:    false,
	}, big.NewInt(1234567890))
	if err != nil {
		fmt.Println("move player error!!!")
		panic(err)
	}
	fmt.Println("move player transaction: ", tx.Hash())

	receipt, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		fmt.Println("wait move player transaction error!!!")
		panic(err)
	}
	movePlayerReceipt, err := json.Marshal(receipt)
	if err != nil {
		fmt.Println("move player receipt error!!!")
		panic(err)
	}
	fmt.Println("move player transaction receipt: ", string(movePlayerReceipt))

	// 查询位置信息
	positionRes, err := gameContract.Position(&bind.CallOpts{
		Pending:     false,
		From:        common.Address{},
		BlockNumber: nil,
		Context:     nil,
	})
	if err != nil {
		fmt.Println("get position error!!!")
		panic(err)
	}
	fmt.Println("position: ", positionRes)
}

func main() {
	fmt.Println("main function")
	// 链接服务器
	conn, err := ethclient.Dial("http://192.168.31.115:18001")
	if err != nil {
		fmt.Println("Dial err", err)
		return
	}
	defer conn.Close()

	// game.sol
	gameInfo(conn)

	// 	// 绑定临时私钥
	// 	err = bindTempPrivateKeyCall(conn, gameContractAddress, tempAddress, []byte("Hello World"))
	// 	if nil != err {
	// 		fmt.Println(err)
	// 	}

	// 	// 合约调用代签
	// 	input := tempPk.MovePlayer(big.NewInt(12345))
	// 	err = behalfSignatureCall(conn, workAddress, gameContractAddress, []byte("Hello World"), input)
	// 	if nil != err {
	// 		fmt.Println(err)
	// 	}

	// 	// 增加授信额度
	// 	err = addLineOfCreditCall(conn, gameContractAddress, workAddress, big.NewInt(1234567890))
	// 	if nil != err {
	// 		fmt.Println(err)
	// 	}

	// 	// 作废临时私钥
	// 	err = invalidateTempPrivateKeyCall(conn, gameContractAddress, tempAddress)
	// 	if nil != err {
	// 		fmt.Println(err)
	// 	}
}
