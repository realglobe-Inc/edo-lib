package driver

import (
	"testing"
)

func TestCachingDatedIdProviderLister(t *testing.T) {
	// ////////////////////////////////
	// hndl := util.InitLog("github.com/realglobe-Inc")
	// hndl.SetLevel(level.ALL)
	// defer hndl.SetLevel(level.INFO)
	// ////////////////////////////////

	reg := NewMemoryDatedIdProviderLister(0)
	reg.AddIdProvider(&IdProvider{
		Uuid: "a_b-c",
		Name: "ABC",
		Uri:  "https://localhost:1234",
	})
	testDatedIdProviderLister(t, NewCachingDatedIdProviderLister(reg))
}
