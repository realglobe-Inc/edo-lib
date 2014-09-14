package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryNameRegistry(t *testing.T) {
	reg := NewMemoryNameRegistry()
	for name, addr := range testNameAddrMap {
		reg.AddAddress(name, addr)
	}
	testNameRegistry(t, reg)
}
