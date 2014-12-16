package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// キャッシュする。
type cachingListedRawDataStore struct {
	base  ListedRawDataStore
	cache util.Cache

	keyCache util.Cache
}

// スレッドセーフではない。
func newCachingListedRawDataStore(base ListedRawDataStore) *cachingListedRawDataStore {
	return &cachingListedRawDataStore{
		base:     base,
		cache:    util.NewCache(stampExpirationDateLess),
		keyCache: util.NewCache(stampExpirationDateLess),
	}
}

func (reg *cachingListedRawDataStore) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
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

func (reg *cachingListedRawDataStore) Get(key string, caStmp *Stamp) (data []byte, newCaStmp *Stamp, err error) {
	now := time.Now()

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: now}
	reg.cache.CleanLower(cleanThres)
	reg.keyCache.CleanLower(cleanThres)

	var buffData []byte
	var buffStmp *Stamp
	val, prio := reg.cache.Get(key)
	if prio != nil {
		// キャッシュしてた。
		buffData = val.([]byte)
		buffStmp = prio.(*Stamp)
		if now.Before(buffStmp.StaleDate) {
			// キャッシュが最新だと思って良い。
			if caStmp != nil && !caStmp.Older(buffStmp) {
				// 要求元のキャッシュより新しそうではなかった。
				return nil, buffStmp, nil
			} else {
				// 要求元のキャッシュより新しそう。
				return buffData, newCaStmp, nil
			}
		} else {
			// キャッシュが古くなっているかも。
		}
	} else {
		// キャッシュしてない。
	}

	// キャッシュしてない。
	data, newCaStmp, err = reg.base.Get(key, buffStmp)
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
	} else if data == nil {
		// キャッシュと同じ。
		reg.cache.Update(key, newCaStmp)
		buffStmp = newCaStmp
	} else {
		// あった、または、新しくなってた。
		reg.cache.Put(key, data, newCaStmp)
		// キー集合キャッシュの更新。
		if v, _ := reg.keyCache.Get(""); v != nil {
			v.(map[string]bool)[key] = true
		}
		buffData = data
		buffStmp = newCaStmp
	}

	if caStmp != nil && !caStmp.Older(buffStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, buffStmp, nil
	} else {
		// 要求元のキャッシュより新しそう。
		return buffData, buffStmp, nil
	}
}

func (reg *cachingListedRawDataStore) Put(key string, data []byte) (*Stamp, error) {
	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: time.Now()}
	reg.cache.CleanLower(cleanThres)
	reg.keyCache.CleanLower(cleanThres)

	if newCaStmp, err := reg.base.Put(key, data); err != nil {
		return nil, erro.Wrap(err)
	} else {
		// キャッシュの更新。
		reg.cache.Put(key, data, newCaStmp)
		if v, _ := reg.keyCache.Get(""); v != nil {
			v.(map[string]bool)[key] = true
		}
		return newCaStmp, nil
	}
}

func (reg *cachingListedRawDataStore) Remove(key string) error {
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
