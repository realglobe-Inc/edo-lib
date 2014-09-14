package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoServiceKeyRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoServiceKeyRegistry(mongoAddr, testLabel, "service-key-registry")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*serviceKeyRegistry).keyValueStore.(*mongoKeyValueStore).DB(testLabel).DropDatabase()

	if err := reg.(*serviceKeyRegistry).put(testServUuid, testPublicKeyPem); err != nil {
		t.Fatal(err)
	}

	testServiceKeyRegistry(t, reg)
}

// キャッシュ用。
func TestMongoDatedServiceKeyRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoDatedServiceKeyRegistry(mongoAddr, testLabel, "service-key-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*datedServiceKeyRegistry).datedKeyValueStore.(*mongoDatedKeyValueStore).DB(testLabel).DropDatabase()

	if _, err := reg.(*datedServiceKeyRegistry).stampedPut(testServUuid, testPublicKeyPem); err != nil {
		t.Fatal(err)
	}

	testDatedServiceKeyRegistry(t, reg)
}
