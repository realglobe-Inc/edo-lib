package driver

import (
	"testing"
)

func TestMemoryVolatileKeyValueStore(t *testing.T) {
	testVolatileKeyValueStore(t, newMemoryConcurrentVolatileKeyValueStore(testStaleDur, testCaExpiDur))
}

func TestMemoryConcurrentVolatileKeyValueStore(t *testing.T) {
	testConcurrentVolatileKeyValueStore(t, newMemoryConcurrentVolatileKeyValueStore(testStaleDur, testCaExpiDur))
}

func TestMemoryConcurrentVolatileKeyValueStoreConsistency(t *testing.T) {
	testConcurrentVolatileKeyValueStoreConsistency(t, NewMemoryConcurrentVolatileKeyValueStore(testStaleDur, testCaExpiDur))
}
