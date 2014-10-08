package driver

import (
	"testing"
)

func TestMongoUserAttributeRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoUserAttributeRegistry(mongoAddr, testLabel, "user-attribute-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*userAttributeRegistry).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*userAttributeRegistry).base.Put(testUsrUuid+"/"+testAttrName, testAttr); err != nil {
		t.Fatal(err)
	}

	testUserAttributeRegistry(t, reg)
}
