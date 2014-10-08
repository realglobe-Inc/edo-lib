package driver

import (
	"testing"
)

func TestMemoryServiceExplorer(t *testing.T) {
	reg := NewMemoryServiceExplorer(0)
	reg.SetServiceUuids(map[string]string{testUri: testServUuid})
	testServiceExplorer(t, reg)
}

func TestMemoryServiceExplorerStamp(t *testing.T) {
	reg := NewMemoryServiceExplorer(0)
	reg.SetServiceUuids(map[string]string{testUri: testServUuid})
	testServiceExplorerStamp(t, reg)
}
