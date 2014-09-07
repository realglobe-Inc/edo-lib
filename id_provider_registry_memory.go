package driver

import (
	"strconv"
	"time"
)

// メモリ上で完結する。デバッグ用。

// 非キャッシュ用。
type MemoryIdProviderRegistry struct {
	idps map[string]string
}

func NewMemoryIdProviderRegistry() *MemoryIdProviderRegistry {
	return &MemoryIdProviderRegistry{map[string]string{}}
}

func (reg *MemoryIdProviderRegistry) IdProviderQueryUri(idpUuid string) (queryUri string, err error) {
	return reg.idps[idpUuid], nil
}
func (reg *MemoryIdProviderRegistry) AddIdProviderQueryUri(idpUuid, queryUri string) {
	reg.idps[idpUuid] = queryUri
}
func (reg *MemoryIdProviderRegistry) RemoveIdProvider(idpUuid string) {
	delete(reg.idps, idpUuid)
}

// キャッシュ用。
type MemoryDatedIdProviderRegistry struct {
	*MemoryIdProviderRegistry
	stmps   map[string]*Stamp
	expiDur time.Duration
}

func NewMemoryDatedIdProviderRegistry(expiDur time.Duration) *MemoryDatedIdProviderRegistry {
	return &MemoryDatedIdProviderRegistry{NewMemoryIdProviderRegistry(), map[string]*Stamp{}, expiDur}
}

func (reg *MemoryDatedIdProviderRegistry) StampedIdProviderQueryUri(idpUuid string, caStmp *Stamp) (queryUri string, newCaStmp *Stamp, err error) {
	stmp := reg.stmps[idpUuid]
	if stmp == nil {
		return "", nil, nil
	}
	newCaStmp = &Stamp{Date: stmp.Date, ExpiDate: time.Now().Add(reg.expiDur), Digest: stmp.Digest}

	if caStmp == nil || caStmp.Date.Before(stmp.Date) || caStmp.Digest != stmp.Digest {
		queryUri, _ = reg.IdProviderQueryUri(idpUuid)
		return queryUri, newCaStmp, nil
	}

	return "", newCaStmp, nil
}
func (reg *MemoryDatedIdProviderRegistry) AddIdProviderQueryUri(idpUuid, queryUri string) {
	reg.MemoryIdProviderRegistry.AddIdProviderQueryUri(idpUuid, queryUri)
	var dig int
	stmp := reg.stmps[idpUuid]
	if stmp == nil {
		dig = 0
	} else {
		dig, _ = strconv.Atoi(stmp.Digest)
	}
	newStmp := &Stamp{Date: time.Now(), Digest: strconv.Itoa(dig + 1)}
	reg.stmps[idpUuid] = newStmp
}
func (reg *MemoryDatedIdProviderRegistry) RemoveIdProviderQueryUri(idpUuid string) {
	reg.MemoryIdProviderRegistry.RemoveIdProvider(idpUuid)
	delete(reg.stmps, idpUuid)
}
