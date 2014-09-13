package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryNameRegistry(t *testing.T) {
	reg := NewMemoryNameRegistry()
	reg.AddAddress("c.b.a", "c.localhost")
	reg.AddAddress("d.b.a", "d.localhost")
	reg.AddAddress("b.a", "localhost")
	testNameRegistry(t, reg)
}
