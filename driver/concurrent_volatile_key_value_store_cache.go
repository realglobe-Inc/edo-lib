// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package driver

import (
	"github.com/realglobe-Inc/edo-lib/cache"
	"github.com/realglobe-Inc/go-lib/erro"
	"time"
)

type cachingConcurrentVolatileKeyValueStore struct {
	base ConcurrentVolatileKeyValueStore
	cac  cache.Cache
}

// スレッドセーフではない。
func newCachingConcurrentVolatileKeyValueStore(base ConcurrentVolatileKeyValueStore) *cachingConcurrentVolatileKeyValueStore {
	return &cachingConcurrentVolatileKeyValueStore{
		base: base,
		cac:  cache.New(stampExpirationDateLess),
	}
}

func (drv *cachingConcurrentVolatileKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	now := time.Now()

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: now}
	drv.cac.CleanLower(cleanThres)

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
		return nil, nil, nil
	} else if val == nil {
		// キャッシュと同じ。
		drv.cac.Update(key, newCaStmp)
		buffStmp = newCaStmp
	} else {
		// あった、または、新しくなってた。
		drv.cac.Put(key, val, newCaStmp)
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
	drv.cac.CleanLower(cleanThres)

	if newCaStmp, err := drv.base.Put(key, val, expiDate); err != nil {
		return nil, erro.Wrap(err)
	} else {
		// キャッシュの更新。
		drv.cac.Put(key, val, newCaStmp)
		return newCaStmp, nil
	}
}

func (drv *cachingConcurrentVolatileKeyValueStore) Remove(key string) error {
	drv.cac.Update(key, nil)

	// 古いキャッシュの削除。
	cleanThres := &Stamp{ExpiDate: time.Now()}
	drv.cac.CleanLower(cleanThres)

	return drv.base.Remove(key)
}

func (drv *cachingConcurrentVolatileKeyValueStore) Close() error {
	if drv.base == nil {
		return nil
	} else if err := drv.base.Close(); err != nil {
		return erro.Wrap(err)
	}
	drv.base = nil
	drv.cac = nil
	return nil
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
