package driver

import (
	"strconv"
	"time"
)

// メモリ上で完結する。デバッグ用。

// 非キャッシュ用。
type memoryKeyValueStore struct {
	cont map[string]interface{}
}

func newMemoryKeyValueStore() *memoryKeyValueStore {
	return &memoryKeyValueStore{map[string]interface{}{}}
}

func (reg *memoryKeyValueStore) get(key string) (value interface{}, err error) {
	return reg.cont[key], nil
}

func (reg *memoryKeyValueStore) put(key string, value interface{}) error {
	reg.cont[key] = value
	return nil
}

func (reg *memoryKeyValueStore) remove(key string) error {
	delete(reg.cont, key)
	return nil
}

// キャッシュ用。
type memoryDatedKeyValueStore struct {
	*memoryKeyValueStore
	stmps   map[string]*Stamp
	expiDur time.Duration
}

func newMemoryDatedKeyValueStore(expiDur time.Duration) *memoryDatedKeyValueStore {
	return &memoryDatedKeyValueStore{newMemoryKeyValueStore(), map[string]*Stamp{}, expiDur}
}

func (reg *memoryDatedKeyValueStore) stampedGet(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	stmp := reg.stmps[key]
	if stmp == nil {
		return nil, nil, nil
	}
	newCaStmp = &Stamp{Date: stmp.Date, ExpiDate: time.Now().Add(reg.expiDur), Digest: stmp.Digest}

	if caStmp == nil || caStmp.Date.Before(stmp.Date) || caStmp.Digest != stmp.Digest {
		value, _ = reg.get(key)
		return value, newCaStmp, nil
	}

	return nil, newCaStmp, nil
}

func (reg *memoryDatedKeyValueStore) stampedPut(key string, value interface{}) (*Stamp, error) {
	reg.memoryKeyValueStore.put(key, value)
	var dig int
	stmp := reg.stmps[key]
	if stmp == nil {
		dig = 0
	} else {
		dig, _ = strconv.Atoi(stmp.Digest)
	}
	newStmp := &Stamp{Date: time.Now(), Digest: strconv.Itoa(dig + 1)}
	reg.stmps[key] = newStmp
	return newStmp, nil
}

func (reg *memoryDatedKeyValueStore) remove(key string) error {
	reg.memoryKeyValueStore.remove(key)
	delete(reg.stmps, key)
	return nil
}
