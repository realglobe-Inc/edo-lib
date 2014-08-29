package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// JavaScript.
func NewFileJsBackendRegistry(path string) JsBackendRegistry {
	return newFileRegistry(path)
}

func (reg *fileRegistry) StampedObject(dir, objName string, caStmp *Stamp) (*Object, *Stamp, error) {
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

	newCaStmp := &Stamp{Date: time.Now(), Digest: stmp.Digest}

	if caStmp != nil && caStmp.Date.After(stmp.Date) && caStmp.Digest == stmp.Digest {
		return nil, newCaStmp, nil
	}

	// 無効なキャッシュだった。

	obj, err := reg.Object(dir, objName)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}
	return obj, newCaStmp, nil
}

// ID プロバイダ。
func NewFileIdProviderBackend(path string) IdProviderBackend {
	return newFileRegistry(path)
}

func (reg *fileRegistry) StampedIdProviders(caStmp *Stamp) ([]*IdProvider, *Stamp, error) {
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

	newCaStmp := &Stamp{Date: time.Now(), Digest: stmp.Digest}

	if caStmp != nil && caStmp.Date.After(stmp.Date) && caStmp.Digest == stmp.Digest {
		return nil, newCaStmp, nil
	}

	// 無効なキャッシュだった。

	var cont []*IdProvider
	if err := readFromJson(path, &cont); err != nil {
		return nil, nil, erro.Wrap(err)
	}

	return cont, newCaStmp, nil
}
