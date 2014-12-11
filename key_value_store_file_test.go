package driver

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileKeyValueStore(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testKeyValueStore(t, newFileKeyValueStore(path, nil, nil, json.Marshal, jsonUnmarshal, 0, 0))
}

func TestFileKeyValueStoreStamp(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testKeyValueStoreStamp(t, newFileKeyValueStore(path, nil, nil, json.Marshal, jsonUnmarshal, 0, 0))
}
