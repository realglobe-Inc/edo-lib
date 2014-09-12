package driver

import (
	"crypto/rsa"
	"time"
)

// メモリ上で完結する。デバッグ用。

// 非キャッシュ用。
type MemoryServiceKeyRegistry struct {
	keyValueStore
}

func NewMemoryServiceKeyRegistry() *MemoryServiceKeyRegistry {
	return &MemoryServiceKeyRegistry{newSynchronizedKeyValueStore(newMemoryKeyValueStore())}
}

func (reg *MemoryServiceKeyRegistry) ServiceKey(servUuid string) (servKey *rsa.PublicKey, err error) {
	val, err := reg.get(servUuid)
	if val != nil {
		servKey = val.(*rsa.PublicKey)
	}
	return servKey, err
}

func (reg *MemoryServiceKeyRegistry) AddServiceKey(servUuid string, servKey *rsa.PublicKey) {
	reg.put(servUuid, servKey)
}

func (reg *MemoryServiceKeyRegistry) RemoveServiceKey(servUuid string) {
	reg.remove(servUuid)
}

// キャッシュ用。
type MemoryDatedServiceKeyRegistry struct {
	datedKeyValueStore
}

func NewMemoryDatedServiceKeyRegistry(expiDur time.Duration) *MemoryDatedServiceKeyRegistry {
	return &MemoryDatedServiceKeyRegistry{newSynchronizedDatedKeyValueStore(newMemoryDatedKeyValueStore(expiDur))}
}

func (reg *MemoryDatedServiceKeyRegistry) StampedServiceKey(servUuid string, caStmp *Stamp) (servKey *rsa.PublicKey, newCaStmp *Stamp, err error) {
	val, newCaStmp, err := reg.stampedGet(servUuid, caStmp)
	if val != nil {
		servKey = val.(*rsa.PublicKey)
	}
	return servKey, newCaStmp, err
}

func (reg *MemoryDatedServiceKeyRegistry) AddServiceKey(servUuid string, servKey *rsa.PublicKey) {
	reg.stampedPut(servUuid, servKey)
}

func (reg *MemoryDatedServiceKeyRegistry) RemoveServiceKey(servUuid string) {
	reg.remove(servUuid)
}
