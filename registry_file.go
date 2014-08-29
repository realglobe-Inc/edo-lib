package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// ファイルを使うモックアップ。
// スレッドセーフではない。

const (
	dirPerm  = 0755
	filePerm = 0644
)

type fileRegistry struct {
	path string
}

func newFileRegistry(path string) *fileRegistry {
	return &fileRegistry{path}
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

// ログイン。
func NewFileLoginRegistry(path string) LoginRegistry {
	return newFileRegistry(path)
}

func (reg *fileRegistry) User(accToken string) (usrUuid string, err error) {
	path := filepath.Join(reg.path, accToken+".json")

	if err := readFromJson(path, &usrUuid); err != nil {
		return "", erro.Wrap(err)
	}
	return usrUuid, nil
}

// JavaScript.
func NewFileJsRegistry(path string) JsRegistry {
	return newFileRegistry(path)
}

type objectHeader struct {
	Service bool     `json:"service,omitempty"`
	Library bool     `json:"library,omitempty"`
	Include []string `json:"include,omitempty"`
}

func (reg *fileRegistry) Object(dir, objName string) (*Object, error) {
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
func (reg *fileRegistry) AddObject(dir, objName string, obj *Object) error {
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
func (reg *fileRegistry) RemoveObject(dir, objName string) error {
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

// ユーザー情報。
func NewFileUserRegistry(path string) UserRegistry {
	return newFileRegistry(path)
}

// 属性名は任意の文字列でファイル名にしづらいのでユーザーごとに 1 ファイル。
func (reg *fileRegistry) Attributes(usrUuid string) (map[string]interface{}, error) {
	path := filepath.Join(reg.path, usrUuid+".json")

	var attrs map[string]interface{}
	if err := readFromJson(path, &attrs); err != nil {
		return nil, erro.Wrap(err)
	}

	return attrs, nil
}
func (reg *fileRegistry) Attribute(usrUuid, attrName string) (interface{}, error) {
	attrs, err := reg.Attributes(usrUuid)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return attrs[attrName], nil
}
func (reg *fileRegistry) AddAttribute(usrUuid, attrName string, attr interface{}) error {
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
func (reg *fileRegistry) RemoveAttribute(usrUuid, attrName string) error {
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

// ジョブ。
func NewFileJobRegistry(path string) JobRegistry {
	return newFileRegistry(path)
}

func (reg *fileRegistry) Result(jobId string) (*JobResult, error) {
	path := filepath.Join(reg.path, jobId+".json")

	var res JobResult
	if err := readFromJson(path, &res); err != nil {
		return nil, erro.Wrap(err)
	}

	if res.Status == 0 {
		return nil, nil
	}
	return &res, nil
}
func (reg *fileRegistry) AddResult(jobId string, res *JobResult, deadline time.Time) error {
	path := filepath.Join(reg.path, jobId+".json")

	if err := writeToJson(path, res); err != nil {
		return erro.Wrap(err)
	}

	return nil
}

// 別名。
func NewFileNameRegistry(path string) NameRegistry {
	return newFileRegistry(path)
}

func (reg *fileRegistry) Address(name string) (addr string, err error) {
	path := filepath.Join(reg.path, name+".json")

	if err := readFromJson(path, &addr); err != nil { // 改行とかに煩わされないので JSON 文字列で。
		return "", erro.Wrap(err)
	}
	return addr, nil
}
func (reg *fileRegistry) Addresses(name string) (addrs []string, err error) {
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

	tree := newNameTree()
	tree.fromContainer(cont)

	return tree.addresses(name), nil
}

// イベント。
func NewFileEventRegistry(path string) EventRegistry {
	return newFileRegistry(path)
}

// イベントは区切りに / を含み、ディレクトリを掘るのは面倒なので、ユーザーごとに 1 ファイル。
func (reg *fileRegistry) Handler(usrUuid, event string) (Handler, error) {
	path := filepath.Join(reg.path, usrUuid+".json")

	var cont map[string]Handler
	if err := readFromJson(path, &cont); err != nil {
		return nil, erro.Wrap(err)
	}

	tree := newEventTree()
	tree.fromContainer(cont)

	return tree.handler(event), nil
}
func (reg *fileRegistry) AddHandler(usrUuid, event string, hndl Handler) error {
	path := filepath.Join(reg.path, usrUuid+".json")

	var cont map[string]Handler
	if err := readFromJson(path, &cont); err != nil {
		return erro.Wrap(err)
	}

	if reflect.DeepEqual(cont[event], hndl) {
		return nil
	}

	if cont == nil {
		cont = map[string]Handler{event: hndl}
	} else {
		cont[event] = hndl
	}
	if err := writeToJson(path, cont); err != nil {
		return erro.Wrap(err)
	}

	return nil
}
func (reg *fileRegistry) RemoveHandler(usrUuid, event string) error {
	path := filepath.Join(reg.path, usrUuid+".json")

	var cont map[string]Handler
	if err := readFromJson(path, &cont); err != nil {
		return erro.Wrap(err)
	}

	if _, ok := cont[event]; !ok {
		return nil
	}

	delete(cont, event)
	if err := writeToJson(path, cont); err != nil {
		return erro.Wrap(err)
	}

	return nil
}

// サービス。
func NewFileServiceRegistry(path string) ServiceRegistry {
	return newFileRegistry(path)
}

func (reg *fileRegistry) Service(endPt string) (servUuid string, err error) {
	path := filepath.Join(reg.path, "uuid.json")

	var cont map[string]string
	if err := readFromJson(path, &cont); err != nil {
		return "", erro.Wrap(err)
	}

	tree := newServiceTree()
	tree.fromContainer(cont)

	return tree.service(endPt), nil
}

// サービス。
func NewFileIdProviderRegistry(path string) IdProviderRegistry {
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
