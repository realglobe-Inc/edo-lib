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

	reg, err := NewMongoServiceExplorer(mongoAddr, testLabel, "service-explorer")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB(testLabel).DropDatabase()

	if err := reg.(*mongoDriver).DB(testLabel).C("service-explorer").Insert(bson.M{
		"service": bson.M{
			"uuid": testServUuid,
			"uri":  testUri,
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

	reg, err := NewMongoDatedServiceExplorer(mongoAddr, testLabel, "service-explorer", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*datedMongoDriver).DB(testLabel).DropDatabase()

	if err := reg.(*datedMongoDriver).DB(testLabel).C("service-explorer").Insert(bson.M{
		"service": bson.M{
			"uuid": testServUuid,
			"uri":  testUri,
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
