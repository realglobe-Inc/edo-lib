package driver

import (
	"time"
)

type MemoryIdpLister struct {
	base KeyValueStore
}

// スレッドセーフ。
func NewMemoryIdpLister(expiDur time.Duration) *MemoryIdpLister {
	return &MemoryIdpLister{NewMemoryKeyValueStore(expiDur)}
}

func (reg *MemoryIdpLister) IdProviders(caStmp *Stamp) (idps []*IdProvider, newCaStmp *Stamp, err error) {
	value, newCaStmp, _ := reg.base.Get("list", caStmp)
	if value == nil {
		return nil, newCaStmp, nil
	}
	return value.([]*IdProvider), newCaStmp, nil
}

func (reg *MemoryIdpLister) SetIdProviders(idps []*IdProvider) {
	reg.base.Put("list", idps)
}
