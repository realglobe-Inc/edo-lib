package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// ファイルを使うモックアップ。
// スレッドセーフではない。

const (
	dirPerm  = 0755
	filePerm = 0644
)

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

// JavaScript.
type fileJsRegistry struct {
	path string
}

func NewFileJsRegistry(path string) JsRegistry {
	return &fileJsRegistry{path}
}

type objectHeader struct {
	Service bool     `json:"service,omitempty"`
	Library bool     `json:"library,omitempty"`
	Include []string `json:"include,omitempty"`
}

func (reg *fileJsRegistry) Object(dir, objName string) (*Object, error) {
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
func (reg *fileJsRegistry) AddObject(dir, objName string, obj *Object) error {
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
func (reg *fileJsRegistry) RemoveObject(dir, objName string) error {
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

// ログイン。
type fileLoginRegistry struct {
	path string
}

func NewFileLoginRegistry(path string) LoginRegistry {
	return &fileLoginRegistry{path}
}

func (reg *fileLoginRegistry) User(accToken string) (usrUuid string, err error) {
	path := filepath.Join(reg.path, accToken+".json")

	if err := readFromJson(path, &usrUuid); err != nil {
		return "", erro.Wrap(err)
	}
	return usrUuid, nil
}

// ユーザー情報。
type fileUserRegistry struct {
	path string
}

func NewFileUserRegistry(path string) UserRegistry {
	return &fileUserRegistry{path}
}

// 属性名は任意の文字列でファイル名にしづらいのでユーザーごとに 1 ファイル。
func (reg *fileUserRegistry) Attributes(usrUuid string) (map[string]interface{}, error) {
	path := filepath.Join(reg.path, usrUuid+".json")

	var attrs map[string]interface{}
	if err := readFromJson(path, &attrs); err != nil {
		return nil, erro.Wrap(err)
	}

	return attrs, nil
}
func (reg *fileUserRegistry) Attribute(usrUuid, attrName string) (interface{}, error) {
	attrs, err := reg.Attributes(usrUuid)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return attrs[attrName], nil
}
func (reg *fileUserRegistry) AddAttribute(usrUuid, attrName string, attr interface{}) error {
	attrs, err := reg.Attributes(usrUuid)
	if err != nil {
		return erro.Wrap(err)
	} else if attrs == nil {
		attrs = map[string]interface{}{}
	}

	if reflect.DeepEqual(attrs[attrName], attr) {
		return nil
	}

	path := filepath.Join(reg.path, usrUuid+".json")
	attrs[attrName] = attr
	if err := writeToJson(path, attrs); err != nil {
		return erro.Wrap(err)
	}

	return nil
}
func (reg *fileUserRegistry) RemoveAttribute(usrUuid, attrName string) error {
	attrs, err := reg.Attributes(usrUuid)
	if err != nil {
		return erro.Wrap(err)
	}

	if _, ok := attrs[attrName]; !ok {
		return nil
	}

	path := filepath.Join(reg.path, usrUuid+".json")
	delete(attrs, attrName)
	if err := writeToJson(path, attrs); err != nil {
		return erro.Wrap(err)
	}

	return nil
}

// 別名。
type fileNameRegistry struct {
	path string
}

func NewFileNameRegistry(path string) NameRegistry {
	return &fileNameRegistry{path}
}

func (reg *fileNameRegistry) Address(name string) (addr string, err error) {
	path := filepath.Join(reg.path, name+".json")

	if err := readFromJson(path, &addr); err != nil { // 改行とかに煩わされないので JSON 文字列で。
		return "", erro.Wrap(err)
	}
	return addr, nil
}
func (reg *fileNameRegistry) Addresses(name string) (addrs []string, err error) {
	cont := map[string]string{}

	fis, err := ioutil.ReadDir(reg.path)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	for _, fi := range fis {
		if fi.IsDir() {
			continue
		} else if !strings.HasSuffix(fi.Name(), ".json") {
			continue
		}

		curName := strings.TrimSuffix(fi.Name(), ".json")

		if !strings.HasSuffix(curName, name) {
			// 部分木以外はスルー。
			continue
		}

		path := filepath.Join(reg.path, fi.Name())

		var addr string
		if err := readFromJson(path, &addr); err != nil {
			return nil, erro.Wrap(err)
		}

		cont[curName] = addr
	}

	tree := nameTree{}
	tree.fromContainer(cont)

	return tree.addresses(name), nil
}

// イベント.
type fileEventRegistry struct {
	path string
}

func NewFileEventRegistry(path string) EventRegistry {
	return &fileEventRegistry{path}
}

func (reg *fileEventRegistry) Handler(usrUuid, event string) (Handler, error) {
	cont := map[string]Handler{}

	dir := filepath.Join(reg.path, usrUuid)
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	for _, fi := range fis {
		if fi.IsDir() {
			continue
		} else if !strings.HasSuffix(fi.Name(), ".json") {
			continue
		}

		curEvent := strings.TrimSuffix(fi.Name(), ".json")

		if !strings.HasPrefix(curEvent, event) {
			// 部分木以外はスルー。
			continue
		}

		path := filepath.Join(dir, fi.Name())

		var hndl Handler
		if err := readFromJson(path, &hndl); err != nil {
			return nil, erro.Wrap(err)
		}

		cont[curEvent] = hndl
	}

	tree := eventTree{}
	tree.fromContainer(cont)

	return tree.handler(event), nil
}
func (reg *fileEventRegistry) AddHandler(usrUuid, event string, hndl Handler) error {
	dir := filepath.Join(reg.path, usrUuid)
	path := filepath.Join(usrUuid, event+".json")

	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return erro.Wrap(err)
	}
	if err := writeToJson(path, hndl); err != nil {
		return erro.Wrap(err)
	}

	return nil
}
func (reg *fileEventRegistry) RemoveHandler(usrUuid, event string) error {
	path := filepath.Join(reg.path, usrUuid, event+".json")

	if err := os.Remove(path); err != nil {
		if !os.IsNotExist(err) {
			return erro.Wrap(err)
		}
	}

	return nil
}
