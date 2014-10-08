package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileRawDataStore(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testRawDataStore(t, newFileRawDataStore(path, nil, 0))
}

func TestFileRawDataStoreStamp(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testRawDataStoreStamp(t, newFileRawDataStore(path, nil, 0))
}
