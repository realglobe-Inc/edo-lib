package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// 非キャッシュ用。
func TestFileUserAttributeRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, escapeToFileName(userAttributeKey("a_b-c", "attribute"))+".json"), "abcd"); err != nil {
		t.Fatal(err)
	}

	testUserAttributeRegistry(t, NewFileUserAttributeRegistry(path))
}
