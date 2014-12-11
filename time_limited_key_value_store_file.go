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
func NewFileTimeLimitedKeyValueStore(path string, keyToPath, pathToKey func(string) string, marshal Marshal, unmarshal Unmarshal, staleDur, expiDur time.Duration) TimeLimitedKeyValueStore {
	return newSynchronizedTimeLimitedKeyValueStore(newCachingTimeLimitedKeyValueStore(newFileTimeLimitedKeyValueStore(path, keyToPath, pathToKey, marshal, unmarshal, staleDur, expiDur)))
}

// スレッドセーフではない。
func newFileTimeLimitedKeyValueStore(path string, keyToPath, pathToKey func(string) string, marshal Marshal, unmarshal Unmarshal, staleDur, expiDur time.Duration) *fileTimeLimitedKeyValueStore {
	return &fileTimeLimitedKeyValueStore{
		NewFileKeyValueStore(path, keyToPath, pathToKey, marshal, unmarshal, staleDur, expiDur),
		NewFileKeyValueStore(path+".expires", keyToPath, pathToKey,
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
	if value, _, err := reg.expires.Get(key, caStmp); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if value == nil {
		return nil, nil, nil
	} else if expires := value.(time.Time); time.Now().After(expires) {
		reg.expires.Remove(key)
		reg.base.Remove(key)
		return nil, nil, nil
	}

	return reg.base.Get(key, caStmp)
}

func (reg *fileTimeLimitedKeyValueStore) Put(key string, value interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
	if _, err := reg.expires.Put(key, expiDate); err != nil {
		return nil, erro.Wrap(err)
	}

	return reg.base.Put(key, value)
}

func (reg *fileTimeLimitedKeyValueStore) Remove(key string) error {
	if err := reg.expires.Remove(key); err != nil {
		return erro.Wrap(err)
	}

	return reg.base.Remove(key)
}
