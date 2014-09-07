package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryIdProviderRegistry(t *testing.T) {
	reg := NewMemoryIdProviderRegistry()
	reg.AddIdProviderQueryUri("a_b-c", "https://localhost:1234/query")
	testIdProviderRegistry(t, reg)
}

// キャッシュ用。
func TestMemoryDatedIdProviderRegistry(t *testing.T) {
	reg := NewMemoryDatedIdProviderRegistry(0)
	reg.AddIdProviderQueryUri("a_b-c", "https://localhost:1234/query")
	testDatedIdProviderRegistry(t, reg)
}
