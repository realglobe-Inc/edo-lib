package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoLoginRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoLoginRegistry(mongoAddr, testLabel, "login-registry")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB(testLabel).DropDatabase()

	if err := reg.(*mongoDriver).DB(testLabel).C("login-registry").Insert(
		&mongoUser{testAccToken, testUsrName},
	); err != nil {
		t.Fatal(err)
	}

	testLoginRegistry(t, reg)
}
