package driver

import (
	"testing"
)

func TestMemoryKeyValueStore(t *testing.T) {
	testKeyValueStore(t, newMemoryKeyValueStore(0))
}

func TestMemoryKeyValueStoreStamp(t *testing.T) {
	testKeyValueStoreStamp(t, newMemoryKeyValueStore(0))
}
