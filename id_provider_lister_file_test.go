package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// 非キャッシュ用。
func TestFileIdProviderLister(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "idp.json"), []*IdProvider{&IdProvider{"a_b-c", "ABC", "https://localhost:1234"}}); err != nil {
		t.Fatal(err)
	}

	testIdProviderLister(t, NewFileIdProviderLister(path))
}

// キャッシュ用。
func TestFileDatedIdProviderLister(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "idp.json"), []*IdProvider{&IdProvider{"a_b-c", "ABC", "https://localhost:1234"}}); err != nil {
		t.Fatal(err)
	}

	testDatedIdProviderLister(t, NewFileDatedIdProviderLister(path, 0))
}
