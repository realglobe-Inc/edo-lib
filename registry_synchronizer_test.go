package driver

import (
	"testing"
)

func TestSynchronizedLoginRegistry(t *testing.T) {
	reg := NewMemoryLoginRegistry()
	reg.AddUser("abc-012", "a_b-c")
	testLoginRegistry(t, NewSynchronizedLoginRegistry(reg))
}

func TestSynchronizedJsRegistry(t *testing.T) {
	testJsRegistry(t, NewSynchronizedJsRegistry(NewMemoryJsRegistry()))
}

func TestSynchronizedUserRegistry(t *testing.T) {
	testUserRegistry(t, NewSynchronizedUserRegistry(NewMemoryUserRegistry()))
}

func TestSynchronizedJobRegistry(t *testing.T) {
	testJobRegistry(t, NewSynchronizedJobRegistry(NewMemoryJobRegistry()))
}

func TestSynchronizedNameRegistry(t *testing.T) {
	reg := NewMemoryNameRegistry()
	reg.AddAddress("c.b.a", "c.localhost")
	reg.AddAddress("d.b.a", "d.localhost")
	reg.AddAddress("b.a", "localhost")
	testNameRegistry(t, NewSynchronizedNameRegistry(reg))
}

func TestSynchronizedEventRegistry(t *testing.T) {
	testEventRegistry(t, NewSynchronizedEventRegistry(NewMemoryEventRegistry()))
}

func TestSynchronizedServiceRegistry(t *testing.T) {
	reg := NewMemoryServiceRegistry()
	reg.AddService("localhost:1234", "a_b-c")
	testServiceRegistry(t, NewSynchronizedServiceRegistry(reg))
}

func TestSynchronizedIdProviderLister(t *testing.T) {
	reg := NewMemoryIdProviderLister()
	reg.AddIdProvider(&IdProvider{
		Uuid: "a_b-c",
		Name: "ABC",
		Uri:  "https://localhost:1234",
	})
	testIdProviderLister(t, NewSynchronizedIdProviderLister(reg))
}
