package driver

import (
	"testing"
)

func TestCachingKeyValueStoreStamp(t *testing.T) {
	// ////////////////////////////////
	// util.SetupConsoleLog("github.com/realglobe-Inc", level.ALL)
	// defer util.SetupConsoleLog("github.com/realglobe-Inc", level.OFF)
	// ////////////////////////////////

	testKeyValueStoreStamp(t, newCachingKeyValueStore(newMemoryKeyValueStore(0)))
}
