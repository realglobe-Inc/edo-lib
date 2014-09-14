package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoNameRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoNameRegistry(mongoAddr, testLabel, "name-registry")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB(testLabel).DropDatabase()

	for name, addr := range testNameAddrMap {
		if err := reg.(*mongoDriver).DB(testLabel).C("name-registry").Insert(&mongoAddress{name, addr}); err != nil {
			t.Fatal(err)
		}
	}

	testNameRegistry(t, reg)
}
