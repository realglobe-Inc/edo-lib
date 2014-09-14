package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// 非キャッシュ用。
func TestFileLoginRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, testAccToken+".json"), testUsrName); err != nil {
		t.Fatal(err)
	}

	testLoginRegistry(t, NewFileLoginRegistry(path))
}
