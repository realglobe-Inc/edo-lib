package driver

import (
	"time"
)

type MemoryIdProviderAttributeRegistry struct {
	base KeyValueStore
}

// スレッドセーフ。
func NewMemoryIdProviderAttributeRegistry(expiDur time.Duration) *MemoryIdProviderAttributeRegistry {
	return &MemoryIdProviderAttributeRegistry{NewMemoryKeyValueStore(expiDur)}
}

func (reg *MemoryIdProviderAttributeRegistry) IdProviderAttribute(idpUuid, attrName string, caStmp *Stamp) (idpAttr interface{}, newCaStmp *Stamp, err error) {
	return reg.base.Get(idpUuid+"/"+attrName, caStmp)
}

func (reg *MemoryIdProviderAttributeRegistry) AddIdProviderAttribute(idpUuid, attrName string, idpAttr interface{}) {
	reg.base.Put(idpUuid+"/"+attrName, idpAttr)
}

func (reg *MemoryIdProviderAttributeRegistry) RemoveIdProviderAttribute(idpUuid, attrName string) {
	reg.base.Remove(idpUuid + "/" + attrName)
}
