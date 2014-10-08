package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileLoginRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	reg := NewFileLoginRegistry(path, 0)
	if _, err := reg.(*loginRegistry).base.Put(testAccToken, testUsrName); err != nil {
		t.Fatal(err)
	}

	testLoginRegistry(t, reg)
}
