package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestSynchronizedIdProviderLister(t *testing.T) {
	reg := NewMemoryIdProviderLister()
	reg.AddIdProvider(&IdProvider{
		Uuid:     "a_b-c",
		Name:     "ABC",
		LoginUri: "https://localhost:1234",
	})
	testIdProviderLister(t, NewSynchronizedIdProviderLister(reg))
}

// キャッシュ用。
func TestSynchronizedDatedIdProviderLister(t *testing.T) {
	reg := NewMemoryDatedIdProviderLister(0)
	reg.AddIdProvider(&IdProvider{
		Uuid:     "a_b-c",
		Name:     "ABC",
		LoginUri: "https://localhost:1234",
	})
	testDatedIdProviderLister(t, NewSynchronizedDatedIdProviderLister(reg))
}
