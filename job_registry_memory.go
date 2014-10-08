package driver

import (
	"time"
)

// スレッドセーフ。
func NewMemoryJobRegistry(expiDur time.Duration) JobRegistry {
	return newJobRegistry(NewMemoryTimeLimitedKeyValueStore(expiDur))
}
