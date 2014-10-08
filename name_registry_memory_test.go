package driver

import (
	"testing"
)

func TestMemoryNameRegistry(t *testing.T) {
	reg := NewMemoryNameRegistry(0)
	reg.SetAddresses(testNameTree.toContainer())
	testNameRegistry(t, reg)
}
