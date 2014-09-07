package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// 非キャッシュ用。
func TestFileIdProviderRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "a_b-c.json"), "https://localhost:1234/query"); err != nil {
		t.Fatal(err)
	}

	testIdProviderRegistry(t, NewFileIdProviderRegistry(path))
}

// キャッシュ用。
func TestFileDatedIdProviderRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "a_b-c.json"), "https://localhost:1234/query"); err != nil {
		t.Fatal(err)
	}

	testDatedIdProviderRegistry(t, NewFileDatedIdProviderRegistry(path, 0))
}
