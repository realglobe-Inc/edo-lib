package driver

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

// 非キャッシュ用。
func TestMongoIdProviderRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdProviderRegistry(mongoAddr, "test_driver_mongo", "idp")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoRegistry).DB("test_driver_mongo").C("idp").Insert(bson.M{
		"id_provider": bson.M{
			"uuid":      "a_b-c",
			"query_uri": "https://localhost:1234/query",
		},
	}); err != nil {
		t.Fatal(err)
	}

	testIdProviderRegistry(t, reg)
}

// キャッシュ用。
func TestMongoDatedIdProviderRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoDatedIdProviderRegistry(mongoAddr, "test_driver_mongo", "idp", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoBackend).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoBackend).DB("test_driver_mongo").C("idp").Insert(bson.M{
		"id_provider": bson.M{
			"uuid":      "a_b-c",
			"query_uri": "https://localhost:1234/query",
		},
		"stamp": &Stamp{
			Date:   time.Now(),
			Digest: "0",
		},
	}); err != nil {
		t.Fatal(err)
	}

	testDatedIdProviderRegistry(t, reg)
}
