package driver

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestMongoTaExplorer(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoTaExplorer(mongoAddr, testLabel, "ta-explorer", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*taExplorer).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*taExplorer).base.Put("list", bson.M{testUri: testServUuid}); err != nil {
		t.Fatal(err)
	}

	testTaExplorer(t, reg)
}

func TestMongoTaExplorerStamp(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoTaExplorer(mongoAddr, testLabel, "ta-explorer", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*taExplorer).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*taExplorer).base.Put("list", bson.M{testUri: testServUuid}); err != nil {
		t.Fatal(err)
	}

	testTaExplorerStamp(t, reg)
}
