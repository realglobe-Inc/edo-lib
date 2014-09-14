package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryServiceExplorer(t *testing.T) {
	reg := NewMemoryServiceExplorer()
	reg.AddServiceUuid(testUri, testServUuid)
	testServiceExplorer(t, reg)
}

// キャッシュ用。
func TestMemoryDatedServiceExplorer(t *testing.T) {
	reg := NewMemoryDatedServiceExplorer(0)
	reg.AddServiceUuid(testUri, testServUuid)
	testDatedServiceExplorer(t, reg)
}
