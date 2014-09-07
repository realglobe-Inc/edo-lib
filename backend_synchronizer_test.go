package driver

import (
	"testing"
)

func TestSynchronizedJsBackendRegistry(t *testing.T) {
	testJsBackendRegistry(t, NewSynchronizedJsBackendRegistry(NewMemoryJsBackendRegistry(0)))
}
