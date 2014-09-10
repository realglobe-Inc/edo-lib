package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

// 非キャッシュ用。
func TestFileKeyValueStore(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testKeyValueStore(t, newFileKeyValueStore(path))
}

// キャッシュ用。
func TestFileDatedKeyValueStore(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testDatedKeyValueStore(t, newFileDatedKeyValueStore(path, 0))
}
