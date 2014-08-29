package driver

import (
	"testing"
)

func TestCachingIdProviderBackend(t *testing.T) {
	// ////////////////////////////////
	// hndl := util.InitLog("github.com/realglobe-Inc")
	// hndl.SetLevel(level.ALL)
	// defer hndl.SetLevel(level.INFO)
	// ////////////////////////////////

	reg := NewMemoryIdProviderBackend(0)
	reg.AddIdProvider(&IdProvider{
		IdpUuid: "a_b-c",
		Name:    "ABC",
		Uri:     "https://localhost:1234",
	})
	testIdProviderBackend(t, NewCachingIdProviderBackend(reg))
}
