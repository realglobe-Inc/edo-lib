package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestSynchronizedLoginRegistry(t *testing.T) {
	reg := NewMemoryLoginRegistry()
	reg.AddUser(testAccToken, testUsrName)
	testLoginRegistry(t, NewSynchronizedLoginRegistry(reg))
}
