package driver

import ()

// ID プロバイダ選択時に列挙する情報。
type IdProvider struct {
	Name string `json:"name" bson:"name"`
	Uuid string `json:"uuid" bson:"uuid"`
}

func (idp *IdProvider) String() string {
	return idp.Uuid + "," + idp.Name
}

// ID プロバイダの列挙。
type IdProviderLister interface {
	IdProviders() ([]*IdProvider, error)
}

// ID プロバイダの列挙。キャッシュ用。
type DatedIdProviderLister interface {
	StampedIdProviders(caStmp *Stamp) ([]*IdProvider, *Stamp, error)
}
