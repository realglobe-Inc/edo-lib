package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoLoginRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoLoginRegistry(mongoAddr, "test_driver_mongo", "login")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoDriver).DB("test_driver_mongo").C("login").Insert(
		&mongoUser{"abc-012", "a_b-c"},
	); err != nil {
		t.Fatal(err)
	}

	testLoginRegistry(t, reg)
}
