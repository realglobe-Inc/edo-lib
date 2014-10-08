package driver

import (
	"testing"
)

func TestSynchronizedKeyValueStore(t *testing.T) {
	testKeyValueStore(t, newSynchronizedKeyValueStore(NewMemoryKeyValueStore(0)))
}

func TestSynchronizedKeyValueStoreStamp(t *testing.T) {
	testKeyValueStoreStamp(t, newSynchronizedKeyValueStore(NewMemoryKeyValueStore(0)))
}
