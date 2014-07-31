package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileJsBackendRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testJsBackendRegistry(t, NewFileJsBackendRegistry(path))
}
