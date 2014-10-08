package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileJsRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testJsRegistry(t, NewFileJsRegistry(path, 0))
}

func TestFileJsRegistryStamp(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testJsRegistryStamp(t, NewFileJsRegistry(path, 0))
}
