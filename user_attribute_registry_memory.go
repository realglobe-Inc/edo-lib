package driver

import (
	"time"
)

type MemoryUserAttributeRegistry struct {
	KeyValueStore
}

// スレッドセーフ。
func NewMemoryUserAttributeRegistry(expiDur time.Duration) *MemoryUserAttributeRegistry {
	return &MemoryUserAttributeRegistry{NewMemoryKeyValueStore(expiDur)}
}

func (reg *MemoryUserAttributeRegistry) UserAttribute(usrUuid, attrName string, caStmp *Stamp) (usrAttr interface{}, newCaStmp *Stamp, err error) {
	return reg.Get(usrUuid+"/"+attrName, caStmp)
}

func (reg *MemoryUserAttributeRegistry) AddUserAttribute(usrUuid, attrName string, usrAttr interface{}) {
	reg.Put(usrUuid+"/"+attrName, usrAttr)
}

func (reg *MemoryUserAttributeRegistry) RemoveIdProvider(usrUuid, attrName string) {
	reg.Remove(usrUuid + "/" + attrName)
}
