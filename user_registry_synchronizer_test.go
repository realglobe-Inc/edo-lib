package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestSynchronizedUserRegistry(t *testing.T) {
	testUserRegistry(t, NewSynchronizedUserRegistry(NewMemoryUserRegistry()))
}
