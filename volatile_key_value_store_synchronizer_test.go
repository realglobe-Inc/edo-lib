package driver

import (
	"testing"
)

func TestWebVolatileKeyValueStore(t *testing.T) {
	testVolatileKeyValueStore(t, newSynchronizedVolatileKeyValueStore(newMemoryVolatileKeyValueStore(testStaleDur, testCaExpiDur)))
}

func TestSynchronizedConcurrentVolatileKeyValueStore(t *testing.T) {
	testConcurrentVolatileKeyValueStore(t, newSynchronizedVolatileKeyValueStore(newMemoryVolatileKeyValueStore(testStaleDur, testCaExpiDur)))
}

func TestSynchronizedConcurrentVolatileKeyValueStoreConsistency(t *testing.T) {
	testConcurrentVolatileKeyValueStoreConsistency(t, newSynchronizedVolatileKeyValueStore(newMemoryVolatileKeyValueStore(testStaleDur, testCaExpiDur)))
}
