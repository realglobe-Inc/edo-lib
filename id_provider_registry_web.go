package driver

import ()

// {
//   "id_provider: {
//     "query_uri": "https://aaa.bbb.ccc/query"
//   }
// }

// 非キャッシュ用。
func NewWebIdProviderRegistry(prefix string) IdProviderRegistry {
	return newWebIdProviderRegistry(newWebKeyValueStore(prefix))
}

type webIdProviderRegistry struct {
	keyValueStore
}

func newWebIdProviderRegistry(base keyValueStore) *webIdProviderRegistry {
	return &webIdProviderRegistry{base}
}

func (reg *webIdProviderRegistry) IdProviderQueryUri(idpUuid string) (queryUri string, err error) {
	val, err := reg.get(idpUuid)
	if val != nil {
		queryUri = val.(map[string]interface{})["id_provider"].(map[string]interface{})["query_uri"].(string)
	}
	return queryUri, err
}

// キャッシュ用。
func NewWebDatedIdProviderRegistry(prefix string) DatedIdProviderRegistry {
	// TODO キャッシュの並列化。
	return newWebDatedIdProviderRegistry(newSynchronizedDatedKeyValueStore(newCachingDatedKeyValueStore(newWebDatedKeyValueStore(prefix))))
}

type webDatedIdProviderRegistry struct {
	datedKeyValueStore
}

func newWebDatedIdProviderRegistry(base datedKeyValueStore) *webDatedIdProviderRegistry {
	return &webDatedIdProviderRegistry{base}
}

func (reg *webDatedIdProviderRegistry) StampedIdProviderQueryUri(idpUuid string, caStmp *Stamp) (queryUri string, newCaStmp *Stamp, err error) {
	val, newCaStmp, err := reg.stampedGet(idpUuid, caStmp)
	if val != nil {
		queryUri = val.(map[string]interface{})["id_provider"].(map[string]interface{})["query_uri"].(string)
	}
	return queryUri, newCaStmp, err
}
