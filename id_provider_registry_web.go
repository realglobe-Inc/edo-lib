package driver

import ()

// 非キャッシュ用。
func NewWebIdProviderRegistry(prefix string) IdProviderRegistry {
	return newIdProviderRegistry(newWebKeyValueStore(prefix))
}

// キャッシュ用。
func NewWebDatedIdProviderRegistry(prefix string) DatedIdProviderRegistry {
	// TODO キャッシュの並列化。
	return newDatedIdProviderRegistry(newSynchronizedDatedKeyValueStore(newCachingDatedKeyValueStore(newWebDatedKeyValueStore(prefix))))
}
