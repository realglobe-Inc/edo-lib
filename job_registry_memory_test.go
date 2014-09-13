package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryJobRegistry(t *testing.T) {
	testJobRegistry(t, NewMemoryJobRegistry())
}
