package driver

import (
	"testing"
)

// テストするなら、ローカルにデフォルトポートで mongodb をたてる必要あり。

func _TestMongoJsBackendRegistry(t *testing.T) {
	reg, err := NewMongoJsBackendRegistry(mongoAddr, "test_driver_mongo", "js")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	testJsBackendRegistry(t, reg)
}
