package driver

import (
	"crypto/rsa"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

type ServiceKeyRegistry interface {
	// サービスの公開鍵を返す。
	ServiceKey(servUuid string, caStmp *Stamp) (servKey *rsa.PublicKey, newCaStmp *Stamp, err error)
}

// 骨組み。
type serviceKeyRegistry struct {
	base KeyValueStore
}

func newServiceKeyRegistry(base KeyValueStore) *serviceKeyRegistry {
	return &serviceKeyRegistry{base}
}

func (reg *serviceKeyRegistry) ServiceKey(servUuid string, caStmp *Stamp) (servKey *rsa.PublicKey, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get(servUuid, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if value == nil || value == "" {
		return nil, newCaStmp, nil
	}

	return value.(*rsa.PublicKey), newCaStmp, nil
}
