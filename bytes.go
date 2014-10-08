package util

import (
	"crypto/rand"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"math/big"
)

func SecureRandomBytes(length int) ([]byte, error) {
	maxVal := big.NewInt(0).Lsh(big.NewInt(1), uint(length*8))
	value, err := rand.Int(rand.Reader, maxVal)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	buff := value.Bytes()

	// 上位 8 ビット以上が全部 0 だと短い値になり得るため。
	for len(buff) < length {
		buff = append([]byte{0}, buff...) // ビッグエンディアンなので、前に 0 を足さないと、一様性に影響が出る。
	}

	return buff, nil
}
