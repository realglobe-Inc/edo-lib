package driver

import ()

func NewWebUserAttributeRegistry(prefix string) UserAttributeRegistry {
	return newUserAttributeRegistry(newWebKeyValueStore(prefix))
}
