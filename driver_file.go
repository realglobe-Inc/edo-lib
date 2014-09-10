package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io/ioutil"
	"os"
	"time"
)

// ファイルを使うモックアップ。
// スレッドセーフではない。

const (
	dirPerm  = 0755
	filePerm = 0644
)

// 非キャッシュ用。
type fileDriver struct {
	path string
}

func newFileDriver(path string) *fileDriver {
	return &fileDriver{path}
}

func readFromJson(path string, v interface{}) error {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return erro.Wrap(err)
	}
	if err := json.Unmarshal(buff, v); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func writeToJson(path string, v interface{}) error {
	buff, err := json.Marshal(v)
	if err != nil {
		return erro.Wrap(err)
	}
	if err := ioutil.WriteFile(path, buff, filePerm); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

// キャッシュ用。
type datedFileDriver struct {
	*fileDriver
	expiDur time.Duration
}

func newDatedFileDriver(path string, expiDur time.Duration) *datedFileDriver {
	return &datedFileDriver{newFileDriver(path), expiDur}
}
