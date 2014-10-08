package driver

import (
	"testing"
)

func TestMemoryJobRegistry(t *testing.T) {
	testJobRegistry(t, NewMemoryJobRegistry(0))
}
