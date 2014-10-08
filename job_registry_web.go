package driver

import (
	"encoding/json"
)

// スレッドセーフ。
func NewWebJobRegistry(prefix string) JobRegistry {
	return newJobRegistry(NewWebTimeLimitedKeyValueStore(prefix, json.Marshal, jobResultUnmarshal))
}
