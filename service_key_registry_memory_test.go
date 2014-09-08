package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryServiceKeyRegistry(t *testing.T) {
	reg := NewMemoryServiceKeyRegistry()
	reg.AddServiceKey("a_b-c", "kore ga kagi dayo.")
	testServiceKeyRegistry(t, reg)
}

// キャッシュ用。
func TestMemoryDatedServiceKeyRegistry(t *testing.T) {
	reg := NewMemoryDatedServiceKeyRegistry(0)
	reg.AddServiceKey("a_b-c", "kore ga kagi dayo.")
	testDatedServiceKeyRegistry(t, reg)
}
