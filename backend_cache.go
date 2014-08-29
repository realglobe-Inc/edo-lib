package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// backend をキャッシュでラップする。

type cachingIdProviderBackend struct {
	IdProviderBackend
	cache  []*IdProvider
	caStmp *Stamp
}

func NewCachingIdProviderBackend(backend IdProviderBackend) IdProviderBackend {
	return &cachingIdProviderBackend{IdProviderBackend: backend}
}

func (reg *cachingIdProviderBackend) StampedIdProviders(caStmp *Stamp) ([]*IdProvider, *Stamp, error) {
	if reg.caStmp != nil && !time.Now().After(reg.caStmp.ExpiDate) {
		return reg.cache, reg.caStmp, nil
	}

	// キャッシュは有効期限切れ。

	idps, newCaStmp, err := reg.IdProviderBackend.StampedIdProviders(reg.caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	if newCaStmp == nil {
		return nil, nil, nil
	}

	// あった。

	log.Err("China ", newCaStmp)
	reg.caStmp = newCaStmp

	if idps != nil {
		reg.cache = idps
	}

	log.Err("Debu ", caStmp)
	if caStmp != nil && !newCaStmp.Date.After(caStmp.Date) && caStmp.Digest == newCaStmp.Digest {
		return nil, newCaStmp, nil
	}
	log.Err("Ero ", reg.cache)
	return reg.cache, newCaStmp, nil
}
