package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestSynchronizedServiceKeyRegistry(t *testing.T) {
	reg := NewMemoryServiceKeyRegistry()
	reg.AddServiceKey("a_b-c", "kore ga kagi dayo.")
	testServiceKeyRegistry(t, NewSynchronizedServiceKeyRegistry(reg))
}

// キャッシュ用。
func TestSynchronizedDatedServiceKeyRegistry(t *testing.T) {
	reg := NewMemoryDatedServiceKeyRegistry(0)
	reg.AddServiceKey("a_b-c", "kore ga kagi dayo.")
	testDatedServiceKeyRegistry(t, NewSynchronizedDatedServiceKeyRegistry(reg))
}
