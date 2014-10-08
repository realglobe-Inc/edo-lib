package driver

import (
	"time"
)

type memoryRawDataStore memoryKeyValueStore

// スレッドセーフ。
func NewMemoryRawDataStore(expiDur time.Duration) RawDataStore {
	return newSynchronizedRawDataStore(newMemoryRawDataStore(expiDur))
}

// スレッドセーフではない。
func newMemoryRawDataStore(expiDur time.Duration) *memoryRawDataStore {
	return (*memoryRawDataStore)(newMemoryKeyValueStore(expiDur))
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
