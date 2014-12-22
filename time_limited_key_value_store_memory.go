package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"strconv"
	"time"
)

func stampExpirationDateLess(a1 interface{}, a2 interface{}) bool {
	return a1.(*Stamp).ExpiDate.Before(a2.(*Stamp).ExpiDate)
}

type memoryTimeLimitedKeyValueStore struct {
	base     util.Cache
	staleDur time.Duration
	expiDur  time.Duration
}

// スレッドセーフ。
func NewMemoryTimeLimitedKeyValueStore(staleDur, expiDur time.Duration) TimeLimitedKeyValueStore {
	return newSynchronizedTimeLimitedKeyValueStore(newMemoryTimeLimitedKeyValueStore(staleDur, expiDur))
}

// スレッドセーフではない。
func newMemoryTimeLimitedKeyValueStore(staleDur, expiDur time.Duration) *memoryTimeLimitedKeyValueStore {
	return &memoryTimeLimitedKeyValueStore{util.NewCache(stampExpirationDateLess), staleDur, expiDur}
}

func (reg *memoryTimeLimitedKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	now := time.Now()
	reg.base.CleanLower(&Stamp{ExpiDate: now})

	val, prio := reg.base.Get(key)
	if prio == nil {
		return nil, nil, nil
	}
	stmp := prio.(*Stamp)

	newCaStmp = &Stamp{
		Date:      stmp.Date,
		StaleDate: now.Add(reg.staleDur),
		ExpiDate:  now.Add(reg.expiDur),
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

func (reg *memoryTimeLimitedKeyValueStore) Put(key string, val interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
	now := time.Now()

	stmp := &Stamp{Date: now, ExpiDate: expiDate, Digest: strconv.FormatInt(int64(now.Nanosecond()), 16)}
	reg.base.Put(key, val, stmp)

	newCaStmp = &Stamp{
		Date:      stmp.Date,
		StaleDate: now.Add(reg.staleDur),
		ExpiDate:  now.Add(reg.expiDur),
		Digest:    stmp.Digest,
	}
	if newCaStmp.ExpiDate.After(expiDate) {
		newCaStmp.ExpiDate = expiDate
		if newCaStmp.StaleDate.After(newCaStmp.ExpiDate) {
			newCaStmp.StaleDate = newCaStmp.ExpiDate
		}
	}

	reg.base.CleanLower(&Stamp{ExpiDate: now})
	return newCaStmp, nil
}

func (reg *memoryTimeLimitedKeyValueStore) Remove(key string) error {
	reg.base.Update(key, nil)
	reg.base.CleanLower(nil)
	return nil
}
