package driver

import ()

// サービスの公開鍵を返す。
type ServiceKeyRegistry interface {
	ServiceKey(servUuid string) (key string, err error)
}

// サービスの公開鍵を返す。キャッシュ用。
type DatedServiceKeyRegistry interface {
	StampedServiceKey(servUuid string, caStmp *Stamp) (key string, newCaStmp *Stamp, err error)
}
