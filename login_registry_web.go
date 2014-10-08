package driver

import ()

// スレッドセーフ。
func NewWebLoginRegistry(prefix string) LoginRegistry {
	return newLoginRegistry(NewWebKeyValueStore(prefix, stringMarshal, stringUnmarshal))
}
