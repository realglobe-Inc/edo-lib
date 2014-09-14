package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// 非キャッシュ用。
func TestFileIdProviderLister(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "list.json"), testIdps); err != nil {
		t.Fatal(err)
	}

	testIdProviderLister(t, NewFileIdProviderLister(path))
}

// キャッシュ用。
func TestFileDatedIdProviderLister(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "list.json"), testIdps); err != nil {
		t.Fatal(err)
	}

	testDatedIdProviderLister(t, NewFileDatedIdProviderLister(path, 0))
}
