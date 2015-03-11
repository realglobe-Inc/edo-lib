package driver

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileListedKeyValueStore(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testListedKeyValueStore(t, newFileListedKeyValueStore(path, nil, nil, json.Marshal, jsonUnmarshal, 0, 0))
}

func TestFileListedKeyValueStoreStamp(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testListedKeyValueStoreStamp(t, newFileListedKeyValueStore(path, nil, nil, json.Marshal, jsonUnmarshal, 0, 0))
}
