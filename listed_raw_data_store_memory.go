package driver

import (
	"time"
)

type memoryListedRawDataStore memoryListedKeyValueStore

// スレッドセーフ。
func NewMemoryListedRawDataStore(staleDur, expiDur time.Duration) ListedRawDataStore {
	return newSynchronizedListedRawDataStore(newMemoryListedRawDataStore(staleDur, expiDur))
}

// スレッドセーフではない。
func newMemoryListedRawDataStore(staleDur, expiDur time.Duration) *memoryListedRawDataStore {
	return (*memoryListedRawDataStore)(newMemoryListedKeyValueStore(staleDur, expiDur))
}

func (reg *memoryListedRawDataStore) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
	return ((*memoryListedKeyValueStore)(reg)).Keys(caStmp)
}

func (reg *memoryListedRawDataStore) Get(key string, caStmp *Stamp) (data []byte, newCaStmp *Stamp, err error) {
	val, newCaStmp, err := ((*memoryListedKeyValueStore)(reg)).Get(key, caStmp)
	if val == nil {
		return nil, newCaStmp, err
	}
	return val.([]byte), newCaStmp, nil
}

func (reg *memoryListedRawDataStore) Put(key string, data []byte) (*Stamp, error) {
	return ((*memoryListedKeyValueStore)(reg)).Put(key, data)
}

func (reg *memoryListedRawDataStore) Remove(key string) error {
	return ((*memoryListedKeyValueStore)(reg)).Remove(key)
}
