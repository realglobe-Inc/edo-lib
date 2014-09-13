package driver

import ()

// 非キャッシュ用。
func NewFileUserAttributeRegistry(path string) UserAttributeRegistry {
	return newUserAttributeRegistry(newSynchronizedKeyValueStore(newFileKeyValueStore(path)))
}
