package driver

import (
	"testing"
)

func TestMemoryIdProviderLister(t *testing.T) {
	reg := NewMemoryIdProviderLister(0)
	reg.SetIdProviders(testIdps)
	testIdProviderLister(t, reg)
}

func TestMemoryIdProviderListerStamp(t *testing.T) {
	reg := NewMemoryIdProviderLister(0)
	reg.SetIdProviders(testIdps)
	testIdProviderListerStamp(t, reg)
}
