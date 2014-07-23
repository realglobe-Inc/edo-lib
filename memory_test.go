package driver

import (
	"testing"
)

func TestMemoryJsRegistry(t *testing.T) {
	testJsRegistry(t, NewMemoryJsRegistry())
}

func TestMemoryUserRegistry(t *testing.T) {
	testUserRegistry(t, NewMemoryUserRegistry())
}

func TestMemoryJobRegistry(t *testing.T) {
	testJobRegistry(t, NewMemoryJobRegistry())
}

func TestMemoryNameRegistry(t *testing.T) {
	reg := NewMemoryNameRegistry()
	reg.AddAddress("c.b.a", "c.localhost")
	reg.AddAddress("d.b.a", "d.localhost")
	reg.AddAddress("b.a", "localhost")
	testNameRegistry(t, reg)
}

func TestMemoryEventRegistry(t *testing.T) {
	testEventRegistry(t, NewMemoryEventRegistry())
}
