package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// TODO 文字列キーに従ってディレクトリを切る。

// 文字列キーをファイル名に相応しいものにする。
func escapeTo(key string) (path string) {
	// TODO URL パラメータのエスケープで代用してるがもっと良い方法はありそう。
	return url.QueryEscape(key)
}

// キャッシュ用。
type fileRawDataStore struct {
	path      string
	keyToPath func(string) string
	pathToKey func(string) string
	staleDur  time.Duration
	expiDur   time.Duration
}

// スレッドセーフ。
func NewFileRawDataStore(path string, keyToPath, pathToKey func(string) string, staleDur, expiDur time.Duration) RawDataStore {
	return newSynchronizedRawDataStore(newCachingRawDataStore(newFileRawDataStore(path, keyToPath, pathToKey, staleDur, expiDur)))
}

// スレッドセーフではない。
func newFileRawDataStore(path string, keyToPath, pathToKey func(string) string, staleDur, expiDur time.Duration) *fileRawDataStore {
	if keyToPath == nil {
		keyToPath = func(key string) string { return key }
	}
	if pathToKey == nil {
		pathToKey = func(path string) string { return path }
	}
	return &fileRawDataStore{path, keyToPath, pathToKey, staleDur, expiDur}
}

// ダイジェストはタイムスタンプの秒未満にファイルサイズを足して 16 進数で表した文字列。
func (reg *fileRawDataStore) getFileStamp(fi os.FileInfo) *Stamp {
	date := fi.ModTime()
	now := time.Now()
	return &Stamp{
		Date:      date,
		StaleDate: now.Add(reg.staleDur),
		ExpiDate:  now.Add(reg.expiDur),
		Digest:    strconv.FormatInt(int64(date.Nanosecond())+fi.Size(), 16),
	}
}

func (reg *fileRawDataStore) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
	fi, err := os.Stat(reg.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		} else {
			return nil, nil, erro.Wrap(err)
		}
	}

	newCaStmp = reg.getFileStamp(fi)

	// 対象のスタンプを取得。

	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。

	keys = map[string]bool{}
	if err := filepath.Walk(reg.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return erro.Wrap(err)
		} else if info.IsDir() {
			return nil
		}

		key := reg.pathToKey(path[len(reg.path)+1:]) // +1 は / の分。
		if key == "" {
			return nil
		}

		keys[key] = true
		return nil
	}); err != nil {
		return nil, nil, erro.Wrap(err)
	}

	return keys, newCaStmp, nil
}

func (reg *fileRawDataStore) Get(key string, caStmp *Stamp) (data []byte, newCaStmp *Stamp, err error) {
	path := filepath.Join(reg.path, reg.keyToPath(key))

	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		} else {
			return nil, nil, erro.Wrap(err)
		}
	}

	newCaStmp = reg.getFileStamp(fi)

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

func (reg *fileRawDataStore) Put(key string, data []byte) (*Stamp, error) {
	path := filepath.Join(reg.path, reg.keyToPath(key))

	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, filePerm)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, erro.Wrap(err)
		}

		// ディレクトリが無かっただけかもしれないので、
		// ディレクトリを掘って再挑戦。
		if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
			return nil, erro.Wrap(err)
		}
		f, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, filePerm)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	}

	if _, err := f.Write(data); err != nil {
		return nil, erro.Wrap(err)
	}

	// 保存できた。

	fi, err := f.Stat()
	if err != nil {
		return nil, erro.Wrap(err)
	}

	return reg.getFileStamp(fi), nil
}

func (reg *fileRawDataStore) Remove(key string) error {
	path := filepath.Join(reg.path, reg.keyToPath(key))

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return erro.Wrap(err)
	}
	return nil
}
