package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoJobRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJobRegistry(mongoAddr, testLabel, "job-registry")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB(testLabel).DropDatabase()

	testJobRegistry(t, reg)
}

func TestMongoJobRegistryRemoveOld(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJobRegistry(mongoAddr, testLabel, "job-registry")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB(testLabel).DropDatabase()

	testJobRegistryRemoveOld(t, reg)
}
