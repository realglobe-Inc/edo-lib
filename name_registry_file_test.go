package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// 非キャッシュ用。
func TestFileNameRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	for name, addr := range testNameAddrMap {
		if err := writeToJson(filepath.Join(path, name+".json"), addr); err != nil {
			t.Fatal(err)
		}
	}
	testNameRegistry(t, NewFileNameRegistry(path))
}
