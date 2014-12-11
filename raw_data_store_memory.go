package driver

import (
	"time"
)

type memoryRawDataStore memoryKeyValueStore

// スレッドセーフ。
func NewMemoryRawDataStore(staleDur, expiDur time.Duration) RawDataStore {
	return newSynchronizedRawDataStore(newMemoryRawDataStore(staleDur, expiDur))
}

// スレッドセーフではない。
func newMemoryRawDataStore(staleDur, expiDur time.Duration) *memoryRawDataStore {
	return (*memoryRawDataStore)(newMemoryKeyValueStore(staleDur, expiDur))
}

func (reg *memoryRawDataStore) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
	return ((*memoryKeyValueStore)(reg)).Keys(caStmp)
}

func (reg *memoryRawDataStore) Get(key string, caStmp *Stamp) (data []byte, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := ((*memoryKeyValueStore)(reg)).Get(key, caStmp)
	if value == nil {
		return nil, newCaStmp, err
	}
	return value.([]byte), newCaStmp, nil
}

func (reg *memoryRawDataStore) Put(key string, data []byte) (*Stamp, error) {
	return ((*memoryKeyValueStore)(reg)).Put(key, data)
}

func (reg *memoryRawDataStore) Remove(key string) error {
	return ((*memoryKeyValueStore)(reg)).Remove(key)
}
