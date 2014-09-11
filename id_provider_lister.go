package driver

import ()

// ID プロバイダ選択時に列挙する情報。
type IdProvider struct {
	Uuid     string `json:"uuid"       bson:"uuid"`
	Name     string `json:"name"       bson:"name"`
	LoginUri string `json:"login_uri"  bson:"login_uri"`
}

func (idp *IdProvider) String() string {
	return idp.Uuid + "," + idp.Name + "," + idp.LoginUri
}

// ID プロバイダの列挙。
type IdProviderLister interface {
	IdProviders() ([]*IdProvider, error)
}

// ID プロバイダの列挙。キャッシュ用。
type DatedIdProviderLister interface {
	StampedIdProviders(caStmp *Stamp) ([]*IdProvider, *Stamp, error)
}
