package driver

import (
	"testing"
)

// キャッシュ用。
func TestCachingDatedIdProviderRegistry(t *testing.T) {
	// ////////////////////////////////
	// hndl := util.InitLog("github.com/realglobe-Inc")
	// hndl.SetLevel(level.ALL)
	// defer hndl.SetLevel(level.INFO)
	// ////////////////////////////////

	reg := NewMemoryDatedIdProviderRegistry(0)
	reg.AddIdProviderQueryUri("a_b-c", "https://localhost:1234/query")
	testDatedIdProviderRegistry(t, NewCachingDatedIdProviderRegistry(reg))
}
