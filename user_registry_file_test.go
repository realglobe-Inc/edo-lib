package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileUserRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testUserRegistry(t, NewFileUserRegistry(path, 0))
}
