package driver

import ()

type UserRegistry interface {
	// ユーザー情報を取得。
	Attribute(usrUuid, attrName string, caStmp *Stamp) (attr interface{}, newCaStmp *Stamp, err error)

	// ユーザー情報を変更。
	AddAttribute(usrUuid, attrName string, attr interface{}) (newCaStmp *Stamp, err error)
	// ユーザー情報を削除。
	RemoveAttribute(usrUuid, attrName string) error
}

// 骨組み。
type userRegistry struct {
	base KeyValueStore
}

func newUserRegistry(base KeyValueStore) *userRegistry {
	return &userRegistry{base}
}

func (reg *userRegistry) Attribute(usrUuid, attrName string, caStmp *Stamp) (attr interface{}, newCaStmp *Stamp, err error) {
	return reg.base.Get(usrUuid+"/"+attrName, caStmp)
}

func (reg *userRegistry) AddAttribute(usrUuid, attrName string, attr interface{}) (newCaStmp *Stamp, err error) {
	return reg.base.Put(usrUuid+"/"+attrName, attr)
}

func (reg *userRegistry) RemoveAttribute(usrUuid, attrName string) error {
	return reg.base.Remove(usrUuid + "/" + attrName)
}
