package driver

import (
	"testing"
)

func TestMongoServiceKeyRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoServiceKeyRegistry(mongoAddr, testLabel, "service-key-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*serviceKeyRegistry).base.(*cachingKeyValueStore).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*serviceKeyRegistry).base.Put(testServUuid, testPublicKey); err != nil {
		t.Fatal(err)
	}

	testServiceKeyRegistry(t, reg)
}

func TestMongoServiceKeyRegistryStamp(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoServiceKeyRegistry(mongoAddr, testLabel, "service-key-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*serviceKeyRegistry).base.(*cachingKeyValueStore).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*serviceKeyRegistry).base.Put(testServUuid, testPublicKey); err != nil {
		t.Fatal(err)
	}

	testServiceKeyRegistryStamp(t, reg)
}
