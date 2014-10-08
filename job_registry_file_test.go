package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileJobRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testJobRegistry(t, NewFileJobRegistry(path, 0))
}
