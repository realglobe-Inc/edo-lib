package driver

import ()

// 非キャッシュ用。
func NewWebServiceKeyRegistry(prefix string) ServiceKeyRegistry {
	return newServiceKeyRegistry(newWebKeyValueStore(prefix))
}

// キャッシュ用。
func NewWebDatedServiceKeyRegistry(prefix string) DatedServiceKeyRegistry {
	// TODO キャッシュの並列化。
	return newDatedServiceKeyRegistry(newSynchronizedDatedKeyValueStore(newCachingDatedKeyValueStore(newWebDatedKeyValueStore(prefix))))
}
