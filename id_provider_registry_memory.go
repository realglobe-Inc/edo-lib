package driver

import (
	"time"
)

// メモリ上で完結する。デバッグ用。

// 非キャッシュ用。
type MemoryIdProviderRegistry struct {
	keyValueStore
}

func NewMemoryIdProviderRegistry() *MemoryIdProviderRegistry {
	return &MemoryIdProviderRegistry{newSynchronizedKeyValueStore(newMemoryKeyValueStore())}
}

func (reg *MemoryIdProviderRegistry) IdProviderQueryUri(idpUuid string) (queryUri string, err error) {
	val, err := reg.get(idpUuid)
	if val != nil && val != "" {
		queryUri = val.(string)
	}
	return queryUri, err
}

func (reg *MemoryIdProviderRegistry) AddIdProviderQueryUri(idpUuid, queryUri string) {
	reg.put(idpUuid, queryUri)
}

func (reg *MemoryIdProviderRegistry) RemoveIdProvider(idpUuid string) {
	reg.remove(idpUuid)
}

// キャッシュ用。
type MemoryDatedIdProviderRegistry struct {
	datedKeyValueStore
}

func NewMemoryDatedIdProviderRegistry(expiDur time.Duration) *MemoryDatedIdProviderRegistry {
	return &MemoryDatedIdProviderRegistry{newSynchronizedDatedKeyValueStore(newMemoryDatedKeyValueStore(expiDur))}
}

func (reg *MemoryDatedIdProviderRegistry) StampedIdProviderQueryUri(idpUuid string, caStmp *Stamp) (queryUri string, newCaStmp *Stamp, err error) {
	val, newCaStmp, err := reg.stampedGet(idpUuid, caStmp)
	if val != nil && val != "" {
		queryUri = val.(string)
	}
	return queryUri, newCaStmp, err
}

func (reg *MemoryDatedIdProviderRegistry) AddIdProviderQueryUri(idpUuid, queryUri string) {
	reg.stampedPut(idpUuid, queryUri)
}

func (reg *MemoryDatedIdProviderRegistry) RemoveIdProviderQueryUri(idpUuid string) {
	reg.remove(idpUuid)
}
