package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestSynchronizedJobRegistry(t *testing.T) {
	testJobRegistry(t, NewSynchronizedJobRegistry(NewMemoryJobRegistry()))
}
