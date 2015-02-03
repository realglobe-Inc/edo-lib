package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"io/ioutil"
	"os"
	"testing"
)

func init() {
	util.SetupConsoleLog("github.com/realglobe-Inc", level.OFF)
}

func TestFileVolatileKeyValueStore(t *testing.T) {
	// ////////////////////////////////
	// util.SetupConsoleLog("github.com/realglobe-Inc", level.ALL)
	// defer util.SetupConsoleLog("github.com/realglobe-Inc", level.OFF)
	// ////////////////////////////////
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)
	expiPath, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(expiPath)

	testVolatileKeyValueStore(t, newFileConcurrentVolatileKeyValueStore(path, expiPath, nil, nil, json.Marshal, jsonUnmarshal, testStaleDur, testCaExpiDur))
}

func TestFileConcurrentVolatileKeyValueStore(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)
	expiPath, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(expiPath)

	testConcurrentVolatileKeyValueStore(t, newFileConcurrentVolatileKeyValueStore(path, expiPath, nil, nil, json.Marshal, jsonUnmarshal, testStaleDur, testCaExpiDur))
}

func TestFileConcurrentVolatileKeyValueStoreConsistency(t *testing.T) {
	// ////////////////////////////////
	// util.SetupConsoleLog("github.com/realglobe-Inc", level.ALL)
	// defer util.SetupConsoleLog("github.com/realglobe-Inc", level.OFF)
	// ////////////////////////////////
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)
	expiPath, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(expiPath)

	testConcurrentVolatileKeyValueStoreConsistency(t, NewFileConcurrentVolatileKeyValueStore(path, expiPath, nil, nil, json.Marshal, jsonUnmarshal, testStaleDur, testCaExpiDur))
}
