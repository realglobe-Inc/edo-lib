package driver

import (
	"time"
)

type MemoryIdpAttributeProvider struct {
	base KeyValueStore
}

// スレッドセーフ。
func NewMemoryIdpAttributeProvider(expiDur time.Duration) *MemoryIdpAttributeProvider {
	return &MemoryIdpAttributeProvider{NewMemoryKeyValueStore(expiDur)}
}

func (reg *MemoryIdpAttributeProvider) IdProviderAttribute(idpUuid, attrName string, caStmp *Stamp) (idpAttr interface{}, newCaStmp *Stamp, err error) {
	return reg.base.Get(idpUuid+"/"+attrName, caStmp)
}

func (reg *MemoryIdpAttributeProvider) AddIdProviderAttribute(idpUuid, attrName string, idpAttr interface{}) {
	reg.base.Put(idpUuid+"/"+attrName, idpAttr)
}

func (reg *MemoryIdpAttributeProvider) RemoveIdProviderAttribute(idpUuid, attrName string) {
	reg.base.Remove(idpUuid + "/" + attrName)
}
