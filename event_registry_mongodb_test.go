package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoEventRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoEventRegistry(mongoAddr, "test_driver_mongo", "event")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB("test_driver_mongo").DropDatabase()

	testEventRegistry(t, reg)
}
