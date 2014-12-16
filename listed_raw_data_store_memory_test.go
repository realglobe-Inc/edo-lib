package driver

import (
	"testing"
)

func TestMemoryListedRawDataStore(t *testing.T) {
	testListedRawDataStore(t, newMemoryListedRawDataStore(0, 0))
}

func TestMemoryListedRawDataStoreStamp(t *testing.T) {
	testListedRawDataStoreStamp(t, newMemoryListedRawDataStore(0, 0))
}
