package shortener

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/itchyny/base58-go"
)

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) (string, error) {
	encoded, err := base58.BitcoinEncoding.Encode(bytes)
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func GenerateShortLink(originUrl string, userId string) (string, error) {
	urlHashByte := sha256Of(originUrl + userId)
	generateNumber := new(big.Int).SetBytes(urlHashByte).Uint64()
	result, err := base58Encoded([]byte(fmt.Sprintf("%d", generateNumber)))
	if err != nil {
		return "", err
	}

	return result, nil
}
