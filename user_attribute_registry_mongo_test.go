package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoUserAttributeRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoUserAttributeRegistry(mongoAddr, testLabel, "user-attribute-registry")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*userAttributeRegistry).keyValueStore.(*mongoKeyValueStore).DB(testLabel).DropDatabase()

	if err := reg.(*userAttributeRegistry).put(userAttributeKey(testUsrUuid, testAttrName), testAttr); err != nil {
		t.Fatal(err)
	}

	testUserAttributeRegistry(t, reg)
}
