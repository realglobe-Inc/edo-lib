package driver

import (
	"testing"
)

func TestMongoJsBackendRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJsBackendRegistry(mongoAddr, "test_driver_mongo", "js", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoBackend).DB("test_driver_mongo").DropDatabase()

	testJsBackendRegistry(t, reg)
}
