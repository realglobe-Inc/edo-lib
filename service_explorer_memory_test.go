package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryServiceExplorer(t *testing.T) {
	reg := NewMemoryServiceExplorer()
	reg.AddServiceUuid("https://localhost:1234/api", "a_b-c")
	testServiceExplorer(t, reg)
}

// キャッシュ用。
func TestMemoryDatedServiceExplorer(t *testing.T) {
	reg := NewMemoryDatedServiceExplorer(0)
	reg.AddServiceUuid("https://localhost:1234/api", "a_b-c")
	testDatedServiceExplorer(t, reg)
}
