package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"time"
)

type memoryTimeLimitedKeyValueStore struct {
	util.Cache
}

func NewMemoryTimeLimitedKeyValueStore() TimeLimitedKeyValueStore {
	return &memoryTimeLimitedKeyValueStore{util.NewCache(func(a1 interface{}, a2 interface{}) bool {
		return a1.(time.Time).Before(a2.(time.Time))
	})}
}

func (this *memoryTimeLimitedKeyValueStore) Get(key string) (value interface{}, err error) {
	this.Cache.CleanLesser(time.Now())

	val, _ := this.Cache.Get(key)
	return val, nil
}

func (this *memoryTimeLimitedKeyValueStore) Put(key string, value interface{}, timLim time.Time) error {
	this.Cache.Put(key, value, timLim)

	this.Cache.CleanLesser(time.Now())
	return nil
}

func (this *memoryTimeLimitedKeyValueStore) Remove(key string) error {
	// 期限切れにして押し出す。
	now := time.Now()
	this.Cache.Update(key, now.Add(-time.Second))

	this.Cache.CleanLesser(now)
	return nil
}
