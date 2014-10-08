package driver

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestMongoServiceExplorer(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoServiceExplorer(mongoAddr, testLabel, "service-explorer", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*serviceExplorer).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*serviceExplorer).base.Put("list", bson.M{testUri: testServUuid}); err != nil {
		t.Fatal(err)
	}

	testServiceExplorer(t, reg)
}

func TestMongoServiceExplorerStamp(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoServiceExplorer(mongoAddr, testLabel, "service-explorer", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*serviceExplorer).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*serviceExplorer).base.Put("list", bson.M{testUri: testServUuid}); err != nil {
		t.Fatal(err)
	}

	testServiceExplorerStamp(t, reg)
}
