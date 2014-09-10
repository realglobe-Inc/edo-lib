package driver

import (
	"testing"
)

// キャッシュ用。
func TestCachingDatedKeyValueStore(t *testing.T) {
	// ////////////////////////////////
	// hndl := util.InitLog("github.com/realglobe-Inc")
	// hndl.SetLevel(level.ALL)
	// defer hndl.SetLevel(level.INFO)
	// ////////////////////////////////

	testDatedKeyValueStore(t, newCachingDatedKeyValueStore(newMemoryDatedKeyValueStore(0)))
}
