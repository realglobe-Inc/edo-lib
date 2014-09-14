package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// 非キャッシュ用。
func TestFileUserNameIndex(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, testUsrName+".json"), testUsrUuid); err != nil {
		t.Fatal(err)
	}

	testUserNameIndex(t, NewFileUserNameIndex(path))
}
