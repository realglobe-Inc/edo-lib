package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// バックエンドにファイルシステムを使う。

// 非キャッシュ用。
func NewFileServiceKeyRegistry(path string) ServiceKeyRegistry {
	return newFileRegistry(path)
}

func (reg *fileRegistry) ServiceKey(servUuid string) (key string, err error) {
	path := filepath.Join(reg.path, servUuid+".json")

	if err := readFromJson(path, &key); err != nil {
		return "", erro.Wrap(err)
	}

	return key, nil
}

// キャッシュ用。
func NewFileDatedServiceKeyRegistry(path string, expiDur time.Duration) DatedServiceKeyRegistry {
	return newFileBackend(path, expiDur)
}

func (reg *fileBackend) StampedServiceKey(servUuid string, caStmp *Stamp) (key string, newCaStmp *Stamp, err error) {
	path := filepath.Join(reg.path, servUuid+".json")

	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil, nil
		} else {
			return "", nil, erro.Wrap(err)
		}
	}

	// 対象のスタンプを取得。

	newCaStmp = &Stamp{Date: fi.ModTime(), ExpiDate: time.Now().Add(reg.expiDur), Digest: strconv.FormatInt(fi.Size(), 10)}

	if caStmp != nil && !newCaStmp.Date.After(caStmp.Date) && caStmp.Digest == newCaStmp.Digest {
		return "", newCaStmp, nil
	}

	// 無効なキャッシュだった。

	if err := readFromJson(path, &key); err != nil {
		return "", nil, erro.Wrap(err)
	}

	return key, newCaStmp, nil
}
