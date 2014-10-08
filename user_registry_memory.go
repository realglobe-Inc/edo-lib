package driver

import (
	"time"
)

// スレッドセーフ。
func NewMemoryUserRegistry(expiDur time.Duration) UserRegistry {
	return newUserRegistry(NewMemoryKeyValueStore(expiDur))
}
