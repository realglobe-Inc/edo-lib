package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"path/filepath"
	"reflect"
)

// 非キャッシュ用。
func NewFileEventRegistry(path string) EventRegistry {
	return newFileDriver(path)
}

// イベントは区切りに / を含み、ディレクトリを掘るのは面倒なので、ユーザーごとに 1 ファイル。
func (reg *fileDriver) Handler(usrUuid, event string) (Handler, error) {
	path := filepath.Join(reg.path, usrUuid+".json")

	var cont map[string]Handler
	if err := readFromJson(path, &cont); err != nil {
		return nil, erro.Wrap(err)
	}

	tree := newEventTree()
	tree.fromContainer(cont)

	return tree.handler(event), nil
}
func (reg *fileDriver) AddHandler(usrUuid, event string, hndl Handler) error {
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
func (reg *fileDriver) RemoveHandler(usrUuid, event string) error {
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
