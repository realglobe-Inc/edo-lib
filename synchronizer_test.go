package driver

import (
	"testing"
)

func TestSynchronizedJsRegistry(t *testing.T) {
	testJsRegistry(t, NewSynchronizedJsRegistry(NewMemoryJsRegistry()))
}

func TestSynchronizedUserRegistry(t *testing.T) {
	testUserRegistry(t, NewSynchronizedUserRegistry(NewMemoryUserRegistry()))
}

func TestSynchronizedJobRegistry(t *testing.T) {
	testJobRegistry(t, NewSynchronizedJobRegistry(NewMemoryJobRegistry()))
}

func TestSynchronizedEventRegistry(t *testing.T) {
	testEventRegistry(t, NewSynchronizedEventRegistry(NewMemoryEventRegistry()))
}
