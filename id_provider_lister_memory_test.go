package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryIdProviderLister(t *testing.T) {
	reg := NewMemoryIdProviderLister()
	reg.SetIdProviders(testIdps)
	testIdProviderLister(t, reg)
}

// キャッシュ用。
func TestMemoryDatedIdProviderLister(t *testing.T) {
	reg := NewMemoryDatedIdProviderLister(0)
	reg.SetIdProviders(testIdps)
	testDatedIdProviderLister(t, reg)
}
