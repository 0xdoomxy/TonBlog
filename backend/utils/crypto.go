package utils

import (
	"crypto/ecdsa"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func Ecrecover(message string, signature string) (address string, err error) {
	var t = new(TypedData)
	err = json.Unmarshal([]byte(message), t)
	if err != nil {
		return
	}
	var hash []byte
	hash, message, err = TypedDataAndHash(*t)
	if err != nil {
		return
	}
	var pk *ecdsa.PublicKey
	byt := common.Hex2Bytes(signature[2:])
	byt[64] -= 27
	pk, err = crypto.SigToPub(hash, byt)
	if err != nil {
		return
	}
	address = crypto.PubkeyToAddress(*pk).String()
	return
}
