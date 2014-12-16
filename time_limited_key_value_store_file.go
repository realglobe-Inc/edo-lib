package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// TODO 今は手抜きで古いファイルを無視するだけ。どんどん溜まっていく。

type fileTimeLimitedKeyValueStore struct {
	base    KeyValueStore
	expires KeyValueStore
}

// スレッドセーフ。
func NewFileTimeLimitedKeyValueStore(path, expiPath string, keyToPath, pathToKey func(string) string, marshal Marshal, unmarshal Unmarshal, staleDur, expiDur time.Duration) TimeLimitedKeyValueStore {
	return newSynchronizedTimeLimitedKeyValueStore(newCachingTimeLimitedKeyValueStore(newFileTimeLimitedKeyValueStore(path, expiPath, keyToPath, pathToKey, marshal, unmarshal, staleDur, expiDur)))
}

// スレッドセーフではない。
func newFileTimeLimitedKeyValueStore(path, expiPath string, keyToPath, pathToKey func(string) string, marshal Marshal, unmarshal Unmarshal, staleDur, expiDur time.Duration) *fileTimeLimitedKeyValueStore {
	return &fileTimeLimitedKeyValueStore{
		NewFileListedKeyValueStore(path, keyToPath, pathToKey, marshal, unmarshal, staleDur, expiDur),
		NewFileListedKeyValueStore(expiPath, keyToPath, pathToKey,
			func(value interface{}) ([]byte, error) {
				return []byte(value.(time.Time).Format(time.RFC3339Nano)), nil
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

func (reg *fileTimeLimitedKeyValueStore) Get(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	var expiDate time.Time
	if value, newCaStmp, err := reg.expires.Get(key, nil); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if newCaStmp == nil {
		return nil, nil, nil
	} else {
		expiDate = value.(time.Time)
	}

	if time.Now().After(expiDate) {
		reg.expires.Remove(key)
		reg.base.Remove(key)
		return nil, nil, nil
	}

	value, newCaStmp, err = reg.base.Get(key, caStmp)
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
	return value, newCaStmp, nil
}

func (reg *fileTimeLimitedKeyValueStore) Put(key string, value interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
	if _, err := reg.expires.Put(key, expiDate); err != nil {
		return nil, erro.Wrap(err)
	}

	newCaStmp, err = reg.base.Put(key, value)
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

func (reg *fileTimeLimitedKeyValueStore) Remove(key string) error {
	if err := reg.expires.Remove(key); err != nil {
		return erro.Wrap(err)
	}
	return reg.base.Remove(key)
}
