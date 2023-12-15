package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/bubbleDevelop/tempory_key/contract"
	"github.com/bubbleDevelop/tempory_key/tempPk"
	"github.com/ethereum/go-ethereum"
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
	tempAddress                   common.Address // 0xA2088F51Ea1f9BA308F5014150961e5a6E0A4E13
	operatorPrivateKey            *ecdsa.PrivateKey
	operatorAddress               common.Address // 0x9FD5bD701Fc8105E46399104AC4B8c1B391df760
	tempPrivateKeyContractAddress common.Address // 0x1000000000000000000000000000000000000021
	gameContractAddress           common.Address
)

func init() {
	var err error

	// work address 47d790a96ca73b23fbb65a6b911b8b57a1d915d364f12e2bc7fae83c196c9c97
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

func tempPrivateKeyContractEstimateGas(client *ethclient.Client, input []byte, workAddress common.Address) (uint64, error) {

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return uint64(0), err
	}
	value := big.NewInt(0)
	gasLimit := uint64(3000000)
	msg := ethereum.CallMsg{Data: input, Gas: gasLimit, GasPrice: gasPrice, To: &tempPrivateKeyContractAddress, From: workAddress, Value: value}
	return client.EstimateGas(context.Background(), msg)
}

func simpleTransfer(client *ethclient.Client, privateKey *ecdsa.PrivateKey, to common.Address, amount *big.Int) error {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	from := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		fmt.Println("cannot get pending nonce")
		return err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("cannot get suggest gas price")
		return err
	}
	gasLimit := uint64(21000)
	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		fmt.Println("cannot get chain id")
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		fmt.Println("cannot get chain id")
		return err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		fmt.Println("cannot send transaction")
		return err
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())

	receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		fmt.Println("wait transaction error!!!")
		return err
	}
	setLineOfCreditReceipt, err := json.Marshal(receipt)
	if err != nil {
		fmt.Println("marshal error!!!")
		return err
	}
	fmt.Println("transaction receipt: ", string(setLineOfCreditReceipt))
	return nil
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

	receipt, err := bind.WaitMined(context.Background(), client, sigTransaction)
	if err != nil {
		fmt.Println("wait transaction error!!!")
		return err
	}
	setLineOfCreditReceipt, err := json.Marshal(receipt)
	if err != nil {
		fmt.Println("marshal error!!!")
		return err
	}
	fmt.Println("transaction receipt: ", string(setLineOfCreditReceipt))

	return nil
}

func bindTempPrivateKeyCall(client *ethclient.Client, gameContractAddress, tempAddress common.Address, period []byte) error {
	input := tempPk.BindTempPrivateKey(gameContractAddress, tempAddress, period)
	return sendTempPrivateKeyContractTx(client, workPrivateKey, workAddress, input)
}

func bindTempPrivateKeyEstimate(client *ethclient.Client, gameContractAddress, tempAddress common.Address, period []byte) (uint64, error) {
	input := tempPk.BindTempPrivateKey(gameContractAddress, tempAddress, period)
	return tempPrivateKeyContractEstimateGas(client, input, workAddress)
}

func invalidateTempPrivateKeyCall(client *ethclient.Client, gameContractAddress, tempAddress common.Address) error {
	input := tempPk.InvalidateTempPrivateKey(gameContractAddress, tempAddress)
	return sendTempPrivateKeyContractTx(client, workPrivateKey, workAddress, input)
}

func invalidateTempPrivateKeyEstimate(client *ethclient.Client, gameContractAddress, tempAddress common.Address) (uint64, error) {
	input := tempPk.InvalidateTempPrivateKey(gameContractAddress, operatorAddress)
	return tempPrivateKeyContractEstimateGas(client, input, workAddress)
}

func behalfSignatureCall(client *ethclient.Client, workAddress, gameContractAddress common.Address, periodArg, input []byte) error {
	paras := tempPk.BehalfSignature(workAddress, gameContractAddress, periodArg, input)
	return sendTempPrivateKeyContractTx(client, tempPrivateKey, tempAddress, paras)
}

func behalfSignatureEstimate(client *ethclient.Client, workAddress, gameContractAddress common.Address, periodArg, input []byte) (uint64, error) {
	paras := tempPk.BehalfSignature(workAddress, gameContractAddress, periodArg, input)
	return tempPrivateKeyContractEstimateGas(client, paras, tempAddress)
}

func addLineOfCreditCall(client *ethclient.Client, gameContractAddress, workAddress common.Address, addValue *big.Int) error {
	input := tempPk.AddLineOfCredit(gameContractAddress, workAddress, addValue)
	return sendTempPrivateKeyContractTx(client, operatorPrivateKey, operatorAddress, input)
}

func addLineOfCreditEstimate(client *ethclient.Client, gameContractAddress, workAddress common.Address, addValue *big.Int) (uint64, error) {
	input := tempPk.AddLineOfCredit(gameContractAddress, workAddress, addValue)
	return tempPrivateKeyContractEstimateGas(client, input, operatorAddress)
}

