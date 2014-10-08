package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// キャッシュする。
type cachingKeyValueStore struct {
	base  KeyValueStore
	cache util.Cache
}

// スレッドセーフではない。
func newCachingKeyValueStore(base KeyValueStore) *cachingKeyValueStore {
	return &cachingKeyValueStore{base: base, cache: util.NewCache(stampExpirationDateLess)}
}

func (reg *cachingKeyValueStore) Get(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	now := time.Now()
	reg.cache.CleanLower(&Stamp{ExpiDate: now})

	// 残ってるキャッシュは有効。

	value, prio := reg.cache.Get(key)
	if prio != nil {
		// キャッシュしてた。
		newCaStmp = prio.(*Stamp)
		if caStmp != nil && !newCaStmp.Date.After(caStmp.Date) && caStmp.Digest == newCaStmp.Digest {
			return nil, newCaStmp, nil
		}
		return value, newCaStmp, nil
	}

	// キャッシュしてない。
	value, newCaStmp, err = reg.base.Get(key, nil)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if newCaStmp == nil {
		// 無い。
		return nil, nil, nil
	}

	// あった。
	reg.cache.Put(key, value, newCaStmp)
	if caStmp != nil && !newCaStmp.Date.After(caStmp.Date) && caStmp.Digest == newCaStmp.Digest {
		// 要求元のキャッシュと同じだった。
		return nil, newCaStmp, nil
	} else {
		return value, newCaStmp, nil
	}
}

func (reg *cachingKeyValueStore) Put(key string, value interface{}) (*Stamp, error) {
	if newCaStmp, err := reg.base.Put(key, value); err != nil {
		return nil, erro.Wrap(err)
	} else {
		reg.cache.Put(key, value, newCaStmp)
		return newCaStmp, nil
	}
}

func (reg *cachingKeyValueStore) Remove(key string) error {
	reg.cache.Update(key, nil)
	reg.cache.CleanLower(nil)
	return reg.base.Remove(key)
}
