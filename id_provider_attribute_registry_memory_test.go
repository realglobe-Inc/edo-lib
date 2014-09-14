package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryIdProviderAttributeRegistry(t *testing.T) {
	reg := NewMemoryIdProviderAttributeRegistry()
	reg.AddIdProviderAttribute(testIdpUuid, testAttrName, testAttr)
	testIdProviderAttributeRegistry(t, reg)
}

// キャッシュ用。
func TestMemoryDatedIdProviderAttributeRegistry(t *testing.T) {
	reg := NewMemoryDatedIdProviderAttributeRegistry(0)
	reg.AddIdProviderAttribute(testIdpUuid, testAttrName, testAttr)
	testDatedIdProviderAttributeRegistry(t, reg)
}
