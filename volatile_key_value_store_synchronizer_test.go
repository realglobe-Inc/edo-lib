package driver

import (
	"testing"
)

func TestWebVolatileKeyValueStore(t *testing.T) {
	testVolatileKeyValueStore(t, newSynchronizedVolatileKeyValueStore(newMemoryConcurrentVolatileKeyValueStore(testStaleDur, testCaExpiDur)))
}

func TestSynchronizedConcurrentVolatileKeyValueStore(t *testing.T) {
	testConcurrentVolatileKeyValueStore(t, newSynchronizedVolatileKeyValueStore(newMemoryConcurrentVolatileKeyValueStore(testStaleDur, testCaExpiDur)))
}

func TestSynchronizedConcurrentVolatileKeyValueStoreConsistency(t *testing.T) {
	testConcurrentVolatileKeyValueStoreConsistency(t, newSynchronizedVolatileKeyValueStore(newMemoryConcurrentVolatileKeyValueStore(testStaleDur, testCaExpiDur)))
}
