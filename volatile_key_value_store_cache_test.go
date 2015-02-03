package driver

import (
	"testing"
)

func TestCachingVolatileKeyValueStore(t *testing.T) {
	testVolatileKeyValueStore(t, newCachingVolatileKeyValueStore(newMemoryConcurrentVolatileKeyValueStore(0, 0)))
}

func TestCachingConcurrentVolatileKeyValueStore(t *testing.T) {
	testConcurrentVolatileKeyValueStore(t, newCachingVolatileKeyValueStore(newMemoryConcurrentVolatileKeyValueStore(testStaleDur, testCaExpiDur)))
}

func TestCachingConcurrentVolatileKeyValueStoreConsistency(t *testing.T) {
	testConcurrentVolatileKeyValueStoreConsistency(t, newSynchronizedVolatileKeyValueStore(newCachingVolatileKeyValueStore(newMemoryConcurrentVolatileKeyValueStore(testStaleDur, testCaExpiDur))))
}
