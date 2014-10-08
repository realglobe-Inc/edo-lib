package driver

import (
	"testing"
)

func TestMongoNameRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoNameRegistry(mongoAddr, testLabel, "name-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*nameRegistry).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*nameRegistry).base.Put("names", testNameTree); err != nil {
		t.Fatal(err)
	}

	testNameRegistry(t, reg)
}
