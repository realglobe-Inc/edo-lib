package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

// 非キャッシュ用。
func TestFileJobRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testJobRegistry(t, NewFileJobRegistry(path))
}
