package driver

import (
	"testing"
)

// テストするなら、ローカルにデフォルトポートで mongodb をたてる必要あり。

var mongoAddr string = "localhost"

func _TestMongoLoginRegistry(t *testing.T) {
	reg, err := NewMongoLoginRegistry(mongoAddr, "test_driver_mongo", "login")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoRegistry).DB("test_driver_mongo").C("login").Insert(
		&mongoUser{"abc-012", "a_b-c"},
	); err != nil {
		t.Fatal(err)
	}

	testLoginRegistry(t, reg)
}

func _TestMongoJsRegistry(t *testing.T) {
	reg, err := NewMongoJsRegistry(mongoAddr, "test_driver_mongo", "js")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	testJsRegistry(t, reg)
}

func _TestMongoUserRegistry(t *testing.T) {
	reg, err := NewMongoUserRegistry(mongoAddr, "test_driver_mongo", "user")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	testUserRegistry(t, reg)
}

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

func _TestMongoEventRegistry(t *testing.T) {
	reg, err := NewMongoEventRegistry(mongoAddr, "test_driver_mongo", "event")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	testEventRegistry(t, reg)
}

func _TestMongoServiceRegistry(t *testing.T) {
	reg, err := NewMongoServiceRegistry(mongoAddr, "test_driver_mongo", "service")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoRegistry).DB("test_driver_mongo").C("service").Insert(
		&mongoService{map[string]string{
			"localhost:1234": "a_b-c",
		}},
	); err != nil {
		t.Fatal(err)
	}

	testServiceRegistry(t, reg)
}
