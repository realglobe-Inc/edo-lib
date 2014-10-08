package driver

import (
	"testing"
)

func TestMemoryIdProviderAttributeRegistry(t *testing.T) {
	reg := NewMemoryIdProviderAttributeRegistry(0)
	reg.AddIdProviderAttribute(testIdpUuid, testAttrName, testAttr)
	testIdProviderAttributeRegistry(t, reg)
}

func TestMemoryIdProviderAttributeRegistryStamp(t *testing.T) {
	reg := NewMemoryIdProviderAttributeRegistry(0)
	reg.AddIdProviderAttribute(testIdpUuid, testAttrName, testAttr)
	testIdProviderAttributeRegistryStamp(t, reg)
}
