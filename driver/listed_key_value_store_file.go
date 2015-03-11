package driver

import (
	"github.com/realglobe-Inc/go-lib/erro"
	"time"
)

type Marshal func(interface{}) ([]byte, error)
type Unmarshal func([]byte) (interface{}, error)

type fileListedKeyValueStore struct {
	base ListedRawDataStore
	Marshal
	Unmarshal
}

// スレッドセーフ。
func NewFileListedKeyValueStore(path string, keyToPath, pathToKey func(string) string, marshal Marshal, unmarshal Unmarshal, staleDur, expiDur time.Duration) ListedKeyValueStore {
	return newSynchronizedListedKeyValueStore(newCachingListedKeyValueStore(newFileListedKeyValueStore(path, keyToPath, pathToKey, marshal, unmarshal, staleDur, expiDur)))
}

// スレッドセーフではない。
func newFileListedKeyValueStore(path string, keyToPath, pathToKey func(string) string, marshal Marshal, unmarshal Unmarshal, staleDur, expiDur time.Duration) *fileListedKeyValueStore {
	return &fileListedKeyValueStore{
		newFileListedRawDataStore(path, keyToPath, pathToKey, staleDur, expiDur),
		marshal, unmarshal,
	}
}

func (drv *fileListedKeyValueStore) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
	return drv.base.Keys(caStmp)
}

func (drv *fileListedKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	buff, newCaStmp, err := drv.base.Get(key, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if buff == nil {
		return nil, newCaStmp, nil
	}

	val, err = drv.Unmarshal(buff)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}
	return val, newCaStmp, nil
}

func (drv *fileListedKeyValueStore) Put(key string, val interface{}) (*Stamp, error) {
	buff, err := drv.Marshal(val)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	return drv.base.Put(key, buff)
}

func (drv *fileListedKeyValueStore) Remove(key string) error {
	return drv.base.Remove(key)
}

func (drv *fileListedKeyValueStore) Close() error {
	return drv.base.Close()
}
