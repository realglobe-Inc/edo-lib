package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// キャッシュする。
type cachingRawDataStore struct {
	base  RawDataStore
	cache util.Cache
}

// スレッドセーフではない。
func newCachingRawDataStore(base RawDataStore) *cachingRawDataStore {
	return &cachingRawDataStore{base: base, cache: util.NewCache(stampExpirationDateLess)}
}

func (reg *cachingRawDataStore) Get(key string, caStmp *Stamp) (data []byte, newCaStmp *Stamp, err error) {
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
		return value.([]byte), newCaStmp, nil
	}

	// キャッシュしてない。
	data, newCaStmp, err = reg.base.Get(key, nil)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if newCaStmp == nil {
		// 無い。
		return nil, nil, nil
	}

	// あった。
	reg.cache.Put(key, data, newCaStmp)
	if caStmp != nil && !newCaStmp.Date.After(caStmp.Date) && caStmp.Digest == newCaStmp.Digest {
		// 要求元のキャッシュと同じだった。
		return nil, newCaStmp, nil
	} else {
		return data, newCaStmp, nil
	}
}

func (reg *cachingRawDataStore) Put(key string, data []byte) (*Stamp, error) {
	if newCaStmp, err := reg.base.Put(key, data); err != nil {
		return nil, erro.Wrap(err)
	} else {
		reg.cache.Put(key, data, newCaStmp)
		return newCaStmp, nil
	}
}

func (reg *cachingRawDataStore) Remove(key string) error {
	reg.cache.Update(key, nil)
	reg.cache.CleanLower(nil)
	return reg.base.Remove(key)
}
