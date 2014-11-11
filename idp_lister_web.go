package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// {
//     "id_providers": [
//         id-provider-no-hitotsu,
//         ...
//     ]
// }
func webIdProvidersUnmarshal(data []byte) (interface{}, error) {
	var res struct {
		Id_providers []*IdProvider
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return res.Id_providers, nil
}

// スレッドセーフ。
func NewWebIdpLister(prefix string) IdpLister {
	return newIdpLister(NewWebKeyValueStore(prefix, nil, webIdProvidersUnmarshal))
}
