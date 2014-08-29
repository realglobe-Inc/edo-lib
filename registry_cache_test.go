package driver

import (
	"testing"
)

func TestCachingIdProviderRegistry(t *testing.T) {
	reg := NewMemoryIdProviderBackend()
	reg.AddIdProvider(&IdProvider{"a_b-c", "ABC", "https://localhost:1234"})
	testIdProviderRegistry(t, NewCachingIdProviderRegistry(reg))
}
