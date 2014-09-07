package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryIdProviderLister(t *testing.T) {
	reg := NewMemoryIdProviderLister()
	reg.AddIdProvider(&IdProvider{"a_b-c", "ABC", "https://localhost:1234"})
	testIdProviderLister(t, reg)
}

// キャッシュ用。
func TestMemoryDatedIdProviderLister(t *testing.T) {
	reg := NewMemoryDatedIdProviderLister(0)
	reg.AddIdProvider(&IdProvider{"a_b-c", "ABC", "https://localhost:1234"})
	testDatedIdProviderLister(t, reg)
}
