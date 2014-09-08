package driver

import (
	"strconv"
	"time"
)

// メモリ上で完結する。デバッグ用。

// 非キャッシュ用。
type MemoryServiceKeyRegistry struct {
	keys map[string]string
}

func NewMemoryServiceKeyRegistry() *MemoryServiceKeyRegistry {
	return &MemoryServiceKeyRegistry{map[string]string{}}
}

func (reg *MemoryServiceKeyRegistry) ServiceKey(servUuid string) (key string, err error) {
	return reg.keys[servUuid], nil
}
func (reg *MemoryServiceKeyRegistry) AddServiceKey(servUuid, key string) {
	reg.keys[servUuid] = key
}
func (reg *MemoryServiceKeyRegistry) RemoveServiceKey(servUuid string) {
	delete(reg.keys, servUuid)
}

// キャッシュ用。
type MemoryDatedServiceKeyRegistry struct {
	*MemoryServiceKeyRegistry
	stmps   map[string]*Stamp
	expiDur time.Duration
}

func NewMemoryDatedServiceKeyRegistry(expiDur time.Duration) *MemoryDatedServiceKeyRegistry {
	return &MemoryDatedServiceKeyRegistry{NewMemoryServiceKeyRegistry(), map[string]*Stamp{}, expiDur}
}

func (reg *MemoryDatedServiceKeyRegistry) StampedServiceKey(servUuid string, caStmp *Stamp) (key string, newCaStmp *Stamp, err error) {
	stmp := reg.stmps[servUuid]
	if stmp == nil {
		return "", nil, nil
	}
	newCaStmp = &Stamp{Date: stmp.Date, ExpiDate: time.Now().Add(reg.expiDur), Digest: stmp.Digest}

	if caStmp == nil || caStmp.Date.Before(stmp.Date) || caStmp.Digest != stmp.Digest {
		key, _ = reg.ServiceKey(servUuid)
		return key, newCaStmp, nil
	}

	return "", newCaStmp, nil
}
func (reg *MemoryDatedServiceKeyRegistry) AddServiceKey(servUuid, key string) {
	reg.MemoryServiceKeyRegistry.AddServiceKey(servUuid, key)
	var dig int
	stmp := reg.stmps[servUuid]
	if stmp == nil {
		dig = 0
	} else {
		dig, _ = strconv.Atoi(stmp.Digest)
	}
	newStmp := &Stamp{Date: time.Now(), Digest: strconv.Itoa(dig + 1)}
	reg.stmps[servUuid] = newStmp
}
func (reg *MemoryDatedServiceKeyRegistry) RemoveServiceKey(servUuid string) {
	reg.MemoryServiceKeyRegistry.RemoveServiceKey(servUuid)
	delete(reg.stmps, servUuid)
}
