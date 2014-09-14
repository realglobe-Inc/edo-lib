package driver

import (
	"testing"
)

// キャッシュ用。
func TestCachingDatedIdProviderLister(t *testing.T) {
	// ////////////////////////////////
	// hndl := util.InitLog("github.com/realglobe-Inc")
	// hndl.SetLevel(level.ALL)
	// defer hndl.SetLevel(level.INFO)
	// ////////////////////////////////

	reg := NewMemoryDatedIdProviderLister(0)
	reg.SetIdProviders(testIdps)
	testDatedIdProviderLister(t, newCachingDatedIdProviderLister(reg))
}
