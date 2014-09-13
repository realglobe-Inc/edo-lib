package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryUserAttributeRegistry(t *testing.T) {
	reg := NewMemoryUserAttributeRegistry()
	reg.AddUserAttribute("a_b-c", "attribute", "abcd")
	testUserAttributeRegistry(t, reg)
}
