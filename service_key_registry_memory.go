package driver

import (
	"crypto/rsa"
	"time"
)

type MemoryServiceKeyRegistry struct {
	base KeyValueStore
}

// スレッドセーフ。
func NewMemoryServiceKeyRegistry(expiDur time.Duration) *MemoryServiceKeyRegistry {
	return &MemoryServiceKeyRegistry{NewMemoryKeyValueStore(expiDur)}
}

func (reg *MemoryServiceKeyRegistry) ServiceKey(servUuid string, caStmp *Stamp) (servKey *rsa.PublicKey, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get(servUuid, caStmp)
	if value != nil {
		servKey = value.(*rsa.PublicKey)
	}
	return servKey, newCaStmp, err
}

func (reg *MemoryServiceKeyRegistry) AddServiceKey(servUuid string, servKey *rsa.PublicKey) {
	reg.base.Put(servUuid, servKey)
}

func (reg *MemoryServiceKeyRegistry) RemoveServiceKey(servUuid string) {
	reg.base.Remove(servUuid)
}
