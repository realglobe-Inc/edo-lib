package driver

import (
	"time"
)

// スレッドセーフ。
func NewMemoryJsRegistry(expiDur time.Duration) JsRegistry {
	return newJsRegistry(NewMemoryKeyValueStore(expiDur))
}
