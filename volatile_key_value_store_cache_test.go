package driver

import (
	"testing"
)

func TestCachingVolatileKeyValueStore(t *testing.T) {
	testVolatileKeyValueStore(t, newCachingVolatileKeyValueStore(newMemoryVolatileKeyValueStore(0, 0)))
}

func TestCachingConcurrentVolatileKeyValueStore(t *testing.T) {
	testConcurrentVolatileKeyValueStore(t, newCachingVolatileKeyValueStore(newMemoryVolatileKeyValueStore(testStaleDur, testCaExpiDur)))
}

func TestCachingConcurrentVolatileKeyValueStoreConsistency(t *testing.T) {
	testConcurrentVolatileKeyValueStoreConsistency(t, newSynchronizedVolatileKeyValueStore(newCachingVolatileKeyValueStore(newMemoryVolatileKeyValueStore(testStaleDur, testCaExpiDur))))
}
