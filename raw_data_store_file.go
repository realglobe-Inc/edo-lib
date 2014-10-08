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

// TODO key に従ってディレクトリを切る。

// 任意の文字列をファイル名に相応しいものにする。
func escapeToFileName(before string) (after string) {
	// TODO URL パラメータのエスケープで代用してるがもっと良い方法はありそう。
	return url.QueryEscape(before)
}

// キャッシュ用。
type fileRawDataStore struct {
	path    string
	keyGen  func(string) string
	expiDur time.Duration
}

// スレッドセーフ。
func NewFileRawDataStore(path string, keyGen func(string) string, expiDur time.Duration) RawDataStore {
	return newSynchronizedRawDataStore(newCachingRawDataStore(newFileRawDataStore(path, keyGen, expiDur)))
}

// スレッドセーフではない。
func newFileRawDataStore(path string, keyGen func(string) string, expiDur time.Duration) *fileRawDataStore {
	if keyGen == nil {
		keyGen = escapeToFileName
	} else {
		oldKeyGen := keyGen
		keyGen = func(before string) string {
			return escapeToFileName(oldKeyGen(before))
		}
	}
	return &fileRawDataStore{path, keyGen, expiDur}
}

// ダイジェストはファイルサイズを 10 進数で表したもの。
// ダイジェストが違えばファイルが違うような値として、他に手軽に取れるものがないため。

func (reg *fileRawDataStore) Get(key string, caStmp *Stamp) (data []byte, newCaStmp *Stamp, err error) {
	path := filepath.Join(reg.path, reg.keyGen(key))

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
	path := filepath.Join(reg.path, reg.keyGen(key))

	for {
		if err := ioutil.WriteFile(path, data, filePerm); err != nil {
			return nil, erro.Wrap(err)
		}

		// 保存できた。

		fi, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				// 保存したあとに消された。使い方が悪い。やり直し。
				continue
			} else {
				return nil, erro.Wrap(err)
			}
		}

		// 対象のスタンプを取得。

		newCaStmp := &Stamp{Date: fi.ModTime(), ExpiDate: time.Now().Add(reg.expiDur), Digest: strconv.FormatInt(fi.Size(), 10)}
		return newCaStmp, nil
	}
}

func (reg *fileRawDataStore) Remove(key string) error {
	path := filepath.Join(reg.path, reg.keyGen(key))

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return erro.Wrap(err)
	}
	return nil
}
