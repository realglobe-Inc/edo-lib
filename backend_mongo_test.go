package driver

import (
	"testing"
)

// テストするなら、ローカルにデフォルトポートで mongodb をたてる必要あり。

func TestMongoJsBackendRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJsBackendRegistry(mongoAddr, "test_driver_mongo", "js")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	testJsBackendRegistry(t, reg)
}
