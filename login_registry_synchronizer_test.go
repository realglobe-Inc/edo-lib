package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestSynchronizedLoginRegistry(t *testing.T) {
	reg := NewMemoryLoginRegistry()
	reg.AddUser("abc-012", "a_b-c")
	testLoginRegistry(t, NewSynchronizedLoginRegistry(reg))
}
