package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoUserNameIndex(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoUserNameIndex(mongoAddr, testLabel, "user-name-index")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*userNameIndex).keyValueStore.(*mongoKeyValueStore).DB(testLabel).DropDatabase()

	if err := reg.(*userNameIndex).put(testUsrName, testUsrUuid); err != nil {
		t.Fatal(err)
	}

	testUserNameIndex(t, reg)
}
