package secrand

import (
	"crypto/rand"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io"
)

func Bytes(length int) ([]byte, error) {
	buff := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, buff); err != nil {
		return nil, erro.Wrap(err)
	}
	return buff, nil
}
