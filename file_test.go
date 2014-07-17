package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileJsRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_registry_file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testJsRegistry(t, NewFileJsRegistry(path))
}

func TestFileUserRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_registry_file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testUserRegistry(t, NewFileUserRegistry(path))
}

func TestFileEventRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_registry_file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testEventRegistry(t, NewFileEventRegistry(path))
}
