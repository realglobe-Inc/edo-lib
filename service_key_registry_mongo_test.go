package driver

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

// 非キャッシュ用。
func TestMongoServiceKeyRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoServiceKeyRegistry(mongoAddr, "test_driver_mongo", "key")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoRegistry).DB("test_driver_mongo").C("key").Insert(bson.M{
		"service": bson.M{
			"uuid":       "a_b-c",
			"public_key": "kore ga kagi dayo.",
		},
	}); err != nil {
		t.Fatal(err)
	}

	testServiceKeyRegistry(t, reg)
}

// キャッシュ用。
func TestMongoDatedServiceKeyRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoDatedServiceKeyRegistry(mongoAddr, "test_driver_mongo", "key", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoBackend).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoBackend).DB("test_driver_mongo").C("key").Insert(bson.M{
		"service": bson.M{
			"uuid":       "a_b-c",
			"public_key": "kore ga kagi dayo.",
		},
		"stamp": &Stamp{
			Date:   time.Now(),
			Digest: "0",
		},
	}); err != nil {
		t.Fatal(err)
	}

	testDatedServiceKeyRegistry(t, reg)
}
