package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestSynchronizedJsRegistry(t *testing.T) {
	testJsRegistry(t, NewSynchronizedJsRegistry(NewMemoryJsRegistry()))
}

// キャッシュ用。
func TestSynchronizedJsBackendRegistry(t *testing.T) {
	testJsBackendRegistry(t, NewSynchronizedJsBackendRegistry(NewMemoryJsBackendRegistry(0)))
}
