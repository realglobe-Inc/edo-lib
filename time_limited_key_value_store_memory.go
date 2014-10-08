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
	base    util.Cache
	expiDur time.Duration
}

// スレッドセーフ。
func NewMemoryTimeLimitedKeyValueStore(expiDur time.Duration) TimeLimitedKeyValueStore {
	return newSynchronizedTimeLimitedKeyValueStore(newMemoryTimeLimitedKeyValueStore(expiDur))
}

// スレッドセーフではない。
func newMemoryTimeLimitedKeyValueStore(expiDur time.Duration) *memoryTimeLimitedKeyValueStore {
	return &memoryTimeLimitedKeyValueStore{util.NewCache(stampExpirationDateLess), expiDur}
}

func (reg *memoryTimeLimitedKeyValueStore) Get(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	now := time.Now()
	reg.base.CleanLower(&Stamp{ExpiDate: now})

	value, prio := reg.base.Get(key)
	if prio == nil {
		return nil, nil, nil
	}
	stmp := prio.(*Stamp)

	newCaStmp = &Stamp{Date: stmp.Date, ExpiDate: now.Add(reg.expiDur), Digest: stmp.Digest}
	if newCaStmp.ExpiDate.After(stmp.ExpiDate) {
		// newCaStmp.ExpiDate は stmp.ExpiDate 以前。
		newCaStmp.ExpiDate = stmp.ExpiDate
	}

	if caStmp == nil || caStmp.Date.Before(newCaStmp.Date) || caStmp.Digest != newCaStmp.Digest {
		return value, newCaStmp, nil
	}
	return nil, newCaStmp, nil
}

func (reg *memoryTimeLimitedKeyValueStore) Put(key string, value interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
	now := time.Now()

	stmp := &Stamp{Date: now, ExpiDate: expiDate, Digest: strconv.FormatInt(now.UnixNano(), 10)}
	reg.base.Put(key, value, stmp)

	newCaStmp = &Stamp{Date: stmp.Date, ExpiDate: now.Add(reg.expiDur), Digest: stmp.Digest}
	if newCaStmp.ExpiDate.After(expiDate) {
		// newCaStmp.ExpiDate は expiDate 以前。
		newCaStmp.ExpiDate = expiDate
	}

	reg.base.CleanLower(&Stamp{ExpiDate: now})
	return newCaStmp, nil
}

func (reg *memoryTimeLimitedKeyValueStore) Remove(key string) error {
	reg.base.Update(key, nil)
	reg.base.CleanLower(nil)
	return nil
}
