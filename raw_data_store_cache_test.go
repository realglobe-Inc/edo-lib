package driver

import (
	"testing"
)

func TestCachingRawDataStore(t *testing.T) {
	testRawDataStore(t, newCachingRawDataStore(newMemoryRawDataStore(0)))
}

func TestCachingRawDataStoreStamp(t *testing.T) {
	testRawDataStoreStamp(t, newCachingRawDataStore(newMemoryRawDataStore(0)))
}
