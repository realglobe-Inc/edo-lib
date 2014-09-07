package driver

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

// 非キャッシュ用。
func TestMongoIdProviderLister(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdProviderLister(mongoAddr, "test_driver_mongo", "idp")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoRegistry).DB("test_driver_mongo").C("idp").Insert(bson.M{
		"id_provider": &IdProvider{
			Uuid:     "a_b-c",
			Name:     "ABC",
			LoginUri: "https://localhost:1234",
		},
	}); err != nil {
		t.Fatal(err)
	}

	testIdProviderLister(t, reg)
}

// キャッシュ用。
func TestMongoDatedIdProviderLister(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoDatedIdProviderLister(mongoAddr, "test_driver_mongo", "idp", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoBackend).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoBackend).DB("test_driver_mongo").C("idp").Insert(
		bson.M{
			"id_provider": &IdProvider{
				Uuid:     "a_b-c",
				Name:     "ABC",
				LoginUri: "https://localhost:1234",
			},
		},
		bson.M{
			"stamp": &Stamp{
				Date:   time.Now(),
				Digest: "0",
			},
		},
	); err != nil {
		t.Fatal(err)
	}

	testDatedIdProviderLister(t, reg)
}
