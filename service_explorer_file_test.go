package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileServiceExplorer(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	reg := NewFileServiceExplorer(path, 0)
	if _, err := reg.(*serviceExplorer).base.Put("list", testServExpTree); err != nil {
		t.Fatal(err)
	}

	testServiceExplorer(t, reg)
}

func TestFileServiceExplorerStamp(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	reg := NewFileServiceExplorer(path, 0)
	if _, err := reg.(*serviceExplorer).base.Put("list", testServExpTree); err != nil {
		t.Fatal(err)
	}

	testServiceExplorerStamp(t, reg)
}
