package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// JavaScript.
func NewFileJsBackendRegistry(path string, expiDur time.Duration) JsBackendRegistry {
	return newDatedFileDriver(path, expiDur)
}

func (reg *datedFileDriver) StampedObject(dir, objName string, caStmp *Stamp) (*Object, *Stamp, error) {
	headPath := filepath.Join(reg.path, dir, objName+".json")
	codePath := filepath.Join(reg.path, dir, objName+".js")

	headFi, err := os.Stat(headPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, nil, erro.Wrap(err)
		}
	}
	codeFi, err := os.Stat(codePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, nil, erro.Wrap(err)
		}
	}

	if codeFi == nil {
		return nil, nil, nil
	}

	stmp := &Stamp{}
	if headFi == nil {
		stmp.Date = codeFi.ModTime()
		stmp.Digest = strconv.FormatInt(codeFi.Size(), 10)
	} else {
		if headFi.ModTime().After(codeFi.ModTime()) {
			stmp.Date = headFi.ModTime()
		} else {
			stmp.Date = codeFi.ModTime()
		}
		stmp.Digest = strconv.FormatInt(headFi.Size()+codeFi.Size(), 10)
	}

	// 対象のスタンプを取得。

	newCaStmp := &Stamp{Date: stmp.Date, ExpiDate: time.Now().Add(reg.expiDur), Digest: stmp.Digest}

	if caStmp != nil && !stmp.Date.After(caStmp.Date) && caStmp.Digest == stmp.Digest {
		return nil, newCaStmp, nil
	}

	// 無効なキャッシュだった。

	obj, err := reg.Object(dir, objName)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}
	return obj, newCaStmp, nil
}
