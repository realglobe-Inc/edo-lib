package driver

import (
	"github.com/realglobe-Inc/edo/util/cache"
	"github.com/realglobe-Inc/go-lib/erro"
	"time"
)

// キャッシュする。
type cachingListedKeyValueStore struct {
	base   ListedKeyValueStore
	cac    cache.Cache
	keyCac cache.Cache
}

// スレッドセーフではない。
func NewCachingListedKeyValueStore(base ListedKeyValueStore) ListedKeyValueStore {
	return newCachingListedKeyValueStore(base)
}

// スレッドセーフではない。
func newCachingListedKeyValueStore(base ListedKeyValueStore) *cachingListedKeyValueStore {
	return &cachingListedKeyValueStore{
		base:   base,
		cac:    cache.New(stampExpirationDateLess),
		keyCac: cache.New(stampExpirationDateLess),
	}
}

func (drv *cachingListedKeyValueStore) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
	now := time.Now()

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: now}
	drv.cac.CleanLower(cleanThres)
	drv.keyCac.CleanLower(cleanThres)

	var buffKeys map[string]bool
	var buffStmp *Stamp
	val, prio := drv.keyCac.Get("")
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
		drv.keyCac.Update("", nil)
		return nil, nil, nil
	} else if keys == nil {
		// キャッシュと同じ。
		drv.keyCac.Update("", newCaStmp)
		buffStmp = newCaStmp
	} else {
		// あった、または、新しくなってた。
		drv.keyCac.Put("", keys, newCaStmp)
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
	drv.cac.CleanLower(cleanThres)
	drv.keyCac.CleanLower(cleanThres)

	var buffVal interface{}
	var buffStmp *Stamp
	val, prio := drv.cac.Get(key)
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
		drv.cac.Update(key, nil)
		// キー集合キャッシュの更新。
		if v, _ := drv.keyCac.Get(""); v != nil {
			delete(v.(map[string]bool), key)
		}
		return nil, nil, nil
	} else if val == nil {
		// キャッシュと同じ。
		drv.cac.Update(key, newCaStmp)
		buffStmp = newCaStmp
	} else {
		// あった、または、新しくなってた。
		drv.cac.Put(key, val, newCaStmp)
		// キー集合キャッシュの更新。
		if v, _ := drv.keyCac.Get(""); v != nil {
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
	drv.cac.CleanLower(cleanThres)
	drv.keyCac.CleanLower(cleanThres)

	if newCaStmp, err := drv.base.Put(key, val); err != nil {
		return nil, erro.Wrap(err)
	} else {
		// キャッシュの更新。
		drv.cac.Put(key, val, newCaStmp)
		if v, _ := drv.keyCac.Get(""); v != nil {
			v.(map[string]bool)[key] = true
		}
		return newCaStmp, nil
	}
}

func (drv *cachingListedKeyValueStore) Remove(key string) error {
	drv.cac.Update(key, nil)

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: time.Now()}
	drv.cac.CleanLower(cleanThres)
	drv.keyCac.CleanLower(cleanThres)

	// キャッシュの更新。
	if v, _ := drv.keyCac.Get(""); v != nil {
		delete(v.(map[string]bool), key)
	}
	return drv.base.Remove(key)
}

func (drv *cachingListedKeyValueStore) Close() error {
	if drv.base == nil {
		return nil
	} else if err := drv.base.Close(); err != nil {
		return erro.Wrap(err)
		return erro.Wrap(err)
	}
	drv.base = nil
	drv.cac = nil
	drv.keyCac = nil
	return nil
}
