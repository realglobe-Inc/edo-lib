package driver

import (
	"testing"
)

func TestSynchronizedVolatileKeyValueStore(t *testing.T) {
	testVolatileKeyValueStore(t, newSynchronizedConcurrentVolatileKeyValueStore(newMemoryConcurrentVolatileKeyValueStore(testStaleDur, testCaExpiDur)))
}

func TestSynchronizedConcurrentVolatileKeyValueStore(t *testing.T) {
	testConcurrentVolatileKeyValueStore(t, newSynchronizedConcurrentVolatileKeyValueStore(newMemoryConcurrentVolatileKeyValueStore(testStaleDur, testCaExpiDur)))
}

func TestSynchronizedConcurrentVolatileKeyValueStoreConsistency(t *testing.T) {
	testConcurrentVolatileKeyValueStoreConsistency(t, newSynchronizedConcurrentVolatileKeyValueStore(newMemoryConcurrentVolatileKeyValueStore(testStaleDur, testCaExpiDur)))
}
