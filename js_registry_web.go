package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// data を JSON として、Object にデコードする。
func objectUnmarshal(data []byte) (interface{}, error) {
	var res Object
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return &res, nil
}

// スレッドセーフ。
func NewWebJsRegistry(prefix string) JsRegistry {
	return newJsRegistry(NewWebKeyValueStore(prefix, json.Marshal, objectUnmarshal))
}
