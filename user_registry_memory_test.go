package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryUserRegistry(t *testing.T) {
	testUserRegistry(t, NewMemoryUserRegistry())
}
