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

func (drv *memoryListedRawDataStore) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
	return ((*memoryListedKeyValueStore)(drv)).Keys(caStmp)
}

func (drv *memoryListedRawDataStore) Get(key string, caStmp *Stamp) (data []byte, newCaStmp *Stamp, err error) {
	val, newCaStmp, err := ((*memoryListedKeyValueStore)(drv)).Get(key, caStmp)
	if val == nil {
		return nil, newCaStmp, err
	}
	return val.([]byte), newCaStmp, nil
}

func (drv *memoryListedRawDataStore) Put(key string, data []byte) (*Stamp, error) {
	return ((*memoryListedKeyValueStore)(drv)).Put(key, data)
}

func (drv *memoryListedRawDataStore) Remove(key string) error {
	return ((*memoryListedKeyValueStore)(drv)).Remove(key)
}

func (drv *memoryListedRawDataStore) Close() error {
	return ((*memoryListedKeyValueStore)(drv)).Close()
}
