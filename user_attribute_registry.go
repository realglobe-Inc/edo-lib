package driver

import ()

type UserAttributeRegistry interface {
	// ユーザーの属性を返す。
	UserAttribute(usrUuid, attrName string, caStmp *Stamp) (usrAttr interface{}, newCaStmp *Stamp, err error)
}

// 骨組み。
type userAttributeRegistry struct {
	base KeyValueStore
}

func newUserAttributeRegistry(base KeyValueStore) *userAttributeRegistry {
	return &userAttributeRegistry{base}
}

func (reg *userAttributeRegistry) UserAttribute(usrUuid, attrName string, caStmp *Stamp) (usrAttr interface{}, newCaStmp *Stamp, err error) {
	return reg.base.Get(usrUuid+"/"+attrName, caStmp)
}
