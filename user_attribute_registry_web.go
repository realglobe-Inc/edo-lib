package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// {
//   "user": {
//     attrNameX: XXX,
//   }
// }

// スレッドセーフ。
func NewWebUserAttributeRegistry(prefix string) UserAttributeRegistry {
	return &webUserAttributeRegistry{NewWebKeyValueStore(prefix, nil, func(data []byte) (interface{}, error) {
		var res struct {
			User map[string]interface{}
		}
		if err := json.Unmarshal(data, &res); err != nil {
			return nil, erro.Wrap(err)
		}
		return res.User, nil
	})}
}

type webUserAttributeRegistry struct {
	base KeyValueStore
}

func (reg *webUserAttributeRegistry) UserAttribute(usrUuid, attrName string, caStmp *Stamp) (usrAttr interface{}, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get(usrUuid+"/"+attrName, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if value == nil {
		return nil, newCaStmp, nil
	}
	return value.(map[string]interface{})[attrName], newCaStmp, nil
}
