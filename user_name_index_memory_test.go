package driver

import (
	"testing"
)

func TestMemoryUserNameIndex(t *testing.T) {
	reg := NewMemoryUserNameIndex()
	reg.AddUserUuid("a_b-c", "aaaa-bbbb-cccc")
	testUserNameIndex(t, reg)
}
