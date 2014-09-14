package driver

import (
	"strconv"
	"time"
)

// メモリ上で完結する。デバッグ用。

// 非キャッシュ用。
type MemoryIdProviderLister struct {
	idps []*IdProvider
}

func NewMemoryIdProviderLister() *MemoryIdProviderLister {
	return &MemoryIdProviderLister{[]*IdProvider{}}
}

func (reg *MemoryIdProviderLister) IdProviders() ([]*IdProvider, error) {
	return reg.idps, nil
}
func (reg *MemoryIdProviderLister) SetIdProviders(idps []*IdProvider) {
	reg.idps = idps
}

// キャッシュ用。
type MemoryDatedIdProviderLister struct {
	*MemoryIdProviderLister
	stmp    *Stamp
	expiDur time.Duration
}

func NewMemoryDatedIdProviderLister(expiDur time.Duration) *MemoryDatedIdProviderLister {
	return &MemoryDatedIdProviderLister{NewMemoryIdProviderLister(), &Stamp{Date: time.Now(), Digest: strconv.Itoa(0)}, expiDur}
}

func (reg *MemoryDatedIdProviderLister) StampedIdProviders(caStmp *Stamp) ([]*IdProvider, *Stamp, error) {
	newCaStmp := &Stamp{Date: reg.stmp.Date, ExpiDate: time.Now().Add(reg.expiDur), Digest: reg.stmp.Digest}

	if caStmp == nil || caStmp.Date.Before(reg.stmp.Date) || caStmp.Digest != reg.stmp.Digest {
		idps, _ := reg.IdProviders()
		return idps, newCaStmp, nil
	}

	return nil, newCaStmp, nil
}
func (reg *MemoryDatedIdProviderLister) SetIdProviders(idps []*IdProvider) {
	reg.MemoryIdProviderLister.SetIdProviders(idps)
	dig, _ := strconv.Atoi(reg.stmp.Digest)
	reg.stmp = &Stamp{Date: time.Now(), Digest: strconv.Itoa(dig + 1)}
}
