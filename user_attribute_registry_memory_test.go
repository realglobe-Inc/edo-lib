package driver

import (
	"testing"
)

func TestMemoryUserAttributeRegistry(t *testing.T) {
	reg := NewMemoryUserAttributeRegistry()
	reg.AddUserAttribute("a_b-c", "attribute", "abcd")
	testUserAttributeRegistry(t, reg)
}
