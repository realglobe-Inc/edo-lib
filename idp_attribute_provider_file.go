package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// data を JSON として、encoding/json の標準データ型にデコードする。
func jsonUnmarshal(data []byte) (interface{}, error) {
	var res interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return res, nil
}

// スレッドセーフ。
func NewFileIdpAttributeProvider(path string, expiDur time.Duration) IdpAttributeProvider {
	return newIdpAttributeProvider(NewFileKeyValueStore(path, jsonKeyGen, json.Marshal, jsonUnmarshal, expiDur))
}
