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
