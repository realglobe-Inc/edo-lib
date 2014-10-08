package driver

import (
	"testing"
)

func TestMongoJsRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJsRegistry(mongoAddr, testLabel, "js-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*jsRegistry).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	testJsRegistry(t, reg)
}

func TestMongoJsRegistryStamp(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJsRegistry(mongoAddr, testLabel, "js-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*jsRegistry).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	testJsRegistryStamp(t, reg)
}
