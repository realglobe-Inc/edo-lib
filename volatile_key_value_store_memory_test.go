package driver

import (
	"testing"
)

func TestMemoryVolatileKeyValueStore(t *testing.T) {
	testVolatileKeyValueStore(t, newMemoryVolatileKeyValueStore(0, 0))
}
