package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// キャッシュ用。
type fileRawDataStore struct {
	path      string
	keyToPath func(string) string
	staleDur  time.Duration
	expiDur   time.Duration
}

// スレッドセーフではない。
func newFileRawDataStore(path string, keyToPath func(string) string, staleDur, expiDur time.Duration) *fileRawDataStore {
	if keyToPath == nil {
		keyToPath = func(key string) string { return key }
	}
	return &fileRawDataStore{path, keyToPath, staleDur, expiDur}
}

func (drv *fileRawDataStore) getStamp(fi os.FileInfo) *Stamp {
	now := time.Now()
	stmp := getFileStamp(fi)
	stmp.StaleDate = now.Add(drv.staleDur)
	stmp.ExpiDate = now.Add(drv.expiDur)
	return stmp
}

func (drv *fileRawDataStore) Get(key string, caStmp *Stamp) (data []byte, newCaStmp *Stamp, err error) {
	path := filepath.Join(drv.path, drv.keyToPath(key))

	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		} else {
			return nil, nil, erro.Wrap(err)
		}
	}

	newCaStmp = drv.getStamp(fi)

	// 対象のスタンプを取得。

	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。

	data, err = ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, erro.Wrap(err)
	}
	return data, newCaStmp, nil
}

func (drv *fileRawDataStore) Put(key string, data []byte) (*Stamp, error) {
	path := filepath.Join(drv.path, drv.keyToPath(key))

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, erro.Wrap(err)
		}

		// ディレクトリが無かっただけかもしれないので、
		// ディレクトリを掘って再挑戦。
		if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
			return nil, erro.Wrap(err)
		}
		f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, filePerm)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return nil, erro.Wrap(err)
	}

	// 保存できた。

	fi, err := f.Stat()
	if err != nil {
		return nil, erro.Wrap(err)
	} else if int64(len(data)) < fi.Size() {
		// 前の内容の方が大きかった。
		if err := f.Truncate(int64(len(data))); err != nil {
			return nil, erro.Wrap(err)
		}
		fi, err = f.Stat()
		if err != nil {
			return nil, erro.Wrap(err)
		}
	}
	return drv.getStamp(fi), nil
}

func (drv *fileRawDataStore) Remove(key string) error {
	path := filepath.Join(drv.path, drv.keyToPath(key))

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return erro.Wrap(err)
	}
	return nil
}

func (drv *fileRawDataStore) Close() error {
	return nil
}
