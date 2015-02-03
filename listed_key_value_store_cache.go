package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// キャッシュする。
type cachingListedKeyValueStore struct {
	base     ListedKeyValueStore
	cache    util.Cache
	keyCache util.Cache
}

// スレッドセーフではない。
func NewCachingListedKeyValueStore(base ListedKeyValueStore) ListedKeyValueStore {
	return newCachingListedKeyValueStore(base)
}

// スレッドセーフではない。
func newCachingListedKeyValueStore(base ListedKeyValueStore) *cachingListedKeyValueStore {
	return &cachingListedKeyValueStore{
		base:     base,
		cache:    util.NewCache(stampExpirationDateLess),
		keyCache: util.NewCache(stampExpirationDateLess),
	}
}

func (drv *cachingListedKeyValueStore) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
	now := time.Now()

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: now}
	drv.cache.CleanLower(cleanThres)
	drv.keyCache.CleanLower(cleanThres)

	var buffKeys map[string]bool
	var buffStmp *Stamp
	val, prio := drv.keyCache.Get("")
	if prio != nil {
		// キャッシュしてた。
		buffKeys = val.(map[string]bool)
		buffStmp = prio.(*Stamp)
		if now.Before(buffStmp.StaleDate) {
			// キャッシュが最新だと思って良い。
			if caStmp != nil && !caStmp.Older(buffStmp) {
				// 要求元のキャッシュより新しそうではなかった。
				return nil, buffStmp, nil
			}
			// 要求元のキャッシュより新しそう。
			keys := map[string]bool{}
			for k, v := range buffKeys {
				keys[k] = v
			}
			return keys, buffStmp, nil
		} else {
			// キャッシュが古くなっているかも。
		}
	} else {
		// キャッシュしてない。
	}

	keys, newCaStmp, err = drv.base.Keys(buffStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if newCaStmp == nil {
		// 1 つも無い。
		drv.keyCache.Update("", nil)
		return nil, nil, nil
	} else if keys == nil {
		// キャッシュと同じ。
		drv.keyCache.Update("", newCaStmp)
		buffStmp = newCaStmp
	} else {
		// あった、または、新しくなってた。
		drv.keyCache.Put("", keys, newCaStmp)
		buffKeys = keys
		buffStmp = newCaStmp
	}

	if caStmp != nil && !caStmp.Older(buffStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, buffStmp, nil
	} else {
		// 要求元のキャッシュより新しそう。
		keys := map[string]bool{}
		for k, v := range buffKeys {
			keys[k] = v
		}
		return keys, buffStmp, nil
	}
}

func (drv *cachingListedKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	now := time.Now()

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: now}
	drv.cache.CleanLower(cleanThres)
	drv.keyCache.CleanLower(cleanThres)

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
		// キー集合キャッシュの更新。
		if v, _ := drv.keyCache.Get(""); v != nil {
			delete(v.(map[string]bool), key)
		}
		return nil, nil, nil
	} else if val == nil {
		// キャッシュと同じ。
		drv.cache.Update(key, newCaStmp)
		buffStmp = newCaStmp
	} else {
		// あった、または、新しくなってた。
		drv.cache.Put(key, val, newCaStmp)
		// キー集合キャッシュの更新。
		if v, _ := drv.keyCache.Get(""); v != nil {
			v.(map[string]bool)[key] = true
		}
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

func (drv *cachingListedKeyValueStore) Put(key string, val interface{}) (*Stamp, error) {
	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: time.Now()}
	drv.cache.CleanLower(cleanThres)
	drv.keyCache.CleanLower(cleanThres)

	if newCaStmp, err := drv.base.Put(key, val); err != nil {
		return nil, erro.Wrap(err)
	} else {
		// キャッシュの更新。
		drv.cache.Put(key, val, newCaStmp)
		if v, _ := drv.keyCache.Get(""); v != nil {
			v.(map[string]bool)[key] = true
		}
		return newCaStmp, nil
	}
}

func (drv *cachingListedKeyValueStore) Remove(key string) error {
	drv.cache.Update(key, nil)

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: time.Now()}
	drv.cache.CleanLower(cleanThres)
	drv.keyCache.CleanLower(cleanThres)

	// キャッシュの更新。
	if v, _ := drv.keyCache.Get(""); v != nil {
		delete(v.(map[string]bool), key)
	}
	return drv.base.Remove(key)
}
