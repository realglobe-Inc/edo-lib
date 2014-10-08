package driver

import (
	"encoding/json"
)

// スレッドセーフ。
func NewWebUserRegistry(prefix string) UserRegistry {
	return newUserRegistry(NewWebKeyValueStore(prefix, json.Marshal, jsonUnmarshal))
}
