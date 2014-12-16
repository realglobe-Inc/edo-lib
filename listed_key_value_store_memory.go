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

func (reg *memoryListedKeyValueStore) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
	newCaStmp = &Stamp{Date: reg.date, Digest: strconv.FormatInt(int64(reg.digest), 16)}
	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。

	keys = map[string]bool{}
	for key, _ := range reg.keyToVal {
		keys[key] = true
	}
	return keys, newCaStmp, nil
}

func (reg *memoryListedKeyValueStore) Get(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	stmp := reg.keyToStmp[key]
	if stmp == nil {
		return nil, nil, nil
	}
	now := time.Now()
	newCaStmp = &Stamp{
		Date:      stmp.Date,
		StaleDate: now.Add(reg.staleDur),
		ExpiDate:  now.Add(reg.expiDur),
		Digest:    stmp.Digest,
	}

	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。

	return reg.keyToVal[key], newCaStmp, nil
}

func (reg *memoryListedKeyValueStore) Put(key string, val interface{}) (newCaStmp *Stamp, err error) {
	now := time.Now()
	stmp := &Stamp{Date: now, Digest: strconv.FormatInt(int64(now.Nanosecond()), 16)}
	reg.keyToVal[key] = val
	reg.keyToStmp[key] = stmp
	s := *stmp
	s.StaleDate = now.Add(reg.staleDur)
	s.ExpiDate = now.Add(reg.expiDur)
	reg.date = now
	reg.digest++
	return &s, nil
}

func (reg *memoryListedKeyValueStore) Remove(key string) error {
	if _, ok := reg.keyToVal[key]; !ok {
		return nil
	}

	delete(reg.keyToVal, key)
	delete(reg.keyToStmp, key)
	reg.date = time.Now()
	reg.digest++
	return nil
}
