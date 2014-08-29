package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// backend をキャッシュでラップする。

type cachingIdProviderRegistry struct {
	IdProviderBackend
	cache  []*IdProvider
	caStmp *Stamp
}

func NewCachingIdProviderRegistry(backend IdProviderBackend) IdProviderRegistry {
	return &cachingIdProviderRegistry{IdProviderBackend: backend}
}

func (reg *cachingIdProviderRegistry) IdProviders() ([]*IdProvider, error) {
	if reg.caStmp != nil && reg.caStmp.ExpiDate.After(time.Now()) {
		return reg.cache, nil
	}

	// キャッシュは有効期限切れ。

	idps, stmp, err := reg.StampedIdProviders(reg.caStmp)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	if stmp == nil {
		return nil, nil
	}

	// あった。

	reg.caStmp = stmp

	if idps == nil {
		return reg.cache, nil
	}

	// 新しくなってた。

	reg.cache = idps

	return idps, nil
}
