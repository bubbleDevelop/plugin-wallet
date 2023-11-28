package tempPk

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func BindTempPrivateKey(gameContractAddress, tempAddress common.Address, period []byte) []byte {
	params := make([][]byte, 0)
	fnType, _ := rlp.EncodeToBytes(uint16(7200))
	gameContractAddressBytes, _ := rlp.EncodeToBytes(gameContractAddress)
	tempAddressBytes, _ := rlp.EncodeToBytes(tempAddress)
	periodBytes, _ := rlp.EncodeToBytes(period)
	params = append(params, fnType)
	params = append(params, gameContractAddressBytes)
	params = append(params, tempAddressBytes)
	params = append(params, periodBytes)
	buf := new(bytes.Buffer)
	rlp.Encode(buf, params)

	return buf.Bytes()
}

func InvalidateTempPrivateKey(gameContractAddress, tempAddress common.Address) []byte {
	params := make([][]byte, 0)
	fnType, _ := rlp.EncodeToBytes(uint16(7202))
	gameContractAddressBytes, _ := rlp.EncodeToBytes(gameContractAddress)
	tempAddressBytes, _ := rlp.EncodeToBytes(tempAddress)
	params = append(params, fnType)
	params = append(params, gameContractAddressBytes)
	params = append(params, tempAddressBytes)
	buf := new(bytes.Buffer)
	rlp.Encode(buf, params)

	return buf.Bytes()
}

func BehalfSignature(workAddress, gameContractAddress common.Address, periodArg, input []byte) []byte {
	params := make([][]byte, 0)
	fnType, _ := rlp.EncodeToBytes(uint16(7201))
	workAddressBytes, _ := rlp.EncodeToBytes(workAddress)
	gameContractAddressBytes, _ := rlp.EncodeToBytes(gameContractAddress)
	periodArgBytes, _ := rlp.EncodeToBytes(periodArg)
	inputBytes, _ := rlp.EncodeToBytes(input)
	params = append(params, fnType)
	params = append(params, workAddressBytes)
	params = append(params, gameContractAddressBytes)
	params = append(params, periodArgBytes)
	params = append(params, inputBytes)
	buf := new(bytes.Buffer)
	rlp.Encode(buf, params)

	return buf.Bytes()
}

func AddLineOfCredit(gameContractAddress, workAddress common.Address, addValue *big.Int) []byte {
	params := make([][]byte, 0)
	fnType, _ := rlp.EncodeToBytes(uint16(7203))
	gameContractAddressBytes, _ := rlp.EncodeToBytes(gameContractAddress)
	workAddressBytes, _ := rlp.EncodeToBytes(workAddress)
	addValueBytes, _ := rlp.EncodeToBytes(addValue)

	params = append(params, fnType)
	params = append(params, gameContractAddressBytes)
	params = append(params, workAddressBytes)
	params = append(params, addValueBytes)

	buf := new(bytes.Buffer)
	rlp.Encode(buf, params)

	return buf.Bytes()
}

func MovePlayer(postionMove *big.Int) []byte {
	methodId := crypto.Keccak256([]byte("movePlayer(uint256)"))[:4]
	paramValue := math.U256Bytes(postionMove)
	input := append(methodId, paramValue...)
	return input
}
