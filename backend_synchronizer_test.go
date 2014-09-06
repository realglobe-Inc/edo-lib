package driver

import (
	"testing"
)

func TestSynchronizedJsBackendRegistry(t *testing.T) {
	testJsBackendRegistry(t, NewSynchronizedJsBackendRegistry(NewMemoryJsBackendRegistry(0)))
}

func TestSynchronizedDatedIdProviderLister(t *testing.T) {
	reg := NewMemoryDatedIdProviderLister(0)
	reg.AddIdProvider(&IdProvider{
		Uuid: "a_b-c",
		Name: "ABC",
		Uri:  "https://localhost:1234",
	})
	testDatedIdProviderLister(t, NewSynchronizedDatedIdProviderLister(reg))
}
