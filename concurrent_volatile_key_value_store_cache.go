package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

type cachingConcurrentVolatileKeyValueStore struct {
	base  ConcurrentVolatileKeyValueStore
	cache util.Cache
}

// スレッドセーフではない。
func newCachingConcurrentVolatileKeyValueStore(base ConcurrentVolatileKeyValueStore) *cachingConcurrentVolatileKeyValueStore {
	return &cachingConcurrentVolatileKeyValueStore{
		base:  base,
		cache: util.NewCache(stampExpirationDateLess),
	}
}

func (drv *cachingConcurrentVolatileKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	now := time.Now()

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: now}
	drv.cache.CleanLower(cleanThres)

	var buffVal interface{}
	var buffStmp *Stamp
	val, prio := drv.cache.Get(key)
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
	val, newCaStmp, err = drv.base.Get(key, buffStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if newCaStmp == nil {
		// 無い。
		drv.cache.Update(key, nil)
		return nil, nil, nil
	} else if val == nil {
		// キャッシュと同じ。
		drv.cache.Update(key, newCaStmp)
		buffStmp = newCaStmp
	} else {
		// あった、または、新しくなってた。
		drv.cache.Put(key, val, newCaStmp)
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

func (drv *cachingConcurrentVolatileKeyValueStore) Put(key string, val interface{}, expiDate time.Time) (*Stamp, error) {
	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: time.Now()}
	drv.cache.CleanLower(cleanThres)

	if newCaStmp, err := drv.base.Put(key, val, expiDate); err != nil {
		return nil, erro.Wrap(err)
	} else {
		// キャッシュの更新。
		drv.cache.Put(key, val, newCaStmp)
		return newCaStmp, nil
	}
}

func (drv *cachingConcurrentVolatileKeyValueStore) Remove(key string) error {
	drv.cache.Update(key, nil)

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: time.Now()}
	drv.cache.CleanLower(cleanThres)

	return drv.base.Remove(key)
}

func (drv *cachingConcurrentVolatileKeyValueStore) Entry(eKey string) (eVal string, err error) {
	return drv.base.Entry(eKey)
}

func (drv *cachingConcurrentVolatileKeyValueStore) SetEntry(eKey, eVal string, eExpiDate time.Time) error {
	return drv.SetEntry(eKey, eVal, eExpiDate)
}

func (drv *cachingConcurrentVolatileKeyValueStore) GetAndSetEntry(key string, caStmp *Stamp, eKey, eVal string, eExpiDate time.Time) (val interface{}, newCaStmp *Stamp, err error) {
	val, newCaStmp, err = drv.Get(key, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	if err := drv.base.SetEntry(eKey, eVal, eExpiDate); err != nil {
		return nil, nil, erro.Wrap(err)
	}

	return val, newCaStmp, nil
}

func (drv *cachingConcurrentVolatileKeyValueStore) PutIfEntered(key string, val interface{}, expiDate time.Time, eKey, eVal string) (entered bool, newCaStmp *Stamp, err error) {
	eV, err := drv.base.Entry(eKey)
	if err != nil {
		return false, nil, erro.Wrap(err)
	} else if eVal != eV {
		return false, nil, nil
	}

	newCaStmp, err = drv.Put(key, val, expiDate)
	if err != nil {
		return false, nil, erro.Wrap(err)
	}
	return true, newCaStmp, nil
}
