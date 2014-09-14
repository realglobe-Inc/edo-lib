package driver

import (
	"testing"
)

// キャッシュ用。
func TestCachingDatedServiceExplorer(t *testing.T) {
	// ////////////////////////////////
	// hndl := util.InitLog("github.com/realglobe-Inc")
	// hndl.SetLevel(level.ALL)
	// defer hndl.SetLevel(level.INFO)
	// ////////////////////////////////

	reg := NewMemoryDatedServiceExplorer(0)
	reg.AddServiceUuid(testUri, testServUuid)
	testDatedServiceExplorer(t, newCachingDatedServiceExplorer(reg))
}
