package driver

import (
	"testing"
)

func TestSynchronizedListedRawDataStore(t *testing.T) {
	testListedRawDataStore(t, newSynchronizedListedRawDataStore(newMemoryListedRawDataStore(0, 0)))
}

func TestSynchronizedListedRawDataStoreStamp(t *testing.T) {
	testListedRawDataStoreStamp(t, newSynchronizedListedRawDataStore(newMemoryListedRawDataStore(0, 0)))
}
