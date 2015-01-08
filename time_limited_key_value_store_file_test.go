package driver

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileTimeLimitedKeyValueStore(t *testing.T) {
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

	testTimeLimitedKeyValueStore(t, newFileTimeLimitedKeyValueStore(path, expiPath, nil, nil, json.Marshal, jsonUnmarshal, 0, 0))
}
