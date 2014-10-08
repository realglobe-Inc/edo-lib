package driver

import (
	"time"
)

type MemoryIdProviderLister struct {
	base KeyValueStore
}

// スレッドセーフ。
func NewMemoryIdProviderLister(expiDur time.Duration) *MemoryIdProviderLister {
	return &MemoryIdProviderLister{NewMemoryKeyValueStore(expiDur)}
}

func (reg *MemoryIdProviderLister) IdProviders(caStmp *Stamp) (idps []*IdProvider, newCaStmp *Stamp, err error) {
	value, newCaStmp, _ := reg.base.Get("list", caStmp)
	if value == nil {
		return nil, newCaStmp, nil
	}
	return value.([]*IdProvider), newCaStmp, nil
}

func (reg *MemoryIdProviderLister) SetIdProviders(idps []*IdProvider) {
	reg.base.Put("list", idps)
}
