package driver

import (
	"time"
)

// スレッドセーフ。
func NewMemoryEventRegistry(expiDur time.Duration) EventRegistry {
	return newEventRegistry(NewMemoryKeyValueStore(expiDur))
}
