package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// backend をキャッシュでラップする。

type cachingIdProviderLister struct {
	DatedIdProviderLister
	cache  []*IdProvider
	caStmp *Stamp
}

func NewCachingIdProviderLister(backend DatedIdProviderLister) IdProviderLister {
	return &cachingIdProviderLister{DatedIdProviderLister: backend}
}

func (reg *cachingIdProviderLister) IdProviders() ([]*IdProvider, error) {
	if reg.caStmp != nil && reg.caStmp.ExpiDate.After(time.Now()) {
		return reg.cache, nil
	}

	// キャッシュは有効期限切れ。

	idps, newCaStmp, err := reg.StampedIdProviders(reg.caStmp)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	if newCaStmp == nil {
		return nil, nil
	}

	// あった。

	reg.caStmp = newCaStmp

	if idps == nil {
		return reg.cache, nil
	}

	// 新しくなってた。

	reg.cache = idps

	return idps, nil
}
