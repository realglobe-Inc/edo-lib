package driver

import (
	"crypto/rsa"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

type TaKeyProvider interface {
	// サービスの公開鍵を返す。
	ServiceKey(servUuid string, caStmp *Stamp) (servKey *rsa.PublicKey, newCaStmp *Stamp, err error)
}

// 骨組み。
type taKeyProvider struct {
	base KeyValueStore
}

func newTaKeyProvider(base KeyValueStore) *taKeyProvider {
	return &taKeyProvider{base}
}

func (reg *taKeyProvider) ServiceKey(servUuid string, caStmp *Stamp) (servKey *rsa.PublicKey, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get(servUuid, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if value == nil || value == "" {
		return nil, newCaStmp, nil
	}

	return value.(*rsa.PublicKey), newCaStmp, nil
}
