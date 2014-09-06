package driver

import (
	"testing"
)

func TestCachingIdProviderLister(t *testing.T) {
	reg := NewMemoryDatedIdProviderLister(0)
	reg.AddIdProvider(&IdProvider{"a_b-c", "ABC", "https://localhost:1234"})
	testIdProviderLister(t, NewCachingIdProviderLister(reg))
}
