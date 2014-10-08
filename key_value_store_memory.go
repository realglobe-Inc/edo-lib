package driver

import (
	"strconv"
	"time"
)

type memoryKeyValueStore struct {
	keyToValue map[string]interface{}
	keyToStmp  map[string]*Stamp
	expiDur    time.Duration
}

// スレッドセーフ。
func NewMemoryKeyValueStore(expiDur time.Duration) KeyValueStore {
	return newSynchronizedKeyValueStore(newMemoryKeyValueStore(expiDur))
}

// スレッドセーフではない。
func newMemoryKeyValueStore(expiDur time.Duration) *memoryKeyValueStore {
	return &memoryKeyValueStore{
		map[string]interface{}{},
		map[string]*Stamp{},
		expiDur,
	}
}

func (reg *memoryKeyValueStore) Get(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	stmp := reg.keyToStmp[key]
	if stmp == nil {
		return nil, nil, nil
	}
	newCaStmp = &Stamp{Date: stmp.Date, ExpiDate: time.Now().Add(reg.expiDur), Digest: stmp.Digest}

	if caStmp == nil || caStmp.Date.Before(stmp.Date) || caStmp.Digest != stmp.Digest {
		value, _ = reg.keyToValue[key]
		return value, newCaStmp, nil
	}

	return nil, newCaStmp, nil
}

func (reg *memoryKeyValueStore) Put(key string, value interface{}) (newCaStmp *Stamp, err error) {
	reg.keyToValue[key] = value

	newCaStmp = &Stamp{Date: time.Now(), Digest: strconv.FormatInt(time.Now().UnixNano(), 10)}
	reg.keyToStmp[key] = newCaStmp
	return newCaStmp, nil
}

func (reg *memoryKeyValueStore) Remove(key string) error {
	delete(reg.keyToValue, key)
	delete(reg.keyToStmp, key)
	return nil
}
