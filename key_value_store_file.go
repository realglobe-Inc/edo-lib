package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// バックエンドにファイルシステムを使う。

func escapeToFileName(before string) (after string) {
	// TODO URL パラメータのエスケープで代用してるがもっと良い方法はありそう。
	return url.QueryEscape(before)
}

// 非キャッシュ用。
func newFileKeyValueStore(path string) keyValueStore {
	return newFileDriver(path)
}

func (reg *fileDriver) get(key string) (value interface{}, err error) {
	path := filepath.Join(reg.path, escapeToFileName(key)+".json")

	if err := readFromJson(path, &value); err != nil {
		return nil, erro.Wrap(err)
	}

	return value, nil
}

func (reg *fileDriver) put(key string, value interface{}) error {
	path := filepath.Join(reg.path, escapeToFileName(key)+".json")

	if err := writeToJson(path, &value); err != nil {
		return erro.Wrap(err)
	}

	return nil
}

func (reg *fileDriver) remove(key string) error {
	path := filepath.Join(reg.path, escapeToFileName(key)+".json")

	if err := os.Remove(path); err != nil {
		if !os.IsNotExist(err) {
			return erro.Wrap(err)
		}
	}

	return nil
}

// キャッシュ用。
func newFileDatedKeyValueStore(path string, expiDur time.Duration) datedKeyValueStore {
	return newDatedFileDriver(path, expiDur)
}

func (reg *datedFileDriver) stampedGet(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	path := filepath.Join(reg.path, escapeToFileName(key)+".json")

	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		} else {
			return nil, nil, erro.Wrap(err)
		}
	}

	// 対象のスタンプを取得。

	newCaStmp = &Stamp{Date: fi.ModTime(), ExpiDate: time.Now().Add(reg.expiDur), Digest: strconv.FormatInt(fi.Size(), 10)}

	if caStmp != nil && !newCaStmp.Date.After(caStmp.Date) && caStmp.Digest == newCaStmp.Digest {
		return nil, newCaStmp, nil
	}

	// 無効なキャッシュだった。

	if err := readFromJson(path, &value); err != nil {
		return nil, nil, erro.Wrap(err)
	}

	return value, newCaStmp, nil
}

func (reg *datedFileDriver) stampedPut(key string, value interface{}) (*Stamp, error) {
	path := filepath.Join(reg.path, escapeToFileName(key)+".json")

	if err := writeToJson(path, &value); err != nil {
		return nil, erro.Wrap(err)
	}

	// 保存できた。

	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, erro.Wrap(err)
		}
	}

	// 対象のスタンプを取得。

	newCaStmp := &Stamp{Date: fi.ModTime(), ExpiDate: time.Now().Add(reg.expiDur), Digest: strconv.FormatInt(fi.Size(), 10)}
	return newCaStmp, nil
}
