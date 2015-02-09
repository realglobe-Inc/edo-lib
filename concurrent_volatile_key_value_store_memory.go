package driver

import (
	"github.com/realglobe-Inc/edo/util/cache"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"strconv"
	"time"
)

func stampExpirationDateLess(a1 interface{}, a2 interface{}) bool {
	return a1.(*Stamp).ExpiDate.Before(a2.(*Stamp).ExpiDate)
}

func dateLess(a1 interface{}, a2 interface{}) bool {
	return a1.(time.Time).Before(a2.(time.Time))
}

type memoryConcurrentVolatileKeyValueStore struct {
	base     cache.Cache
	staleDur time.Duration
	expiDur  time.Duration

	ents cache.Cache
}

// スレッドセーフ。
func NewMemoryConcurrentVolatileKeyValueStore(staleDur, expiDur time.Duration) ConcurrentVolatileKeyValueStore {
	return newSynchronizedConcurrentVolatileKeyValueStore(newMemoryConcurrentVolatileKeyValueStore(staleDur, expiDur))
}

// スレッドセーフではない。
func newMemoryConcurrentVolatileKeyValueStore(staleDur, expiDur time.Duration) *memoryConcurrentVolatileKeyValueStore {
	return &memoryConcurrentVolatileKeyValueStore{
		cache.New(stampExpirationDateLess),
		staleDur,
		expiDur,
		cache.New(dateLess),
	}
}

func (drv *memoryConcurrentVolatileKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
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

func (drv *memoryConcurrentVolatileKeyValueStore) Put(key string, val interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
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

func (drv *memoryConcurrentVolatileKeyValueStore) Remove(key string) error {
	drv.base.Update(key, nil)
	drv.base.CleanLower(nil)
	return nil
}

func (drv *memoryConcurrentVolatileKeyValueStore) Entry(eKey string) (eVal string, err error) {
	drv.ents.CleanLower(time.Now())
	eV, _ := drv.ents.Get(eKey)
	eVal, _ = eV.(string)
	return eVal, nil
}

func (drv *memoryConcurrentVolatileKeyValueStore) SetEntry(eKey, eVal string, eExpiDate time.Time) error {
	drv.ents.CleanLower(time.Now())
	drv.ents.Put(eKey, eVal, eExpiDate)
	return nil
}

func (drv *memoryConcurrentVolatileKeyValueStore) GetAndSetEntry(key string, caStmp *Stamp, eKey, eVal string, eExpiDate time.Time) (val interface{}, newCaStmp *Stamp, err error) {
	val, newCaStmp, err = drv.Get(key, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	if err := drv.SetEntry(eKey, eVal, eExpiDate); err != nil {
		return nil, nil, erro.Wrap(err)
	}

	return val, newCaStmp, nil
}

func (drv *memoryConcurrentVolatileKeyValueStore) PutIfEntered(key string, val interface{}, expiDate time.Time, eKey, eVal string) (entered bool, newCaStmp *Stamp, err error) {
	eV, err := drv.Entry(eKey)
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
