package driver

import (
	"testing"
)

func TestMongoLoginRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoLoginRegistry(mongoAddr, testLabel, "login-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*loginRegistry).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*loginRegistry).base.Put(testAccToken, testUsrName); err != nil {
		t.Fatal(err)
	}

	testLoginRegistry(t, reg)
}
