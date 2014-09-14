package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestSynchronizedIdProviderLister(t *testing.T) {
	reg := NewMemoryIdProviderLister()
	reg.SetIdProviders(testIdps)
	testIdProviderLister(t, newSynchronizedIdProviderLister(reg))
}

// キャッシュ用。
func TestSynchronizedDatedIdProviderLister(t *testing.T) {
	reg := NewMemoryDatedIdProviderLister(0)
	reg.SetIdProviders(testIdps)
	testDatedIdProviderLister(t, newSynchronizedDatedIdProviderLister(reg))
}
