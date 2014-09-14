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

	reg, err := NewMongoIdProviderLister(mongoAddr, testLabel, "id-provider-lister")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB(testLabel).DropDatabase()

	for _, idp := range testIdps {
		if err := reg.(*mongoDriver).DB(testLabel).C("id-provider-lister").Insert(bson.M{"id_provider": idp}); err != nil {
			t.Fatal(err)
		}
	}

	testIdProviderLister(t, reg)
}

// キャッシュ用。
func TestMongoDatedIdProviderLister(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoDatedIdProviderLister(mongoAddr, testLabel, "id-provider-lister", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*datedMongoDriver).DB(testLabel).DropDatabase()

	if err := reg.(*datedMongoDriver).DB(testLabel).C("id-provider-lister").Insert(
		bson.M{
			"stamp": &Stamp{
				Date:   time.Now(),
				Digest: "0",
			},
		},
	); err != nil {
		t.Fatal(err)
	}
	for _, idp := range testIdps {
		if err := reg.(*datedMongoDriver).DB(testLabel).C("id-provider-lister").Insert(bson.M{"id_provider": idp}); err != nil {
			t.Fatal(err)
		}
	}

	testDatedIdProviderLister(t, reg)
}
