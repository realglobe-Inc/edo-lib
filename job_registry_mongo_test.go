package driver

import (
	"testing"
)

func TestMongoJobRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJobRegistry(mongoAddr, testLabel, "job-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*jobRegistry).base.(*mongoTimeLimitedKeyValueStore).base.base.DB(testLabel).DropDatabase()

	testJobRegistry(t, reg)
}

func TestMongoJobRegistryRemoveOld(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJobRegistry(mongoAddr, testLabel, "job-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*jobRegistry).base.(*mongoTimeLimitedKeyValueStore).base.base.DB(testLabel).DropDatabase()

	testJobRegistryRemoveOld(t, reg)
}
