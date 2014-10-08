package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileEventRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testEventRegistry(t, NewFileEventRegistry(path, 0))
}
