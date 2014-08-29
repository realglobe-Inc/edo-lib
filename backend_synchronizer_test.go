package driver

import (
	"testing"
)

func TestSynchronizedJsBackendRegistry(t *testing.T) {
	testJsBackendRegistry(t, NewSynchronizedJsBackendRegistry(NewMemoryJsBackendRegistry()))
}

func TestSynchronizedIdProviderBackend(t *testing.T) {
	reg := NewMemoryIdProviderBackend()
	reg.AddIdProvider(&IdProvider{
		IdpUuid: "a_b-c",
		Name:    "ABC",
		Uri:     "https://localhost:1234",
	})
	testIdProviderBackend(t, NewSynchronizedIdProviderBackend(reg))
}
