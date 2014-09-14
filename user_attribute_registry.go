package driver

import ()

// ユーザーの属性を返す。
type UserAttributeRegistry interface {
	UserAttribute(usrUuid, attrName string) (usrAttr interface{}, err error)
}

// 非キャッシュ用。
type userAttributeRegistry struct {
	keyValueStore
}

func newUserAttributeRegistry(base keyValueStore) *userAttributeRegistry {
	return &userAttributeRegistry{base}
}

func userAttributeKey(usrUuid, attrName string) string {
	return usrUuid + "/" + attrName
}

func (reg *userAttributeRegistry) UserAttribute(usrUuid, attrName string) (usrAttr interface{}, err error) {
	return reg.get(userAttributeKey(usrUuid, attrName))
}
