package driver

import ()

type IdProviderAttributeRegistry interface {
	// ID プロバイダの属性を返す。
	IdProviderAttribute(idpUuid, attrName string, caStmp *Stamp) (idpAttr interface{}, newCaStmp *Stamp, err error)
}

// 骨組み。
// バックエンドで ID プロバイダの属性ごとに保存。
type idProviderAttributeRegistry struct {
	base KeyValueStore
}

func newIdProviderAttributeRegistry(base KeyValueStore) *idProviderAttributeRegistry {
	return &idProviderAttributeRegistry{base}
}

func (reg *idProviderAttributeRegistry) IdProviderAttribute(idpUuid, attrName string, caStmp *Stamp) (idpAttr interface{}, newCaStmp *Stamp, err error) {
	return reg.base.Get(idpUuid+"/"+attrName, caStmp)
}
