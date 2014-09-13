package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoNameRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoNameRegistry(mongoAddr, "test_driver_mongo", "name")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoDriver).DB("test_driver_mongo").C("name").Insert(
		&mongoAddress{"c.b.a", "c.localhost"},
		&mongoAddress{"d.b.a", "d.localhost"},
		&mongoAddress{"b.a", "localhost"},
	); err != nil {
		t.Fatal(err)
	}

	testNameRegistry(t, reg)
}
