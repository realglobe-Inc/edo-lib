package driver

import (
	"strconv"
	"time"
)

type memoryListedKeyValueStore struct {
	date      time.Time
	digest    int
	keyToVal  map[string]interface{}
	keyToStmp map[string]*Stamp
	staleDur  time.Duration
	expiDur   time.Duration
}

// スレッドセーフ。
func NewMemoryListedKeyValueStore(staleDur, expiDur time.Duration) ListedKeyValueStore {
	return newSynchronizedListedKeyValueStore(newMemoryListedKeyValueStore(staleDur, expiDur))
}

// スレッドセーフではない。
func newMemoryListedKeyValueStore(staleDur, expiDur time.Duration) *memoryListedKeyValueStore {
	return &memoryListedKeyValueStore{
		date:      time.Now(),
		digest:    0,
		keyToVal:  map[string]interface{}{},
		keyToStmp: map[string]*Stamp{},
		staleDur:  staleDur,
		expiDur:   expiDur,
	}
}

func (drv *memoryListedKeyValueStore) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
	newCaStmp = &Stamp{Date: drv.date, Digest: strconv.FormatInt(int64(drv.digest), 16)}
	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。

	keys = map[string]bool{}
	for key, _ := range drv.keyToVal {
		keys[key] = true
	}
	return keys, newCaStmp, nil
}

func (drv *memoryListedKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	stmp := drv.keyToStmp[key]
	if stmp == nil {
		return nil, nil, nil
	}
	now := time.Now()
	newCaStmp = &Stamp{
		Date:      stmp.Date,
		StaleDate: now.Add(drv.staleDur),
		ExpiDate:  now.Add(drv.expiDur),
		Digest:    stmp.Digest,
	}

	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。

	return drv.keyToVal[key], newCaStmp, nil
}

func (drv *memoryListedKeyValueStore) Put(key string, val interface{}) (newCaStmp *Stamp, err error) {
	now := time.Now()
	stmp := &Stamp{Date: now, Digest: strconv.FormatInt(int64(now.Nanosecond()), 16)}
	drv.keyToVal[key] = val
	drv.keyToStmp[key] = stmp
	s := *stmp
	s.StaleDate = now.Add(drv.staleDur)
	s.ExpiDate = now.Add(drv.expiDur)
	drv.date = now
	drv.digest++
	return &s, nil
}

func (drv *memoryListedKeyValueStore) Remove(key string) error {
	if _, ok := drv.keyToVal[key]; !ok {
		return nil
	}

	delete(drv.keyToVal, key)
	delete(drv.keyToStmp, key)
	drv.date = time.Now()
	drv.digest++
	return nil
}
