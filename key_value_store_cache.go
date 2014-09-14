package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// キャッシュする。

// キャッシュ用。
type cachingDatedKeyValueStore struct {
	datedKeyValueStore
	cache util.Cache
}

func newCachingDatedKeyValueStore(backend datedKeyValueStore) *cachingDatedKeyValueStore {
	return &cachingDatedKeyValueStore{datedKeyValueStore: backend,
		cache: util.NewCache(func(a1 interface{}, a2 interface{}) bool {
			return a1.(*Stamp).ExpiDate.Before(a2.(*Stamp).ExpiDate)
		}),
	}
}

func (reg *cachingDatedKeyValueStore) stampedGet(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	now := time.Now()
	reg.cache.CleanLesser(&Stamp{ExpiDate: now})

	// 残ってるキャッシュは有効。

	val, prio := reg.cache.Get(key)
	if prio == nil {
		// キャッシュしてない。
		value, newCaStmp, err = reg.datedKeyValueStore.stampedGet(key, nil)
		if err != nil {
			return nil, nil, erro.Wrap(err)
		} else if newCaStmp == nil {
			// 無い。
			return nil, nil, nil
		} else {
			// あった。
			reg.cache.Put(key, value, newCaStmp)
			if caStmp != nil && !newCaStmp.Date.After(caStmp.Date) && caStmp.Digest == newCaStmp.Digest {
				// 要求元のキャッシュと同じだった。
				return nil, newCaStmp, nil
			} else {
				return value, newCaStmp, nil
			}
		}
	}

	// キャッシュしてた。

	stmp := prio.(*Stamp)
	if caStmp != nil && !stmp.Date.After(caStmp.Date) && caStmp.Digest == stmp.Digest {
		return nil, stmp, nil
	}
	return val, stmp, nil
}

func (reg *cachingDatedKeyValueStore) stampedPut(key string, value interface{}) (*Stamp, error) {
	if newCaStmp, err := reg.datedKeyValueStore.stampedPut(key, value); err != nil {
		return nil, erro.Wrap(err)
	} else {
		reg.cache.Put(key, value, newCaStmp)
		return newCaStmp, nil
	}
}

func (reg *cachingDatedKeyValueStore) remove(key string) error {
	return reg.datedKeyValueStore.remove(key)
}
