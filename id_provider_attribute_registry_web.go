package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// {
//   "id_provider: {
//     attrNameX: XXX
//   }
// }

// 非キャッシュ用。
func NewWebIdProviderAttributeRegistry(prefix string) IdProviderAttributeRegistry {
	return newWebIdProviderAttributeRegistry(newWebKeyValueStore(prefix))
}

type webIdProviderAttributeRegistry struct {
	keyValueStore
}

func newWebIdProviderAttributeRegistry(base keyValueStore) *webIdProviderAttributeRegistry {
	return &webIdProviderAttributeRegistry{base}
}

func (reg *webIdProviderAttributeRegistry) IdProviderAttribute(idpUuid, attrName string) (idpAttr interface{}, err error) {
	val, err := reg.get(idProviderAttributeKey(idpUuid, attrName))
	if err != nil {
		return nil, erro.Wrap(err)
	} else if val == nil {
		return nil, nil
	}
	return val.(map[string]interface{})["id_provider"].(map[string]interface{})[attrName], nil
}

// キャッシュ用。
func NewWebDatedIdProviderAttributeRegistry(prefix string) DatedIdProviderAttributeRegistry {
	// TODO キャッシュの並列化。
	return newWebDatedIdProviderAttributeRegistry(newSynchronizedDatedKeyValueStore(newCachingDatedKeyValueStore(newWebDatedKeyValueStore(prefix))))
}

type webDatedIdProviderAttributeRegistry struct {
	datedKeyValueStore
}

func newWebDatedIdProviderAttributeRegistry(base datedKeyValueStore) *webDatedIdProviderAttributeRegistry {
	return &webDatedIdProviderAttributeRegistry{base}
}

func (reg *webDatedIdProviderAttributeRegistry) StampedIdProviderAttribute(idpUuid, attrName string, caStmp *Stamp) (idpAttr interface{}, newCaStmp *Stamp, err error) {
	val, newCaStmp, err := reg.stampedGet(idProviderAttributeKey(idpUuid, attrName), caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if newCaStmp == nil {
		return nil, nil, nil
	}
	return val.(map[string]interface{})["id_provider"].(map[string]interface{})[attrName], newCaStmp, nil
}
