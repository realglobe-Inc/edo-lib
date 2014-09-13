package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryLoginRegistry(t *testing.T) {
	reg := NewMemoryLoginRegistry()
	reg.AddUser("abc-012", "a_b-c")
	testLoginRegistry(t, reg)
}
