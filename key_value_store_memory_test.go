package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryKeyValueStore(t *testing.T) {
	testKeyValueStore(t, newMemoryKeyValueStore())
}

// キャッシュ用。
func TestMemoryDatedKeyValueStore(t *testing.T) {
	testDatedKeyValueStore(t, newMemoryDatedKeyValueStore(0))
}
