package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMemoryTimeLimitedKeyValueStore(t *testing.T) {
	testTimeLimitedKeyValueStore(t, NewMemoryTimeLimitedKeyValueStore())
}
