package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// ユーザー名からユーザー UUID を引く。
type UserNameIndex interface {
	UserUuid(usrName string, caStmp *Stamp) (usrUuid string, newCaStmp *Stamp, err error)
}

// 骨組み。
type userNameIndex struct {
	base KeyValueStore
}

func newUserNameIndex(base KeyValueStore) *userNameIndex {
	return &userNameIndex{base}
}

func (reg *userNameIndex) UserUuid(usrName string, caStmp *Stamp) (usrUuid string, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get(usrName, caStmp)
	if err != nil {
		return "", nil, erro.Wrap(err)
	} else if value == nil || value == "" {
		return "", newCaStmp, nil
	}
	return value.(string), newCaStmp, nil
}
