package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryServiceKeyRegistry(t *testing.T) {
	reg := NewMemoryServiceKeyRegistry()
	reg.AddServiceKey(testServUuid, testPublicKey)
	testServiceKeyRegistry(t, reg)
}

// キャッシュ用。
func TestMemoryDatedServiceKeyRegistry(t *testing.T) {
	reg := NewMemoryDatedServiceKeyRegistry(0)
	reg.AddServiceKey(testServUuid, testPublicKey)
	testDatedServiceKeyRegistry(t, reg)
}
