package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

type MemoryUserNameIndex struct {
	base KeyValueStore
}

// スレッドセーフ。
func NewMemoryUserNameIndex(expiDur time.Duration) *MemoryUserNameIndex {
	return &MemoryUserNameIndex{NewMemoryKeyValueStore(expiDur)}
}

func (reg *MemoryUserNameIndex) UserUuid(usrName string, caStmp *Stamp) (usrUuid string, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get(usrName, caStmp)
	if err != nil {
		return "", nil, erro.Wrap(err)
	} else if value == nil || value == "" {
		return "", newCaStmp, nil
	}
	return value.(string), newCaStmp, err
}

func (reg *MemoryUserNameIndex) AddUserUuid(usrName, usrUuid string) {
	reg.base.Put(usrName, usrUuid)
}

func (reg *MemoryUserNameIndex) RemoveIdProvider(usrName string) {
	reg.base.Remove(usrName)
}
