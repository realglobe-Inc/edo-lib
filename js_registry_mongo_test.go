package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoJsRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJsRegistry(mongoAddr, "test_driver_mongo", "js")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB("test_driver_mongo").DropDatabase()

	testJsRegistry(t, reg)
}

// キャッシュ用。
func TestMongoJsBackendRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJsBackendRegistry(mongoAddr, "test_driver_mongo", "js", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*datedMongoDriver).DB("test_driver_mongo").DropDatabase()

	testJsBackendRegistry(t, reg)
}
