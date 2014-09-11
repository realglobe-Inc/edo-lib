package driver

import ()

type MemoryUserAttributeRegistry struct {
	keyValueStore
}

func NewMemoryUserAttributeRegistry() *MemoryUserAttributeRegistry {
	return &MemoryUserAttributeRegistry{newSynchronizedKeyValueStore(newMemoryKeyValueStore())}
}

func (reg *MemoryUserAttributeRegistry) UserAttribute(usrUuid, attrName string) (usrAttr interface{}, err error) {
	return reg.get(userAttributeKey(usrUuid, attrName))
}

func (reg *MemoryUserAttributeRegistry) AddUserAttribute(usrUuid, attrName string, usrAttr interface{}) {
	reg.put(userAttributeKey(usrUuid, attrName), usrAttr)
}

func (reg *MemoryUserAttributeRegistry) RemoveIdProvider(usrUuid, attrName string) {
	reg.remove(userAttributeKey(usrUuid, attrName))
}
