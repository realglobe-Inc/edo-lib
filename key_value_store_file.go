package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

type Marshal func(interface{}) ([]byte, error)
type Unmarshal func([]byte) (interface{}, error)

type fileKeyValueStore struct {
	base RawDataStore
	Marshal
	Unmarshal
}

// スレッドセーフ。
func NewFileKeyValueStore(path string, keyGen func(string) string, marshal Marshal, unmarshal Unmarshal, expiDur time.Duration) KeyValueStore {
	return newSynchronizedKeyValueStore(newCachingKeyValueStore(newFileKeyValueStore(path, keyGen, marshal, unmarshal, expiDur)))
}

// スレッドセーフではない。
func newFileKeyValueStore(path string, keyGen func(string) string, marshal Marshal, unmarshal Unmarshal, expiDur time.Duration) *fileKeyValueStore {
	return &fileKeyValueStore{newFileRawDataStore(path, keyGen, expiDur), marshal, unmarshal}
}
func (reg *fileKeyValueStore) Get(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	buff, newCaStmp, err := reg.base.Get(key, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if buff == nil {
		return nil, newCaStmp, nil
	}

	value, err = reg.Unmarshal(buff)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}
	return value, newCaStmp, nil
}

func (reg *fileKeyValueStore) Put(key string, value interface{}) (*Stamp, error) {
	buff, err := reg.Marshal(value)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	return reg.base.Put(key, buff)
}

func (reg *fileKeyValueStore) Remove(key string) error {
	return reg.base.Remove(key)
}
