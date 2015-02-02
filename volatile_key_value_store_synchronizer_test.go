package driver

import (
	"testing"
)

func TestWebVolatileKeyValueStore(t *testing.T) {
	testVolatileKeyValueStore(t, newSynchronizedVolatileKeyValueStore(newMemoryVolatileKeyValueStore(0, 0)))
}
