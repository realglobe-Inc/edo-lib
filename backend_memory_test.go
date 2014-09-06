package driver

import (
	"testing"
)

func TestMemoryJsBackendRegistry(t *testing.T) {
	testJsBackendRegistry(t, NewMemoryJsBackendRegistry(0))
}

func TestMemoryDatedIdProviderLister(t *testing.T) {
	reg := NewMemoryDatedIdProviderLister(0)
	reg.AddIdProvider(&IdProvider{"a_b-c", "ABC", "https://localhost:1234"})
	testDatedIdProviderLister(t, reg)
}
