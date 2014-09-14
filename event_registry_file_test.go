package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

// 非キャッシュ用。
func TestFileEventRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testEventRegistry(t, NewFileEventRegistry(path))
}
