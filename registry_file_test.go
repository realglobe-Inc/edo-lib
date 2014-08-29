package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFileLoginRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "abc-012.json"), "a_b-c"); err != nil {
		t.Fatal(err)
	}

	testLoginRegistry(t, NewFileLoginRegistry(path))
}

func TestFileJsRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testJsRegistry(t, NewFileJsRegistry(path))
}

func TestFileUserRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testUserRegistry(t, NewFileUserRegistry(path))
}

func TestFileJobRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testJobRegistry(t, NewFileJobRegistry(path))
}

func TestFileNameRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "c.b.a.json"), "c.localhost"); err != nil {
		t.Fatal(err)
	}
	if err := writeToJson(filepath.Join(path, "d.b.a.json"), "d.localhost"); err != nil {
		t.Fatal(err)
	}
	if err := writeToJson(filepath.Join(path, "b.a.json"), "localhost"); err != nil {
		t.Fatal(err)
	}

	testNameRegistry(t, NewFileNameRegistry(path))
}

func TestFileEventRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testEventRegistry(t, NewFileEventRegistry(path))
}

func TestFileServiceRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "uuid.json"), map[string]string{"localhost:1234": "a_b-c"}); err != nil {
		t.Fatal(err)
	}

	testServiceRegistry(t, NewFileServiceRegistry(path))
}

func TestFileIdProviderRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := writeToJson(filepath.Join(path, "idp.json"), []*IdProvider{&IdProvider{"a_b-c", "ABC", "https://localhost:1234"}}); err != nil {
		t.Fatal(err)
	}

	testIdProviderRegistry(t, NewFileIdProviderRegistry(path))
}
