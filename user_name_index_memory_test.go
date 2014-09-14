package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryUserNameIndex(t *testing.T) {
	reg := NewMemoryUserNameIndex()
	reg.AddUserUuid(testUsrName, testUsrUuid)
	testUserNameIndex(t, reg)
}
