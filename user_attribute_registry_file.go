package driver

import ()

func NewFileUserAttributeRegistry(path string) UserAttributeRegistry {
	return newUserAttributeRegistry(newSynchronizedKeyValueStore(newFileKeyValueStore(path)))
}
