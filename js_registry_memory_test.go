package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryJsRegistry(t *testing.T) {
	testJsRegistry(t, NewMemoryJsRegistry())
}

// キャッシュ用。
func TestMemoryJsBackendRegistry(t *testing.T) {
	testJsBackendRegistry(t, NewMemoryJsBackendRegistry(0))
}
