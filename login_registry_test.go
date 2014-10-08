package driver

import (
	"testing"
)

// 要事前登録。

func testLoginRegistry(t *testing.T, reg LoginRegistry) {
	usrUuid, _, err := reg.User(testAccToken, nil)
	if err != nil {
		t.Fatal(err)
	} else if usrUuid != testUsrName {
		t.Error(usrUuid)
	}
}
