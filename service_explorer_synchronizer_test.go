package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestSynchronizedServiceExplorer(t *testing.T) {
	reg := NewMemoryServiceExplorer()
	reg.AddServiceUuid("https://localhost:1234/api", "a_b-c")
	testServiceExplorer(t, NewSynchronizedServiceExplorer(reg))
}

// キャッシュ用。
func TestSynchronizedDatedServiceExplorer(t *testing.T) {
	reg := NewMemoryDatedServiceExplorer(0)
	reg.AddServiceUuid("https://localhost:1234/api", "a_b-c")
	testDatedServiceExplorer(t, NewSynchronizedDatedServiceExplorer(reg))
}
