package driver

import (
	"testing"
)

func TestCachingTimeLimitedKeyValueStore(t *testing.T) {
	testTimeLimitedKeyValueStore(t, newCachingTimeLimitedKeyValueStore(newMemoryTimeLimitedKeyValueStore(0, 0)))
}
