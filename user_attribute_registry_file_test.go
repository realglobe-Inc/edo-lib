package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// 非キャッシュ用。
func TestFileUserAttributeRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, escapeToFileName(userAttributeKey(testUsrUuid, testAttrName))+".json"), testAttr); err != nil {
		t.Fatal(err)
	}

	testUserAttributeRegistry(t, NewFileUserAttributeRegistry(path))
}
