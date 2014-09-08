package driver

import (
	"strconv"
	"time"
)

// メモリ上で完結する。デバッグ用。

// 非キャッシュ用。
type MemoryIdProviderLister struct {
	idps map[string]*IdProvider
}

func NewMemoryIdProviderLister() *MemoryIdProviderLister {
	return &MemoryIdProviderLister{map[string]*IdProvider{}}
}

func (reg *MemoryIdProviderLister) IdProviders() ([]*IdProvider, error) {
	idps := []*IdProvider{}
	for _, idp := range reg.idps {
		idps = append(idps, idp)
	}
	return idps, nil
}
func (reg *MemoryIdProviderLister) AddIdProvider(idp *IdProvider) {
	reg.idps[idp.Uuid] = idp
}
func (reg *MemoryIdProviderLister) RemoveIdProvider(idpUuid string) {
	delete(reg.idps, idpUuid)
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
func (reg *MemoryDatedIdProviderLister) AddIdProvider(idp *IdProvider) {
	reg.MemoryIdProviderLister.AddIdProvider(idp)
	dig, _ := strconv.Atoi(reg.stmp.Digest)
	reg.stmp = &Stamp{Date: time.Now(), Digest: strconv.Itoa(dig + 1)}
}
func (reg *MemoryDatedIdProviderLister) RemoveIdProvider(idpUuid string) {
	reg.MemoryIdProviderLister.RemoveIdProvider(idpUuid)
	dig, _ := strconv.Atoi(reg.stmp.Digest)
	reg.stmp = &Stamp{Date: time.Now(), Digest: strconv.Itoa(dig + 1)}
}
