package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

// キャッシュ用。
func TestFileJsRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testJsRegistry(t, NewFileJsRegistry(path))
}

// 非キャッシュ用。
func TestFileJsBackendRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testJsBackendRegistry(t, NewFileJsBackendRegistry(path, 0))
}
