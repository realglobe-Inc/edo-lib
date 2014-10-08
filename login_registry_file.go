package driver

import (
	"encoding/json"
	"time"
)

// スレッドセーフ。
func NewFileLoginRegistry(path string, expiDur time.Duration) LoginRegistry {
	return newLoginRegistry(NewFileKeyValueStore(path, jsonKeyGen, json.Marshal, jsonUnmarshal, expiDur))
}
