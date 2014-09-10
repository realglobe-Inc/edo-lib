package driver

import ()

// サービスの公開鍵を返す。
type ServiceKeyRegistry interface {
	ServiceKey(servUuid string) (servKey string, err error)
}

// サービスの公開鍵を返す。キャッシュ用。
type DatedServiceKeyRegistry interface {
	StampedServiceKey(servUuid string, caStmp *Stamp) (servKey string, newCaStmp *Stamp, err error)
}

// 非キャッシュ用。
type serviceKeyRegistry struct {
	keyValueStore
}

func newServiceKeyRegistry(base keyValueStore) *serviceKeyRegistry {
	return &serviceKeyRegistry{base}
}

func (reg *serviceKeyRegistry) ServiceKey(servUuid string) (servKey string, err error) {
	val, err := reg.get(servUuid)
	if val != nil && val != "" {
		servKey = val.(string)
	}
	return servKey, err
}

// キャッシュ用。
type datedServiceKeyRegistry struct {
	datedKeyValueStore
}

func newDatedServiceKeyRegistry(base datedKeyValueStore) *datedServiceKeyRegistry {
	return &datedServiceKeyRegistry{base}
}

func (reg *datedServiceKeyRegistry) StampedServiceKey(servUuid string, caStmp *Stamp) (servKey string, newCaStmp *Stamp, err error) {
	val, newCaStmp, err := reg.stampedGet(servUuid, caStmp)
	if val != nil && val != "" {
		servKey = val.(string)
	}
	return servKey, newCaStmp, err
}
