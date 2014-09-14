package driver

import (
	"time"
)

// メモリ上で完結する。デバッグ用。

// 非キャッシュ用。
type MemoryIdProviderAttributeRegistry struct {
	keyValueStore
}

func NewMemoryIdProviderAttributeRegistry() *MemoryIdProviderAttributeRegistry {
	return &MemoryIdProviderAttributeRegistry{newSynchronizedKeyValueStore(newMemoryKeyValueStore())}
}

func (reg *MemoryIdProviderAttributeRegistry) IdProviderAttribute(idpUuid, attrName string) (idpAttr interface{}, err error) {
	return reg.get(idProviderAttributeKey(idpUuid, attrName))
}

func (reg *MemoryIdProviderAttributeRegistry) AddIdProviderAttribute(idpUuid, attrName string, idpAttr interface{}) {
	reg.put(idProviderAttributeKey(idpUuid, attrName), idpAttr)
}

func (reg *MemoryIdProviderAttributeRegistry) RemoveIdProvider(idpUuid, attrName string) {
	reg.remove(idProviderAttributeKey(idpUuid, attrName))
}

// キャッシュ用。
type MemoryDatedIdProviderAttributeRegistry struct {
	datedKeyValueStore
}

func NewMemoryDatedIdProviderAttributeRegistry(expiDur time.Duration) *MemoryDatedIdProviderAttributeRegistry {
	return &MemoryDatedIdProviderAttributeRegistry{newSynchronizedDatedKeyValueStore(newMemoryDatedKeyValueStore(expiDur))}
}

func (reg *MemoryDatedIdProviderAttributeRegistry) StampedIdProviderAttribute(idpUuid, attrName string, caStmp *Stamp) (idpAttr interface{}, newCaStmp *Stamp, err error) {
	return reg.stampedGet(idProviderAttributeKey(idpUuid, attrName), caStmp)
}

func (reg *MemoryDatedIdProviderAttributeRegistry) AddIdProviderAttribute(idpUuid, attrName string, idpAttr interface{}) {
	reg.stampedPut(idProviderAttributeKey(idpUuid, attrName), idpAttr)
}

func (reg *MemoryDatedIdProviderAttributeRegistry) RemoveIdProviderAttribute(idpUuid, attrName string) {
	reg.remove(idProviderAttributeKey(idpUuid, attrName))
}
