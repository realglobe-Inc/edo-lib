package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryEventRegistry(t *testing.T) {
	testEventRegistry(t, NewMemoryEventRegistry())
}
