package driver

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

// 非キャッシュ用。
func TestMongoServiceExplorer(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoServiceExplorer(mongoAddr, "test_driver_mongo", "idp")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoDriver).DB("test_driver_mongo").C("idp").Insert(bson.M{
		"service": bson.M{
			"uuid": "a_b-c",
			"uri":  "https://localhost:1234/api",
		},
	}); err != nil {
		t.Fatal(err)
	}

	testServiceExplorer(t, reg)
}

// キャッシュ用。
func TestMongoDatedServiceExplorer(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoDatedServiceExplorer(mongoAddr, "test_driver_mongo", "idp", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*datedMongoDriver).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*datedMongoDriver).DB("test_driver_mongo").C("idp").Insert(bson.M{
		"service": bson.M{
			"uuid": "a_b-c",
			"uri":  "https://localhost:1234/api",
		},
		"stamp": &Stamp{
			Date:   time.Now(),
			Digest: "0",
		},
	}); err != nil {
		t.Fatal(err)
	}

	testDatedServiceExplorer(t, reg)
}
