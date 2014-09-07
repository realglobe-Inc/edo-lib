package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestSynchronizedIdProviderRegistry(t *testing.T) {
	reg := NewMemoryIdProviderRegistry()
	reg.AddIdProviderQueryUri("a_b-c", "https://localhost:1234/query")
	testIdProviderRegistry(t, NewSynchronizedIdProviderRegistry(reg))
}

// キャッシュ用。
func TestSynchronizedDatedIdProviderRegistry(t *testing.T) {
	reg := NewMemoryDatedIdProviderRegistry(0)
	reg.AddIdProviderQueryUri("a_b-c", "https://localhost:1234/query")
	testDatedIdProviderRegistry(t, NewSynchronizedDatedIdProviderRegistry(reg))
}
