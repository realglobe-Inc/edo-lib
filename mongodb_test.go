package driver

import (
	"testing"
)

var mongoAddr string = "localhost"

func _TestMongoJobRegistry(t *testing.T) {
	reg, err := NewMongoJobRegistry(mongoAddr, "test_driver_mongo", "job")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	testJobRegistry(t, reg)
}

func _TestMongoJobRegistryRemoveOld(t *testing.T) {
	reg, err := NewMongoJobRegistry(mongoAddr, "test_driver_mongo", "job")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	testJobRegistryRemoveOld(t, reg)
}
