package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestSynchronizedEventRegistry(t *testing.T) {
	testEventRegistry(t, NewSynchronizedEventRegistry(NewMemoryEventRegistry()))
}
