package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// 非キャッシュ用。
func TestFileNameRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "c.b.a.json"), "c.localhost"); err != nil {
		t.Fatal(err)
	}
	if err := writeToJson(filepath.Join(path, "d.b.a.json"), "d.localhost"); err != nil {
		t.Fatal(err)
	}
	if err := writeToJson(filepath.Join(path, "b.a.json"), "localhost"); err != nil {
		t.Fatal(err)
	}

	testNameRegistry(t, NewFileNameRegistry(path))
}
