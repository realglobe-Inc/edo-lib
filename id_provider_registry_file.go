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
func NewFileIdProviderRegistry(path string) IdProviderRegistry {
	return newFileRegistry(path)
}

func (reg *fileRegistry) IdProviderQueryUri(idpUuid string) (queryUri string, err error) {
	path := filepath.Join(reg.path, idpUuid+".json")

	if err := readFromJson(path, &queryUri); err != nil {
		return "", erro.Wrap(err)
	}

	return queryUri, nil
}

// キャッシュ用。
func NewFileDatedIdProviderRegistry(path string, expiDur time.Duration) DatedIdProviderRegistry {
	return newFileBackend(path, expiDur)
}

func (reg *fileBackend) StampedIdProviderQueryUri(idpUuid string, caStmp *Stamp) (queryUri string, newCaStmp *Stamp, err error) {
	path := filepath.Join(reg.path, idpUuid+".json")

	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil, nil
		} else {
			return "", nil, erro.Wrap(err)
		}
	}

	stmp := &Stamp{}
	stmp.Date = fi.ModTime()
	stmp.Digest = strconv.FormatInt(fi.Size(), 10)

	// 対象のスタンプを取得。

	newCaStmp = &Stamp{Date: stmp.Date, ExpiDate: time.Now().Add(reg.expiDur), Digest: stmp.Digest}

	if caStmp != nil && !stmp.Date.After(caStmp.Date) && caStmp.Digest == stmp.Digest {
		return "", newCaStmp, nil
	}

	// 無効なキャッシュだった。

	if err := readFromJson(path, &queryUri); err != nil {
		return "", nil, erro.Wrap(err)
	}

	return queryUri, newCaStmp, nil
}
