package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryUserAttributeRegistry(t *testing.T) {
	reg := NewMemoryUserAttributeRegistry()
	reg.AddUserAttribute(testUsrUuid, testAttrName, testAttr)
	testUserAttributeRegistry(t, reg)
}
