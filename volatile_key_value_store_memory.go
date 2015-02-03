package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"strconv"
	"time"
)

func stampExpirationDateLess(a1 interface{}, a2 interface{}) bool {
	return a1.(*Stamp).ExpiDate.Before(a2.(*Stamp).ExpiDate)
}

type memoryVolatileKeyValueStore struct {
	base     util.Cache
	staleDur time.Duration
	expiDur  time.Duration
}

// スレッドセーフ。
func NewMemoryVolatileKeyValueStore(staleDur, expiDur time.Duration) VolatileKeyValueStore {
	return newSynchronizedVolatileKeyValueStore(newMemoryVolatileKeyValueStore(staleDur, expiDur))
}

// スレッドセーフではない。
func newMemoryVolatileKeyValueStore(staleDur, expiDur time.Duration) *memoryVolatileKeyValueStore {
	return &memoryVolatileKeyValueStore{util.NewCache(stampExpirationDateLess), staleDur, expiDur}
}

func (drv *memoryVolatileKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	now := time.Now()
	drv.base.CleanLower(&Stamp{ExpiDate: now})

	val, prio := drv.base.Get(key)
	if prio == nil {
		return nil, nil, nil
	}
	stmp := prio.(*Stamp)

	newCaStmp = &Stamp{
		Date:      stmp.Date,
		StaleDate: now.Add(drv.staleDur),
		ExpiDate:  now.Add(drv.expiDur),
		Digest:    stmp.Digest,
	}

	if newCaStmp.ExpiDate.After(stmp.ExpiDate) {
		newCaStmp.ExpiDate = stmp.ExpiDate
		if newCaStmp.StaleDate.After(newCaStmp.ExpiDate) {
			newCaStmp.StaleDate = newCaStmp.ExpiDate
		}
	}

	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。

	return val, newCaStmp, nil
}

func (drv *memoryVolatileKeyValueStore) Put(key string, val interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
	now := time.Now()

	stmp := &Stamp{Date: now, ExpiDate: expiDate, Digest: strconv.FormatInt(int64(now.Nanosecond()), 16)}
	drv.base.Put(key, val, stmp)

	newCaStmp = &Stamp{
		Date:      stmp.Date,
		StaleDate: now.Add(drv.staleDur),
		ExpiDate:  now.Add(drv.expiDur),
		Digest:    stmp.Digest,
	}
	if newCaStmp.ExpiDate.After(expiDate) {
		newCaStmp.ExpiDate = expiDate
		if newCaStmp.StaleDate.After(newCaStmp.ExpiDate) {
			newCaStmp.StaleDate = newCaStmp.ExpiDate
		}
	}

	drv.base.CleanLower(&Stamp{ExpiDate: now})
	return newCaStmp, nil
}

func (drv *memoryVolatileKeyValueStore) Remove(key string) error {
	drv.base.Update(key, nil)
	drv.base.CleanLower(nil)
	return nil
}
