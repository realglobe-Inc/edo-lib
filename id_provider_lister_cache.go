package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// キャッシュする。

// キャッシュ用。
type cachingDatedIdProviderLister struct {
	DatedIdProviderLister
	cache  []*IdProvider
	caStmp *Stamp
}

func NewCachingDatedIdProviderLister(backend DatedIdProviderLister) DatedIdProviderLister {
	return &cachingDatedIdProviderLister{DatedIdProviderLister: backend}
}

func (reg *cachingDatedIdProviderLister) StampedIdProviders(caStmp *Stamp) ([]*IdProvider, *Stamp, error) {
	if reg.caStmp != nil && !time.Now().After(reg.caStmp.ExpiDate) {
		return reg.cache, reg.caStmp, nil
	}

	// キャッシュは有効期限切れ。

	idps, newCaStmp, err := reg.DatedIdProviderLister.StampedIdProviders(reg.caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	if newCaStmp == nil {
		return nil, nil, nil
	}

	// あった。

	reg.caStmp = newCaStmp

	if idps != nil {
		reg.cache = idps
	}

	if caStmp != nil && !newCaStmp.Date.After(caStmp.Date) && caStmp.Digest == newCaStmp.Digest {
		return nil, newCaStmp, nil
	}
	return reg.cache, newCaStmp, nil
}
