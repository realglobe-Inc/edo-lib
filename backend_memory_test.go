package driver

import (
	"testing"
)

func TestMemoryJsBackendRegistry(t *testing.T) {
	testJsBackendRegistry(t, NewMemoryJsBackendRegistry(0))
}
