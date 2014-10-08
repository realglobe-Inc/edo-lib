package driver

import (
	"testing"
)

func TestSynchronizedRawDataStore(t *testing.T) {
	testRawDataStore(t, newSynchronizedRawDataStore(newMemoryRawDataStore(0)))
}

func TestSynchronizedRawDataStoreStamp(t *testing.T) {
	testRawDataStoreStamp(t, newSynchronizedRawDataStore(newMemoryRawDataStore(0)))
}
