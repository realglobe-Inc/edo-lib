package driver

import ()

// ID プロバイダの属性を返す。
type IdProviderAttributeRegistry interface {
	IdProviderAttribute(idpUuid, attrName string) (idpAttr interface{}, err error)
}

// ID プロバイダのユーザー属性取得用 URI を返す。キャッシュ用。
type DatedIdProviderAttributeRegistry interface {
	StampedIdProviderAttribute(idpUuid, attrName string, caStmp *Stamp) (idpAttr interface{}, newCaStmp *Stamp, err error)
}

func idProviderAttributeKey(idpUuid, attrName string) string {
	return idpUuid + "/" + attrName
}

// 非キャッシュ用。
type idProviderRegistry struct {
	keyValueStore
}

func newIdProviderAttributeRegistry(base keyValueStore) *idProviderRegistry {
	return &idProviderRegistry{base}
}

func (reg *idProviderRegistry) IdProviderAttribute(idpUuid, attrName string) (idpAttr interface{}, err error) {
	return reg.get(idProviderAttributeKey(idpUuid, attrName))
}

// キャッシュ用。
type datedIdProviderAttributeRegistry struct {
	datedKeyValueStore
}

func newDatedIdProviderAttributeRegistry(base datedKeyValueStore) *datedIdProviderAttributeRegistry {
	return &datedIdProviderAttributeRegistry{base}
}

func (reg *datedIdProviderAttributeRegistry) StampedIdProviderAttribute(idpUuid, attrName string, caStmp *Stamp) (idpAttr interface{}, newCaStmp *Stamp, err error) {
	return reg.stampedGet(idProviderAttributeKey(idpUuid, attrName), caStmp)
}
