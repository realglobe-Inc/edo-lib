package driver

import (
	"encoding/json"
	"time"
)

// スレッドセーフ。
func NewFileUserNameIndex(path string, expiDur time.Duration) UserNameIndex {
	return newUserNameIndex(NewFileKeyValueStore(path, jsonKeyGen, json.Marshal, jsonUnmarshal, expiDur))
}
