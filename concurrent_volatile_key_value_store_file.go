package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

func dateMarshal(val interface{}) ([]byte, error) {
	return []byte(val.(time.Time).Format(time.RFC3339Nano)), nil
}

func dateUnmarshal(data []byte) (interface{}, error) {
	date, err := time.Parse(time.RFC3339Nano, string(data))
	if err != nil {
		return time.Time{}, erro.Wrap(err)
	}
	return date, nil
}

type fileEntry struct {
	Val string    `json:"value"`
	Exp time.Time `json:"expires"`
}

func fileEntryUnmarshal(data []byte) (interface{}, error) {
	var ent fileEntry
	if err := json.Unmarshal(data, &ent); err != nil {
		return nil, erro.Wrap(err)
	}
	return &ent, nil
}

// TODO 今は手抜きで古いファイルを無視するだけ。どんどん溜まっていく。

type fileConcurrentVolatileKeyValueStore struct {
	base KeyValueStore
	exps KeyValueStore

	ents KeyValueStore
}

// スレッドセーフ。
func NewFileConcurrentVolatileKeyValueStore(path, expiPath string, keyToPath, pathToKey func(string) string, marshal Marshal, unmarshal Unmarshal, staleDur, expiDur time.Duration) ConcurrentVolatileKeyValueStore {
	return newSynchronizedVolatileKeyValueStore(newCachingConcurrentVolatileKeyValueStore(newFileConcurrentVolatileKeyValueStore(path, expiPath, keyToPath, pathToKey, marshal, unmarshal, staleDur, expiDur)))
}

// スレッドセーフではない。
func newFileConcurrentVolatileKeyValueStore(path, expiPath string, keyToPath, pathToKey func(string) string, marshal Marshal, unmarshal Unmarshal, staleDur, expiDur time.Duration) *fileConcurrentVolatileKeyValueStore {
	return &fileConcurrentVolatileKeyValueStore{
		newFileListedKeyValueStore(path, keyToPath, pathToKey, marshal, unmarshal, staleDur, expiDur),
		newFileListedKeyValueStore(expiPath, keyToPath, pathToKey, dateMarshal, dateUnmarshal, staleDur, expiDur),
		newFileListedKeyValueStore(expiPath, keyToPath, pathToKey, json.Marshal, fileEntryUnmarshal, staleDur, expiDur),
	}
}

func (drv *fileConcurrentVolatileKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	var expiDate time.Time
	if exp, newCaStmp, err := drv.exps.Get(key, nil); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if newCaStmp == nil {
		return nil, nil, nil
	} else {
		expiDate = exp.(time.Time)
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

func (drv *fileConcurrentVolatileKeyValueStore) Put(key string, val interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
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

func (drv *fileConcurrentVolatileKeyValueStore) Remove(key string) error {
	if err := drv.exps.Remove(key); err != nil {
		return erro.Wrap(err)
	}
	return drv.base.Remove(key)
}

func (drv *fileConcurrentVolatileKeyValueStore) Entry(eKey string) (eVal string, err error) {
	v, _, err := drv.ents.Get(eKey, nil)
	if err != nil {
		return "", erro.Wrap(err)
	}
	ent, _ := v.(*fileEntry)

	if ent == nil {
		return "", nil
	} else if time.Now().After(ent.Exp) {
		drv.ents.Remove(eKey)
		return "", nil
	}
	return ent.Val, nil
}

func (drv *fileConcurrentVolatileKeyValueStore) SetEntry(eKey, eVal string, eExpiDate time.Time) error {
	if _, err := drv.ents.Put(eKey, &fileEntry{eVal, eExpiDate}); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func (drv *fileConcurrentVolatileKeyValueStore) GetAndSetEntry(key string, caStmp *Stamp, eKey, eVal string, eExpiDate time.Time) (val interface{}, newCaStmp *Stamp, err error) {
	val, newCaStmp, err = drv.Get(key, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	if err := drv.SetEntry(eKey, eVal, eExpiDate); err != nil {
		return nil, nil, erro.Wrap(err)
	}

	return val, newCaStmp, nil
}

func (drv *fileConcurrentVolatileKeyValueStore) PutIfEntered(key string, val interface{}, expiDate time.Time, eKey, eVal string) (entered bool, newCaStmp *Stamp, err error) {
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
