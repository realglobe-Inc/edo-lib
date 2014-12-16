package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileListedRawDataStore(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testListedRawDataStore(t, newFileListedRawDataStore(path, nil, nil, 0, 0))
}

func TestFileListedRawDataStoreStamp(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testListedRawDataStoreStamp(t, newFileListedRawDataStore(path, nil, nil, 0, 0))
}
