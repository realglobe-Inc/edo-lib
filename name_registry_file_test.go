package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileNameRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	reg := NewFileNameRegistry(path, 0)
	if _, err := reg.(*nameRegistry).base.Put("names", testNameTree); err != nil {
		t.Fatal(err)
	}
	testNameRegistry(t, reg)
}