func getLineOfCreditCall(client *ethclient.Client, gameContractAddress, workAddress common.Address) (string, error) {
	input := tempPk.GetLineOfCredit(gameContractAddress)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	value := big.NewInt(0)
	gasLimit := uint64(3000000)
	msg := ethereum.CallMsg{Data: input, Gas: gasLimit, GasPrice: gasPrice, To: &tempPrivateKeyContractAddress, From: workAddress, Value: value}
	resBytes, err := client.PendingCallContract(context.Background(), msg)
	if err != nil {
		return "", err
	}

	type Result struct {
		Code uint32
		Ret  interface{}
	}

	var res Result
	fmt.Println(string(resBytes))
	json.Unmarshal(resBytes, &res)
	fmt.Println(res)

	str, ok := res.Ret.(string)
	if !ok {
		return "", fmt.Errorf("invalid type")
	}

	return str, nil
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
	}, operatorAddress)
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
	}, big.NewInt(3000000000000000000))
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

func simpleTransferTest() {

	fmt.Println("main function")
	// 链接服务器
	conn, err := ethclient.Dial("http://192.168.31.115:18001")
	if err != nil {
		fmt.Println("Dial err", err)
		return
	}
	defer conn.Close()

	err = simpleTransfer(conn, tempPrivateKey, workAddress, big.NewInt(100000000000000))
	if nil != err {
		fmt.Println(err)
	}

	err = simpleTransfer(conn, workPrivateKey, tempAddress, big.NewInt(100000000000000))
	if nil != err {
		fmt.Println(err)
	}

	err = simpleTransfer(conn, tempPrivateKey, workAddress, big.NewInt(100000000000000))
	if nil != err {
		fmt.Println(err)
	}
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

	// 绑定临时私钥
	fmt.Println("bindTempPrivateKeyCall begin")
	err = bindTempPrivateKeyCall(conn, gameContractAddress, tempAddress, []byte("Hello World"))
	if nil != err {
		fmt.Println(err)
	}
	fmt.Println("bindTempPrivateKeyCall end")

	// gasUsed, err := invalidateTempPrivateKeyEstimate(conn, gameContractAddress, tempAddress)
	// if nil != err {
	// 	fmt.Println("invalidateTempPrivateKeyEstimate error: ", err)
	// 	return
	// }
	// fmt.Println("invalidateTempPrivateKeyEstimate gas used: ", gasUsed)

	// // gasUsed, err := bindTempPrivateKeyEstimate(conn, gameContractAddress, tempAddress, []byte("Hello World"))
	// // if nil != err {
	// // 	fmt.Println("bindTempPrivateKeyEstimate error: ", err)
	// // 	return
	// // }
	// // fmt.Println("bindTempPrivateKeyEstimate gas used: ", gasUsed)

	input := tempPk.MovePlayer(big.NewInt(12345))
	gasUsed, err := behalfSignatureEstimate(conn, workAddress, gameContractAddress, []byte("Hello World"), input)
	if nil != err {
		fmt.Println("behalfSignatureEstimate error: ", err)
		return
	}
	fmt.Println("behalfSignatureEstimate gas used: ", gasUsed)

	// 合约调用代签
	lineOfCredit, err := getLineOfCreditCall(conn, gameContractAddress, workAddress)
	if nil != err {
		fmt.Println(err)
	}
	fmt.Println("lineOfCredit: ", lineOfCredit)
	operatorValue, err := conn.PendingBalanceAt(context.Background(), operatorAddress)
	if nil != err {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("operator balance: ", operatorValue.String())
	fmt.Println("behalfSignatureCall begin")
	input = tempPk.MovePlayer(big.NewInt(12345))
	err = behalfSignatureCall(conn, workAddress, gameContractAddress, []byte("Hello World"), input)
	if nil != err {
		fmt.Println(err)
	}
	fmt.Println("behalfSignatureCall end")
	operatorValue, err = conn.PendingBalanceAt(context.Background(), operatorAddress)
	if nil != err {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("operator balance: ", operatorValue.String())
	lineOfCredit, err = getLineOfCreditCall(conn, gameContractAddress, workAddress)
	if nil != err {
		fmt.Println(err)
	}
	fmt.Println("lineOfCredit: ", lineOfCredit)

	// 增加授信额度
	// fmt.Println("addLineOfCreditCall begin")
	// err = addLineOfCreditCall(conn, gameContractAddress, workAddress, big.NewInt(1234567890123456789))
	// if nil != err {
	// 	fmt.Println(err)
	// }
	// fmt.Println("addLineOfCreditCall end")

	// 查询授信额度
	// fmt.Println("getLineOfCreditCall begin")
	// lineOfCredit, err = getLineOfCreditCall(conn, gameContractAddress, workAddress)
	// if nil != err {
	// 	fmt.Println(err)
	// }
	// fmt.Println("line of credit:", lineOfCredit)
	// fmt.Println("getLineOfCreditCall end")

	// fmt.Println("behalfSignatureCall begin")
	// input = tempPk.MovePlayer(big.NewInt(12345))
	// err = behalfSignatureCall(conn, workAddress, gameContractAddress, []byte("Hello World"), input)
	// if nil != err {
	// 	fmt.Println(err)
	// }
	// fmt.Println("behalfSignatureCall end")

	// // 作废临时私钥
	// err = invalidateTempPrivateKeyCall(conn, gameContractAddress, tempAddress)
	// if nil != err {
	// 	fmt.Println(err)
	// }
}
