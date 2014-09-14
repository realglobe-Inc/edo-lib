package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoJsRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJsRegistry(mongoAddr, testLabel, "js-registry")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB(testLabel).DropDatabase()

	testJsRegistry(t, reg)
}

// キャッシュ用。
func TestMongoJsBackendRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJsBackendRegistry(mongoAddr, testLabel, "js-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*datedMongoDriver).DB(testLabel).DropDatabase()

	testJsBackendRegistry(t, reg)
}
