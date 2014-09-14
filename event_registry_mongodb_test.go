package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoEventRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoEventRegistry(mongoAddr, testLabel, "event-registry")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoDriver).DB(testLabel).DropDatabase()

	testEventRegistry(t, reg)
}
