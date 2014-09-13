package driver

import ()

// サービスの URI から UUID を引く。非キャッシュ用。
type ServiceExplorer interface {
	ServiceUuid(servUri string) (servUuid string, err error)
}

// サービスの URI から UUID を引く。キャッシュ用。
type DatedServiceExplorer interface {
	StampedServiceUuid(servUri string, caStmp *Stamp) (servUuid string, newCaStmp *Stamp, err error)
}
