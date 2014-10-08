package driver

import (
	"testing"
)

func TestMemoryJsRegistry(t *testing.T) {
	testJsRegistry(t, NewMemoryJsRegistry(0))
}

func TestMemoryJsRegistryStamp(t *testing.T) {
	testJsRegistryStamp(t, NewMemoryJsRegistry(0))
}
