package driver

import ()

type IdpAttributeProvider interface {
	// ID プロバイダの属性を返す。
	IdProviderAttribute(idpUuid, attrName string, caStmp *Stamp) (idpAttr interface{}, newCaStmp *Stamp, err error)
}

// 骨組み。
// バックエンドで ID プロバイダの属性ごとに保存。
type idpAttributeProvider struct {
	base KeyValueStore
}

func newIdpAttributeProvider(base KeyValueStore) *idpAttributeProvider {
	return &idpAttributeProvider{base}
}

func (reg *idpAttributeProvider) IdProviderAttribute(idpUuid, attrName string, caStmp *Stamp) (idpAttr interface{}, newCaStmp *Stamp, err error) {
	return reg.base.Get(idpUuid+"/"+attrName, caStmp)
}
