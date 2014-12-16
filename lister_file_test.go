package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFileLister(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)
	if err := ioutil.WriteFile(filepath.Join(path, testKey), []byte{}, filePerm); err != nil {
		t.Fatal(err)
	}

	testLister(t, newFileLister(path, nil, 0, 0))
}
