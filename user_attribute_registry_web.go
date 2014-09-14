package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// {
//   "user": {
//     attrNameX: XXX,
//   }
// }

// 非キャッシュ用。
func NewWebUserAttributeRegistry(prefix string) UserAttributeRegistry {
	return newWebUserAttributeRegistry(newWebKeyValueStore(prefix))
}

type webUserAttributeRegistry struct {
	keyValueStore
}

func newWebUserAttributeRegistry(base keyValueStore) *webUserAttributeRegistry {
	return &webUserAttributeRegistry{base}
}

func (reg *webUserAttributeRegistry) UserAttribute(usrUuid, attrName string) (usrAttr interface{}, err error) {
	val, err := reg.get(userAttributeKey(usrUuid, attrName))
	if err != nil {
		return nil, erro.Wrap(err)
	} else if val == nil {
		return nil, nil
	}
	return val.(map[string]interface{})["user"].(map[string]interface{})[attrName], nil
}
