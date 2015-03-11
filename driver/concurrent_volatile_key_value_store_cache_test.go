package driver

import (
	"testing"
)

func TestCachingVolatileKeyValueStore(t *testing.T) {
	testVolatileKeyValueStore(t, newCachingConcurrentVolatileKeyValueStore(newMemoryConcurrentVolatileKeyValueStore(0, 0)))
}

func TestCachingConcurrentVolatileKeyValueStore(t *testing.T) {
	testConcurrentVolatileKeyValueStore(t, newCachingConcurrentVolatileKeyValueStore(newMemoryConcurrentVolatileKeyValueStore(testStaleDur, testCaExpiDur)))
}

func TestCachingConcurrentVolatileKeyValueStoreConsistency(t *testing.T) {
	testConcurrentVolatileKeyValueStoreConsistency(t, newSynchronizedConcurrentVolatileKeyValueStore(newCachingConcurrentVolatileKeyValueStore(newMemoryConcurrentVolatileKeyValueStore(testStaleDur, testCaExpiDur))))
}
