package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// キャッシュする。
type cachingKeyValueStore struct {
	base     KeyValueStore
	cache    util.Cache
	keyCache util.Cache
}

// スレッドセーフではない。
func NewCachingKeyValueStore(base KeyValueStore) KeyValueStore {
	return newCachingKeyValueStore(base)
}

// スレッドセーフではない。
func newCachingKeyValueStore(base KeyValueStore) *cachingKeyValueStore {
	return &cachingKeyValueStore{
		base:     base,
		cache:    util.NewCache(stampExpirationDateLess),
		keyCache: util.NewCache(stampExpirationDateLess),
	}
}

func (reg *cachingKeyValueStore) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
	now := time.Now()

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: now}
	reg.cache.CleanLower(cleanThres)
	reg.keyCache.CleanLower(cleanThres)

	var buffKeys map[string]bool
	var buffStmp *Stamp
	val, prio := reg.keyCache.Get("")
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

	keys, newCaStmp, err = reg.base.Keys(buffStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if newCaStmp == nil {
		// 1 つも無い。
		reg.keyCache.Update("", nil)
		return nil, nil, nil
	} else if keys == nil {
		// キャッシュと同じ。
		reg.keyCache.Update("", newCaStmp)
		buffStmp = newCaStmp
	} else {
		// あった、または、新しくなってた。
		reg.keyCache.Put("", keys, newCaStmp)
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

func (reg *cachingKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	now := time.Now()

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: now}
	reg.cache.CleanLower(cleanThres)
	reg.keyCache.CleanLower(cleanThres)

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
		// キー集合キャッシュの更新。
		if v, _ := reg.keyCache.Get(""); v != nil {
			delete(v.(map[string]bool), key)
		}
		return nil, nil, nil
	} else if val == nil {
		// キャッシュと同じ。
		reg.cache.Update(key, newCaStmp)
		buffStmp = newCaStmp
	} else {
		// あった、または、新しくなってた。
		reg.cache.Put(key, val, newCaStmp)
		// キー集合キャッシュの更新。
		if v, _ := reg.keyCache.Get(""); v != nil {
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

func (reg *cachingKeyValueStore) Put(key string, val interface{}) (*Stamp, error) {
	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: time.Now()}
	reg.cache.CleanLower(cleanThres)
	reg.keyCache.CleanLower(cleanThres)

	if newCaStmp, err := reg.base.Put(key, val); err != nil {
		return nil, erro.Wrap(err)
	} else {
		// キャッシュの更新。
		reg.cache.Put(key, val, newCaStmp)
		if v, _ := reg.keyCache.Get(""); v != nil {
			v.(map[string]bool)[key] = true
		}
		return newCaStmp, nil
	}
}

func (reg *cachingKeyValueStore) Remove(key string) error {
	reg.cache.Update(key, nil)

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: time.Now()}
	reg.cache.CleanLower(cleanThres)
	reg.keyCache.CleanLower(cleanThres)

	// キャッシュの更新。
	if v, _ := reg.keyCache.Get(""); v != nil {
		delete(v.(map[string]bool), key)
	}
	return reg.base.Remove(key)
}
