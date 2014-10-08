package driver

import (
	"encoding/json"
	"time"
)

// スレッドセーフ。
func NewFileUserRegistry(path string, expiDur time.Duration) UserRegistry {
	return newUserRegistry(NewFileKeyValueStore(path, jsonKeyGen, json.Marshal, jsonUnmarshal, expiDur))
}
