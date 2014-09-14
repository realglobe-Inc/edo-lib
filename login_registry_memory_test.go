package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryLoginRegistry(t *testing.T) {
	reg := NewMemoryLoginRegistry()
	reg.AddUser(testAccToken, testUsrName)
	testLoginRegistry(t, reg)
}
