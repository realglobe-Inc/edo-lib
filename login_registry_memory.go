package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

type MemoryLoginRegistry struct {
	base KeyValueStore
}

// スレッドセーフ。
func NewMemoryLoginRegistry(expiDur time.Duration) *MemoryLoginRegistry {
	return &MemoryLoginRegistry{NewMemoryKeyValueStore(expiDur)}
}

func (reg *MemoryLoginRegistry) User(accToken string, caStmp *Stamp) (addr string, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get(accToken, caStmp)
	if err != nil {
		return "", nil, erro.Wrap(err)
	} else if value == nil || value == "" {
		return "", newCaStmp, nil
	}
	return value.(string), newCaStmp, nil
}

func (reg *MemoryLoginRegistry) AddUser(accToken string, addr string) {
	reg.base.Put(accToken, addr)
}

func (reg *MemoryLoginRegistry) RemoveUser(accToken string) {
	reg.base.Remove(accToken)
}
