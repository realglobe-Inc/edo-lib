package driver

import (
	"testing"
)

func TestCachingVolatileKeyValueStore(t *testing.T) {
	testVolatileKeyValueStore(t, newCachingVolatileKeyValueStore(newMemoryVolatileKeyValueStore(0, 0)))
}
