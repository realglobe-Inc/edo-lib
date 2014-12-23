package util

import (
	"encoding/base64"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

func SecureRandomString(length int) (string, error) {
	buff, err := SecureRandomBytes((length*6 + 7) / 8)
	if err != nil {
		return "", erro.Wrap(err)
	}
	return base64.URLEncoding.EncodeToString(buff)[:length], nil
}
