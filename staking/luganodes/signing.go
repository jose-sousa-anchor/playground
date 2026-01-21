package luganodes

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func SignMessage(message []byte, privateKeyHex string) (string, error) {
	key, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	signature, err := crypto.Sign(crypto.Keccak256(message), key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("0x%x", signature), nil
}
