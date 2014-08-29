package driver

import (
	"testing"
)

func TestMemoryJsBackendRegistry(t *testing.T) {
	testJsBackendRegistry(t, NewMemoryJsBackendRegistry())
}

func TestMemoryIdProviderBackend(t *testing.T) {
	reg := NewMemoryIdProviderBackend()
	reg.AddIdProvider(&IdProvider{"a_b-c", "ABC", "https://localhost:1234"})
	testIdProviderBackend(t, reg)
}
