package onchain

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

// 계정 생성 후 하나의 체인에서만 사용할 것을 권장함. 관리의 편의를 위해서
func CreateAccount() (string, string, error) {
	// 새 개인키 생성
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	// 개인키를 hex 문자열로 변환
	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))

	// 공개키로부터 주소 생성
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return address.Hex(), privateKeyHex, nil
}
