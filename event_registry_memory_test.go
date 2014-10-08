package driver

import (
	"testing"
)

func TestMemoryEventRegistry(t *testing.T) {
	testEventRegistry(t, NewMemoryEventRegistry(0))
}
