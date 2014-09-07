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
func NewFileIdProviderLister(path string) IdProviderLister {
	return newFileRegistry(path)
}

func (reg *fileRegistry) IdProviders() ([]*IdProvider, error) {
	path := filepath.Join(reg.path, "idp.json")

	var cont []*IdProvider
	if err := readFromJson(path, &cont); err != nil {
		return nil, erro.Wrap(err)
	}

	return cont, nil
}

// キャッシュ用。
func NewFileDatedIdProviderLister(path string, expiDur time.Duration) DatedIdProviderLister {
	return newFileBackend(path, expiDur)
}

func (reg *fileBackend) StampedIdProviders(caStmp *Stamp) ([]*IdProvider, *Stamp, error) {
	path := filepath.Join(reg.path, "idp.json")

	fi, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, nil, erro.Wrap(err)
		}
	}

	stmp := &Stamp{}
	stmp.Date = fi.ModTime()
	stmp.Digest = strconv.FormatInt(fi.Size(), 10)

	// 対象のスタンプを取得。

	newCaStmp := &Stamp{Date: stmp.Date, ExpiDate: time.Now().Add(reg.expiDur), Digest: stmp.Digest}

	if caStmp != nil && !stmp.Date.After(caStmp.Date) && caStmp.Digest == stmp.Digest {
		return nil, newCaStmp, nil
	}

	// 無効なキャッシュだった。

	var cont []*IdProvider
	if err := readFromJson(path, &cont); err != nil {
		return nil, nil, erro.Wrap(err)
	}

	return cont, newCaStmp, nil
}
