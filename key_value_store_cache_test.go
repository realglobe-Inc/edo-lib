package driver

import (
	"testing"
)

func TestCachingKeyValueStoreStamp(t *testing.T) {
	// ////////////////////////////////
	// hndl := util.InitLog("github.com/realglobe-Inc")
	// hndl.SetLevel(level.ALL)
	// defer hndl.SetLevel(level.INFO)
	// ////////////////////////////////

	testKeyValueStoreStamp(t, newCachingKeyValueStore(newMemoryKeyValueStore(0)))
}
