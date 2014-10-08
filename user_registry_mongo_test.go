package driver

import (
	"testing"
)

func TestMongoUserRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoUserRegistry(mongoAddr, testLabel, "user-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*userRegistry).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	testUserRegistry(t, reg)
}
