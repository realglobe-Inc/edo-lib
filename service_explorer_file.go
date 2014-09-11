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
func NewFileServiceExplorer(path string) ServiceExplorer {
	return newSynchronizedServiceExplorer(newFileDriver(path))
}

func (reg *fileDriver) ServiceUuid(servUri string) (servUuid string, err error) {
	path := filepath.Join(reg.path, "list.json")

	var cont map[string]string
	if err := readFromJson(path, &cont); err != nil {
		return "", erro.Wrap(err)
	}

	tree := newServiceExplorerTree()
	tree.fromContainer(cont)

	return tree.get(servUri), nil
}

// キャッシュ用。
func NewFileDatedServiceExplorer(path string, expiDur time.Duration) DatedServiceExplorer {
	return newSynchronizedDatedServiceExplorer(newDatedFileDriver(path, expiDur))
}

func (reg *datedFileDriver) StampedServiceUuid(servUri string, caStmp *Stamp) (servUuid string, newCaStmp *Stamp, err error) {
	path := filepath.Join(reg.path, "list.json")

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

	var cont map[string]string
	if err := readFromJson(path, &cont); err != nil {
		return "", nil, erro.Wrap(err)
	}

	tree := newServiceExplorerTree()
	tree.fromContainer(cont)

	return tree.get(servUri), newCaStmp, nil
}
