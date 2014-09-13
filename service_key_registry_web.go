package driver

import (
	"crypto/rsa"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// {
//   "service": {
//     "public_key": "XXXXX"
//   }
// }

// 非キャッシュ用。
func NewWebServiceKeyRegistry(prefix string) ServiceKeyRegistry {
	return newWebServiceKeyRegistry(newWebKeyValueStore(prefix))
}

type webServiceKeyRegistry struct {
	keyValueStore
}

func newWebServiceKeyRegistry(base keyValueStore) ServiceKeyRegistry {
	return &webServiceKeyRegistry{base}
}

func (reg *webServiceKeyRegistry) ServiceKey(servUuid string) (servKey *rsa.PublicKey, err error) {
	val, err := reg.get(servUuid)
	if err != nil {
		return nil, erro.Wrap(err)
	} else if val == nil || val == "" {
		return nil, nil
	}
	return parseKey(val.(map[string]interface{})["service"].(map[string]interface{})["public_key"].(string))
}

// キャッシュ用。
func NewWebDatedServiceKeyRegistry(prefix string) DatedServiceKeyRegistry {
	// TODO キャッシュの並列化。
	return newWebDatedServiceKeyRegistry(newSynchronizedDatedKeyValueStore(newCachingDatedKeyValueStore(newWebDatedKeyValueStore(prefix))))
}

type webDatedServiceKeyRegistry struct {
	datedKeyValueStore
}

func newWebDatedServiceKeyRegistry(base datedKeyValueStore) DatedServiceKeyRegistry {
	return &webDatedServiceKeyRegistry{base}
}

func (reg *webDatedServiceKeyRegistry) StampedServiceKey(servUuid string, caStmp *Stamp) (servKey *rsa.PublicKey, newCaStmp *Stamp, err error) {
	val, newCaStmp, err := reg.stampedGet(servUuid, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if newCaStmp == nil {
		return nil, nil, nil
	} else if val == nil || val == "" {
		return nil, newCaStmp, nil
	}

	servKey, err = parseKey(val.(map[string]interface{})["service"].(map[string]interface{})["public_key"].(string))
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}
	return servKey, newCaStmp, nil
}
