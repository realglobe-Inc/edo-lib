package driver

import (
	"testing"
)

func TestMongoEventRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoEventRegistry(mongoAddr, testLabel, "event-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*eventRegistry).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	testEventRegistry(t, reg)
}
