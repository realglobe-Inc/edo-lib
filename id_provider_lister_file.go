package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// data を JSON として、[]*IdProvider にデコードする。
func idProvidersUnmarshal(data []byte) (interface{}, error) {
	var idps []*IdProvider
	if err := json.Unmarshal(data, &idps); err != nil {
		return nil, erro.Wrap(err)
	}
	return idps, nil
}

// スレッドセーフ。
func NewFileIdProviderLister(path string, expiDur time.Duration) IdProviderLister {
	return newIdProviderLister(NewFileKeyValueStore(path, jsonKeyGen, json.Marshal, idProvidersUnmarshal, expiDur))
}
