package driver

import (
	"testing"
)

func TestMemoryServiceKeyRegistry(t *testing.T) {
	reg := NewMemoryServiceKeyRegistry(0)
	reg.AddServiceKey(testServUuid, testPublicKey)
	testServiceKeyRegistry(t, reg)
}

func TestMemoryServiceKeyRegistryStamp(t *testing.T) {
	reg := NewMemoryServiceKeyRegistry(0)
	reg.AddServiceKey(testServUuid, testPublicKey)
	testServiceKeyRegistry(t, reg)
}
