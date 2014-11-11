package driver

import (
	"crypto/rsa"
	"time"
)

type MemoryTaKeyProvider struct {
	base KeyValueStore
}

// スレッドセーフ。
func NewMemoryTaKeyProvider(expiDur time.Duration) *MemoryTaKeyProvider {
	return &MemoryTaKeyProvider{NewMemoryKeyValueStore(expiDur)}
}

func (reg *MemoryTaKeyProvider) ServiceKey(servUuid string, caStmp *Stamp) (servKey *rsa.PublicKey, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get(servUuid, caStmp)
	if value != nil {
		servKey = value.(*rsa.PublicKey)
	}
	return servKey, newCaStmp, err
}

func (reg *MemoryTaKeyProvider) AddServiceKey(servUuid string, servKey *rsa.PublicKey) {
	reg.base.Put(servUuid, servKey)
}

func (reg *MemoryTaKeyProvider) RemoveServiceKey(servUuid string) {
	reg.base.Remove(servUuid)
}
