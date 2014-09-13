package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// 非キャッシュ用。
func NewFileJsRegistry(path string) JsRegistry {
	return newFileDriver(path)
}

type objectHeader struct {
	Service bool     `json:"service,omitempty"`
	Library bool     `json:"library,omitempty"`
	Include []string `json:"include,omitempty"`
}

func (reg *fileDriver) Object(dir, objName string) (*Object, error) {
	headPath := filepath.Join(reg.path, dir, objName+".json")
	codePath := filepath.Join(reg.path, dir, objName+".js")

	var head objectHeader
	if err := readFromJson(headPath, &head); err != nil {
		return nil, erro.Wrap(err)
	}
	code, err := ioutil.ReadFile(codePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, erro.Wrap(err)
	}

	return &Object{head.Service, head.Library, head.Include, string(code)}, nil
}
func (reg *fileDriver) AddObject(dir, objName string, obj *Object) error {
	headPath := filepath.Join(reg.path, dir, objName+".json")
	codePath := filepath.Join(reg.path, dir, objName+".js")

	if err := os.MkdirAll(filepath.Join(reg.path, dir), dirPerm); err != nil {
		return erro.Wrap(err)
	}
	if err := writeToJson(headPath, &objectHeader{obj.Service, obj.Library, obj.Include}); err != nil {
		return erro.Wrap(err)
	}
	if err := ioutil.WriteFile(codePath, []byte(obj.Code), filePerm); err != nil {
		return erro.Wrap(err)
	}

	return nil
}
func (reg *fileDriver) RemoveObject(dir, objName string) error {
	headPath := filepath.Join(reg.path, dir, objName+".json")
	codePath := filepath.Join(reg.path, dir, objName+".js")

	if err := os.Remove(headPath); err != nil {
		if !os.IsNotExist(err) {
			return erro.Wrap(err)
		}
	}
	if err := os.Remove(codePath); err != nil {
		if !os.IsNotExist(err) {
			return erro.Wrap(err)
		}
	}

	return nil
}

// キャッシュ用。
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
