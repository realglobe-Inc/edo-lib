package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// TODO 今は手抜きで古いファイルを無視するだけ。どんどん溜まっていく。

type fileVolatileKeyValueStore struct {
	base KeyValueStore
	exps KeyValueStore
}

// スレッドセーフ。
func NewFileVolatileKeyValueStore(path, expiPath string, keyToPath, pathToKey func(string) string, marshal Marshal, unmarshal Unmarshal, staleDur, expiDur time.Duration) VolatileKeyValueStore {
	return newSynchronizedVolatileKeyValueStore(newCachingVolatileKeyValueStore(newFileVolatileKeyValueStore(path, expiPath, keyToPath, pathToKey, marshal, unmarshal, staleDur, expiDur)))
}

// スレッドセーフではない。
func newFileVolatileKeyValueStore(path, expiPath string, keyToPath, pathToKey func(string) string, marshal Marshal, unmarshal Unmarshal, staleDur, expiDur time.Duration) *fileVolatileKeyValueStore {
	return &fileVolatileKeyValueStore{
		NewFileListedKeyValueStore(path, keyToPath, pathToKey, marshal, unmarshal, staleDur, expiDur),
		NewFileListedKeyValueStore(expiPath, keyToPath, pathToKey,
			func(val interface{}) ([]byte, error) {
				return []byte(val.(time.Time).Format(time.RFC3339Nano)), nil
			},
			func(data []byte) (interface{}, error) {
				date, err := time.Parse(time.RFC3339Nano, string(data))
				if err != nil {
					return time.Time{}, erro.Wrap(err)
				}
				return date, nil
			},
			staleDur, expiDur),
	}
}

func (drv *fileVolatileKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	var expiDate time.Time
	if val, newCaStmp, err := drv.exps.Get(key, nil); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if newCaStmp == nil {
		return nil, nil, nil
	} else {
		expiDate = val.(time.Time)
	}

	if time.Now().After(expiDate) {
		drv.exps.Remove(key)
		drv.base.Remove(key)
		return nil, nil, nil
	}

	val, newCaStmp, err = drv.base.Get(key, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if newCaStmp == nil {
		return nil, nil, nil
	}

	if newCaStmp.ExpiDate.After(expiDate) {
		newCaStmp.ExpiDate = expiDate
		if newCaStmp.StaleDate.After(newCaStmp.ExpiDate) {
			newCaStmp.StaleDate = newCaStmp.ExpiDate
		}
	}
	return val, newCaStmp, nil
}

func (drv *fileVolatileKeyValueStore) Put(key string, val interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
	if _, err := drv.exps.Put(key, expiDate); err != nil {
		return nil, erro.Wrap(err)
	}

	newCaStmp, err = drv.base.Put(key, val)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	if newCaStmp.ExpiDate.After(expiDate) {
		newCaStmp.ExpiDate = expiDate
		if newCaStmp.StaleDate.After(newCaStmp.ExpiDate) {
			newCaStmp.StaleDate = newCaStmp.ExpiDate
		}
	}
	return newCaStmp, nil
}

func (drv *fileVolatileKeyValueStore) Remove(key string) error {
	if err := drv.exps.Remove(key); err != nil {
		return erro.Wrap(err)
	}
	return drv.base.Remove(key)
}
