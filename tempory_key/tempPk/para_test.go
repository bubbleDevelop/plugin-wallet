package tempPk

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
)

func TestIsTxTxBehalfSignature(t *testing.T) {
	addValue := big.NewInt(-1234567890123456)
	addValueBytes, err := rlp.EncodeToBytes(addValue)
	if nil != err {
		panic(err)
	}
	var bigInt *big.Int
	if err := rlp.DecodeBytes(addValueBytes, &bigInt); nil != err {
		panic(err)
	}
}
