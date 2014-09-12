package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// 非キャッシュ用。
func TestFileServiceKeyRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "a_b-c.json"), testPublicKeyPem); err != nil {
		t.Fatal(err)
	}

	testServiceKeyRegistry(t, NewFileServiceKeyRegistry(path))
}

// キャッシュ用。
func TestFileDatedServiceKeyRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "a_b-c.json"), testPublicKeyPem); err != nil {
		t.Fatal(err)
	}

	testDatedServiceKeyRegistry(t, NewFileDatedServiceKeyRegistry(path, 0))
}
