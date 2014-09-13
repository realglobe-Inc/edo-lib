package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// バックエンドにファイルシステムを使う。

// TODO 今は手抜きで古いファイルを無視するだけ。どんどん溜まっていく。

func readDate(path string) (time.Time, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return time.Time{}, nil
		}
		return time.Time{}, erro.Wrap(err)
	}
	date, err := time.Parse(time.RFC3339Nano, string(buff))
	if err != nil {
		return time.Time{}, erro.Wrap(err)
	}
	return date, nil
}

func writeDate(date time.Time, path string) error {
	buff := date.Format(time.RFC3339Nano)
	if err := ioutil.WriteFile(path, []byte(buff), filePerm); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

// 非キャッシュ用。
type fileTimeLimitedKeyValueStore struct {
	*fileDriver
}

func NewFileTimeLimitedKeyValueStore(path string) TimeLimitedKeyValueStore {
	return &fileTimeLimitedKeyValueStore{newFileDriver(path)}
}

func (reg *fileTimeLimitedKeyValueStore) Get(key string) (value interface{}, err error) {
	limPath := filepath.Join(reg.path, escapeToFileName(key)+".deadline")
	date, err := readDate(limPath)
	if err != nil {
		return nil, erro.Wrap(err)
	} else if date.IsZero() || date.Before(time.Now()) {
		return nil, nil
	}

	path := filepath.Join(reg.path, escapeToFileName(key)+".json")
	if err := readFromJson(path, &value); err != nil {
		return nil, erro.Wrap(err)
	}

	return value, nil
}

func (reg *fileTimeLimitedKeyValueStore) Put(key string, value interface{}, timLim time.Time) error {
	limPath := filepath.Join(reg.path, escapeToFileName(key)+".deadline")
	if err := writeDate(timLim, limPath); err != nil {
		return erro.Wrap(err)
	}

	path := filepath.Join(reg.path, escapeToFileName(key)+".json")
	if err := writeToJson(path, &value); err != nil {
		return erro.Wrap(err)
	}

	return nil
}

func (reg *fileTimeLimitedKeyValueStore) Remove(key string) error {
	limPath := filepath.Join(reg.path, escapeToFileName(key)+".deadline")
	if err := os.Remove(limPath); err != nil {
		if !os.IsNotExist(err) {
			return erro.Wrap(err)
		}
	}

	path := filepath.Join(reg.path, escapeToFileName(key)+".json")
	if err := os.Remove(path); err != nil {
		if !os.IsNotExist(err) {
			return erro.Wrap(err)
		}
	}

	return nil
}
