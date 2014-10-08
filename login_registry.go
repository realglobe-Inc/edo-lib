package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

type LoginRegistry interface {
	// ユーザー ID の取得。
	User(accToken string, caStmp *Stamp) (usrUuid string, newCaStmp *Stamp, err error)
}

// 骨組み。
type loginRegistry struct {
	base KeyValueStore
}

func newLoginRegistry(base KeyValueStore) *loginRegistry {
	return &loginRegistry{base}
}

func (reg *loginRegistry) User(accToken string, caStmp *Stamp) (usrUuid string, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get(accToken, caStmp)
	if err != nil {
		return "", nil, erro.Wrap(err)
	} else if value == nil || value == "" {
		return "", newCaStmp, nil
	}
	return value.(string), newCaStmp, nil
}
