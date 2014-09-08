package driver

import (
	"testing"
)

// キャッシュ用。
func TestCachingDatedServiceKeyRegistry(t *testing.T) {
	// ////////////////////////////////
	// hndl := util.InitLog("github.com/realglobe-Inc")
	// hndl.SetLevel(level.ALL)
	// defer hndl.SetLevel(level.INFO)
	// ////////////////////////////////

	reg := NewMemoryDatedServiceKeyRegistry(0)
	reg.AddServiceKey("a_b-c", "kore ga kagi dayo.")
	testDatedServiceKeyRegistry(t, NewCachingDatedServiceKeyRegistry(reg))
}
