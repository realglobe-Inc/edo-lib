package driver

import (
	"testing"
)

func TestMemoryListedKeyValueStore(t *testing.T) {
	testListedKeyValueStore(t, newMemoryListedKeyValueStore(0, 0))
}

func TestMemoryListedKeyValueStoreStamp(t *testing.T) {
	testListedKeyValueStoreStamp(t, newMemoryListedKeyValueStore(0, 0))
}
