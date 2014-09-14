package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// 非キャッシュ用。
func TestFileIdProviderAttributeRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, escapeToFileName(userAttributeKey(testIdpUuid, testAttrName))+".json"), testAttr); err != nil {
		t.Fatal(err)
	}

	testIdProviderAttributeRegistry(t, NewFileIdProviderAttributeRegistry(path))
}

// キャッシュ用。
func TestFileDatedIdProviderAttributeRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, escapeToFileName(userAttributeKey(testIdpUuid, testAttrName))+".json"), testAttr); err != nil {
		t.Fatal(err)
	}

	testDatedIdProviderAttributeRegistry(t, NewFileDatedIdProviderAttributeRegistry(path, 0))
}
