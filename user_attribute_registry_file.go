package driver

import (
	"encoding/json"
	"time"
)

// スレッドセーフ。
func NewFileUserAttributeRegistry(path string, expiDur time.Duration) UserAttributeRegistry {
	return newUserAttributeRegistry(NewFileKeyValueStore(path, jsonKeyGen, json.Marshal, jsonUnmarshal, expiDur))
}
