package driver

import (
	"testing"
)

func TestMemoryVolatileKeyValueStore(t *testing.T) {
	testVolatileKeyValueStore(t, newMemoryVolatileKeyValueStore(testStaleDur, testCaExpiDur))
}

func TestMemoryConcurrentVolatileKeyValueStore(t *testing.T) {
	testConcurrentVolatileKeyValueStore(t, newMemoryVolatileKeyValueStore(testStaleDur, testCaExpiDur))
}

func TestMemoryConcurrentVolatileKeyValueStoreConsistency(t *testing.T) {
	testConcurrentVolatileKeyValueStoreConsistency(t, NewMemoryVolatileKeyValueStore(testStaleDur, testCaExpiDur))
}
