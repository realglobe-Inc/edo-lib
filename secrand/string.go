package secrand

import (
	"encoding/base64"
	"github.com/realglobe-Inc/go-lib/erro"
)

func String(length int) (string, error) {
	buff, err := Bytes((length*6 + 7) / 8)
	if err != nil {
		return "", erro.Wrap(err)
	}
	return base64.URLEncoding.EncodeToString(buff)[:length], nil
}
