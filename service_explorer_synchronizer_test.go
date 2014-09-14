package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestSynchronizedServiceExplorer(t *testing.T) {
	reg := NewMemoryServiceExplorer()
	reg.AddServiceUuid(testUri, testServUuid)
	testServiceExplorer(t, newSynchronizedServiceExplorer(reg))
}

// キャッシュ用。
func TestSynchronizedDatedServiceExplorer(t *testing.T) {
	reg := NewMemoryDatedServiceExplorer(0)
	reg.AddServiceUuid(testUri, testServUuid)
	testDatedServiceExplorer(t, newSynchronizedDatedServiceExplorer(reg))
}
