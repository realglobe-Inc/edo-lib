package driver

import ()

// ID プロバイダのユーザー属性取得用 URI を返す。
type IdProviderRegistry interface {
	IdProviderQueryUri(idpUuid string) (queryUri string, err error)
}

// ID プロバイダのユーザー属性取得用 URI を返す。キャッシュ用。
type DatedIdProviderRegistry interface {
	StampedIdProviderQueryUri(idpUuid string, caStmp *Stamp) (queryUri string, newCaStmp *Stamp, err error)
}

// 非キャッシュ用。
type idProviderRegistry struct {
	keyValueStore
}

func newIdProviderRegistry(base keyValueStore) *idProviderRegistry {
	return &idProviderRegistry{base}
}

func (reg *idProviderRegistry) IdProviderQueryUri(idpUuid string) (queryUri string, err error) {
	val, err := reg.get(idpUuid)
	if val != nil && val != "" {
		queryUri = val.(string)
	}
	return queryUri, err
}

// キャッシュ用。
type datedIdProviderRegistry struct {
	datedKeyValueStore
}

func newDatedIdProviderRegistry(base datedKeyValueStore) *datedIdProviderRegistry {
	return &datedIdProviderRegistry{base}
}

func (reg *datedIdProviderRegistry) StampedIdProviderQueryUri(idpUuid string, caStmp *Stamp) (queryUri string, newCaStmp *Stamp, err error) {
	val, newCaStmp, err := reg.stampedGet(idpUuid, caStmp)
	if val != nil && val != "" {
		queryUri = val.(string)
	}
	return queryUri, newCaStmp, err
}
