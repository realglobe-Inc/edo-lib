package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestSynchronizedKeyValueStore(t *testing.T) {
	testKeyValueStore(t, newSynchronizedKeyValueStore(newMemoryKeyValueStore()))
}

// キャッシュ用。
func TestSynchronizedDatedKeyValueStore(t *testing.T) {
	testDatedKeyValueStore(t, newSynchronizedDatedKeyValueStore(newMemoryDatedKeyValueStore(0)))
}
