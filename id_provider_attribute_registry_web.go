package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// {
//     "id_provider": {
//         attrNameX: XXX
//     }
// }
func webIdProviderAttributeUnmarshal(data []byte) (interface{}, error) {
	var res struct {
		Id_provider map[string]interface{}
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return res.Id_provider, nil
}

type webIdProviderAttributeRegistry struct {
	KeyValueStore
}

// スレッドセーフ。
func NewWebIdProviderAttributeRegistry(prefix string) IdProviderAttributeRegistry {
	return newWebIdProviderAttributeRegistry(NewWebKeyValueStore(prefix, nil, webIdProviderAttributeUnmarshal))
}

func newWebIdProviderAttributeRegistry(base KeyValueStore) *webIdProviderAttributeRegistry {
	return &webIdProviderAttributeRegistry{base}
}

func (reg *webIdProviderAttributeRegistry) IdProviderAttribute(idpUuid, attrName string, caStmp *Stamp) (idpAttr interface{}, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.Get(idpUuid+"/"+attrName, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if value == nil {
		return nil, nil, nil
	}
	return value.(map[string]interface{})[attrName], newCaStmp, nil
}
