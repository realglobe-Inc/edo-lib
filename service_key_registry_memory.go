package driver

import (
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

func (reg *MemoryServiceKeyRegistry) ServiceKey(servUuid string) (servKey string, err error) {
	val, err := reg.get(servUuid)
	if val != nil && val != "" {
		servKey = val.(string)
	}
	return servKey, err
}

func (reg *MemoryServiceKeyRegistry) AddServiceKey(servUuid, servKey string) {
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

func (reg *MemoryDatedServiceKeyRegistry) StampedServiceKey(servUuid string, caStmp *Stamp) (servKey string, newCaStmp *Stamp, err error) {
	val, newCaStmp, err := reg.stampedGet(servUuid, caStmp)
	if val != nil && val != "" {
		servKey = val.(string)
	}
	return servKey, newCaStmp, err
}

func (reg *MemoryDatedServiceKeyRegistry) AddServiceKey(servUuid, servKey string) {
	reg.stampedPut(servUuid, servKey)
}

func (reg *MemoryDatedServiceKeyRegistry) RemoveServiceKey(servUuid string) {
	reg.remove(servUuid)
}
