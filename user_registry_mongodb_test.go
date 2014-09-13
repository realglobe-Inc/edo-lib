package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoUserRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoUserRegistry(mongoAddr, "test_driver_mongo", "user")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB("test_driver_mongo").DropDatabase()

	testUserRegistry(t, reg)
}
