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

func _TestMongoNameRegistry(t *testing.T) {
	reg, err := NewMongoNameRegistry(mongoAddr, "test_driver_mongo", "name")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoRegistry).DB("test_driver_mongo").C("name").Insert(
		&mongoAddress{"c.b.a", "c.localhost"},
		&mongoAddress{"d.b.a", "d.localhost"},
		&mongoAddress{"b.a", "localhost"},
	); err != nil {
		t.Fatal(err)
	}

	testNameRegistry(t, reg)
}
