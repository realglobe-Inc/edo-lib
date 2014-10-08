package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

type webKeyValueStore struct {
	base RawDataStore
	Marshal
	Unmarshal
}

// スレッドセーフ。
func NewWebKeyValueStore(prefix string, marshal Marshal, unmarshal Unmarshal) KeyValueStore {
	// TODO キャッシュの並列化。
	return newSynchronizedKeyValueStore(newCachingKeyValueStore(newWebKeyValueStore(prefix, marshal, unmarshal)))
}

// スレッドセーフ。
func newWebKeyValueStore(prefix string, marshal Marshal, unmarshal Unmarshal) *webKeyValueStore {
	return &webKeyValueStore{NewWebRawDataStore(prefix), marshal, unmarshal}
}

func (reg *webKeyValueStore) Get(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
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

func (reg *webKeyValueStore) Put(key string, value interface{}) (*Stamp, error) {
	buff, err := reg.Marshal(value)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	return reg.base.Put(key, buff)
}

func (reg *webKeyValueStore) Remove(key string) error {
	return reg.base.Remove(key)
}
