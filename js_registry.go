package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

type Object struct {
	Service bool `json:"service,omitempty"` // Web API からの利用を許可するか。
	Library bool `json:"library,omitempty"` // ライブラリとしての利用を許可するか。

	// 直接利用するライブラリとコネクタ。
	// Call するオブジェクトは含まない。
	// ライブラリは $library、コネクタは $$connector 形式で。
	Include []string `json:"include,omitempty"`

	Code string `json:"code"`
}

type JsRegistry interface {
	// オブジェクトのソースを取得する。
	Object(dir, objName string, caStmp *Stamp) (*Object, *Stamp, error)

	// オブジェクトのソースを登録する。一時的な利用を想定。
	AddObject(dir, objName string, obj *Object) (*Stamp, error)
	// オブジェクトのソースを削除する。
	RemoveObject(dir, objName string) error
}

type jsRegistry struct {
	base KeyValueStore
}

func newJsRegistry(base KeyValueStore) *jsRegistry {
	return &jsRegistry{base}
}

func (reg *jsRegistry) Object(dir, objName string, caStmp *Stamp) (*Object, *Stamp, error) {
	value, newCaStmp, err := reg.base.Get(dir+"/"+objName, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if value == nil {
		return nil, newCaStmp, nil
	}
	return value.(*Object), newCaStmp, nil
}

func (reg *jsRegistry) AddObject(dir, objName string, obj *Object) (*Stamp, error) {
	return reg.base.Put(dir+"/"+objName, obj)
}

func (reg *jsRegistry) RemoveObject(dir, objName string) error {
	return reg.base.Remove(dir + "/" + objName)
}
