package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoJobRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJobRegistry(mongoAddr, "test_driver_mongo", "job")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB("test_driver_mongo").DropDatabase()

	testJobRegistry(t, reg)
}

func TestMongoJobRegistryRemoveOld(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJobRegistry(mongoAddr, "test_driver_mongo", "job")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB("test_driver_mongo").DropDatabase()

	testJobRegistryRemoveOld(t, reg)
}
