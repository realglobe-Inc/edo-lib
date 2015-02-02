package driver

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileVolatileKeyValueStore(t *testing.T) {
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

	testVolatileKeyValueStore(t, newFileVolatileKeyValueStore(path, expiPath, nil, nil, json.Marshal, jsonUnmarshal, 0, 0))
}
