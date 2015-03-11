package driver

import (
	"github.com/realglobe-Inc/edo-toolkit/util/cache"
	"github.com/realglobe-Inc/go-lib/erro"
	"time"
)

// キャッシュする。
type cachingListedRawDataStore struct {
	base ListedRawDataStore
	cac  cache.Cache

	keyCac cache.Cache
}

// スレッドセーフではない。
func newCachingListedRawDataStore(base ListedRawDataStore) *cachingListedRawDataStore {
	return &cachingListedRawDataStore{
		base:   base,
		cac:    cache.New(stampExpirationDateLess),
		keyCac: cache.New(stampExpirationDateLess),
	}
}

func (drv *cachingListedRawDataStore) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
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

func (drv *cachingListedRawDataStore) Get(key string, caStmp *Stamp) (data []byte, newCaStmp *Stamp, err error) {
	now := time.Now()

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: now}
	drv.cac.CleanLower(cleanThres)
	drv.keyCac.CleanLower(cleanThres)

	var buffData []byte
	var buffStmp *Stamp
	val, prio := drv.cac.Get(key)
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
	data, newCaStmp, err = drv.base.Get(key, buffStmp)
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
	} else if data == nil {
		// キャッシュと同じ。
		drv.cac.Update(key, newCaStmp)
		buffStmp = newCaStmp
	} else {
		// あった、または、新しくなってた。
		drv.cac.Put(key, data, newCaStmp)
		// キー集合キャッシュの更新。
		if v, _ := drv.keyCac.Get(""); v != nil {
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

func (drv *cachingListedRawDataStore) Put(key string, data []byte) (*Stamp, error) {
	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: time.Now()}
	drv.cac.CleanLower(cleanThres)
	drv.keyCac.CleanLower(cleanThres)

	if newCaStmp, err := drv.base.Put(key, data); err != nil {
		return nil, erro.Wrap(err)
	} else {
		// キャッシュの更新。
		drv.cac.Put(key, data, newCaStmp)
		if v, _ := drv.keyCac.Get(""); v != nil {
			v.(map[string]bool)[key] = true
		}
		return newCaStmp, nil
	}
}

func (drv *cachingListedRawDataStore) Remove(key string) error {
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

func (drv *cachingListedRawDataStore) Close() error {
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
