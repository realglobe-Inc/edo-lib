package driver

import (
	"testing"
)

func TestMemoryRawDataStore(t *testing.T) {
	testRawDataStore(t, newMemoryRawDataStore(0, 0))
}

func TestMemoryRawDataStoreStamp(t *testing.T) {
	testRawDataStoreStamp(t, newMemoryRawDataStore(0, 0))
}
