package driver

import (
	"testing"
)

func TestMemoryLoginRegistry(t *testing.T) {
	reg := NewMemoryLoginRegistry(0)
	reg.AddUser(testAccToken, testUsrName)
	testLoginRegistry(t, reg)
}
