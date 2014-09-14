package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoUserRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoUserRegistry(mongoAddr, testLabel, "user-registry")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB(testLabel).DropDatabase()

	testUserRegistry(t, reg)
}
