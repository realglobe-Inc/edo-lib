package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

type cachingVolatileKeyValueStore struct {
	base  VolatileKeyValueStore
	cache util.Cache
}

// スレッドセーフではない。
func newCachingVolatileKeyValueStore(base VolatileKeyValueStore) *cachingVolatileKeyValueStore {
	return &cachingVolatileKeyValueStore{
		base:  base,
		cache: util.NewCache(stampExpirationDateLess),
	}
}

func (reg *cachingVolatileKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	now := time.Now()

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: now}
	reg.cache.CleanLower(cleanThres)

	var buffVal interface{}
	var buffStmp *Stamp
	val, prio := reg.cache.Get(key)
	if prio != nil {
		// キャッシュしてた。
		buffVal = val.(interface{})
		buffStmp = prio.(*Stamp)
		if now.Before(buffStmp.StaleDate) {
			// キャッシュが最新だと思って良い。
			if caStmp != nil && !caStmp.Older(buffStmp) {
				// 要求元のキャッシュより新しそうではなかった。
				return nil, buffStmp, nil
			} else {
				// 要求元のキャッシュより新しそう。
				return buffVal, newCaStmp, nil
			}
		} else {
			// キャッシュが古くなっているかも。
		}
	} else {
		// キャッシュしてない。
	}

	// キャッシュしてない。
	val, newCaStmp, err = reg.base.Get(key, buffStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if newCaStmp == nil {
		// 無い。
		reg.cache.Update(key, nil)
		return nil, nil, nil
	} else if val == nil {
		// キャッシュと同じ。
		reg.cache.Update(key, newCaStmp)
		buffStmp = newCaStmp
	} else {
		// あった、または、新しくなってた。
		reg.cache.Put(key, val, newCaStmp)
		buffVal = val
		buffStmp = newCaStmp
	}

	if caStmp != nil && !caStmp.Older(buffStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, buffStmp, nil
	} else {
		// 要求元のキャッシュより新しそう。
		return buffVal, buffStmp, nil
	}
}

func (reg *cachingVolatileKeyValueStore) Put(key string, val interface{}, expiDate time.Time) (*Stamp, error) {
	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: time.Now()}
	reg.cache.CleanLower(cleanThres)

	if newCaStmp, err := reg.base.Put(key, val, expiDate); err != nil {
		return nil, erro.Wrap(err)
	} else {
		// キャッシュの更新。
		reg.cache.Put(key, val, newCaStmp)
		return newCaStmp, nil
	}
}

func (reg *cachingVolatileKeyValueStore) Remove(key string) error {
	reg.cache.Update(key, nil)

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: time.Now()}
	reg.cache.CleanLower(cleanThres)

	return reg.base.Remove(key)
}
