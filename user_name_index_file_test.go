package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFileUserNameIndex(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "a_b-c.json"), "aaaa-bbbb-cccc"); err != nil {
		t.Fatal(err)
	}

	testUserNameIndex(t, NewFileUserNameIndex(path))
}
