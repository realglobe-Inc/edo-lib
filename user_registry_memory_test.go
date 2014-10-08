package driver

import (
	"testing"
)

func TestMemoryUserRegistry(t *testing.T) {
	testUserRegistry(t, NewMemoryUserRegistry(0))
}
