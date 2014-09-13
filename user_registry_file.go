package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"path/filepath"
	"reflect"
)

// 非キャッシュ用。
func NewFileUserRegistry(path string) UserRegistry {
	return newFileDriver(path)
}

// 属性名は任意の文字列でファイル名にしづらいのでユーザーごとに 1 ファイル。
func (reg *fileDriver) Attributes(usrUuid string) (map[string]interface{}, error) {
	path := filepath.Join(reg.path, usrUuid+".json")

	var attrs map[string]interface{}
	if err := readFromJson(path, &attrs); err != nil {
		return nil, erro.Wrap(err)
	}

	return attrs, nil
}
func (reg *fileDriver) Attribute(usrUuid, attrName string) (interface{}, error) {
	attrs, err := reg.Attributes(usrUuid)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return attrs[attrName], nil
}
func (reg *fileDriver) AddAttribute(usrUuid, attrName string, attr interface{}) error {
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
func (reg *fileDriver) RemoveAttribute(usrUuid, attrName string) error {
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
