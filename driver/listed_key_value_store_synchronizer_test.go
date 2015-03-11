package driver

import (
	"testing"
)

func TestSynchronizedListedKeyValueStore(t *testing.T) {
	testListedKeyValueStore(t, newSynchronizedListedKeyValueStore(NewMemoryListedKeyValueStore(0, 0)))
}

func TestSynchronizedListedKeyValueStoreStamp(t *testing.T) {
	testListedKeyValueStoreStamp(t, newSynchronizedListedKeyValueStore(NewMemoryListedKeyValueStore(0, 0)))
}
