package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// 非キャッシュ用。
func TestFileServiceExplorer(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "list.json"), map[string]string{"https://localhost:1234/api": "a_b-c"}); err != nil {
		t.Fatal(err)
	}

	testServiceExplorer(t, NewFileServiceExplorer(path))
}

// キャッシュ用。
func TestFileDatedServiceExplorer(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "list.json"), map[string]string{"https://localhost:1234/api": "a_b-c"}); err != nil {
		t.Fatal(err)
	}

	testDatedServiceExplorer(t, NewFileDatedServiceExplorer(path, 0))
}
