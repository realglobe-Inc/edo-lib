package driver

import (
	"testing"
)

func TestMongoKeyValueStore(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg := newMongoKeyValueStore(mongoAddr, testLabel, "key-value-store", nil, nil, nil, 0, 0)
	defer reg.Clear()

	testKeyValueStore(t, reg)
}

func TestMongoKeyValueStoreStamp(t *testing.T) {
	// ////////////////////////////////
	// util.SetupConsoleLog("github.com/realglobe-Inc", level.ALL)
	// defer util.SetupConsoleLog("github.com/realglobe-Inc", level.OFF)
	// ////////////////////////////////
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg := newMongoKeyValueStore(mongoAddr, testLabel, "key-value-store", nil, nil, nil, 0, 0)
	defer reg.Clear()

	testKeyValueStoreStamp(t, reg)
}
