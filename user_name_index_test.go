package driver

import (
	"testing"
)

// 非キャッシュ用。
func testUserNameIndex(t *testing.T, reg UserNameIndex) {
	usrUuid1, err := reg.UserUuid(testUsrName)
	if err != nil {
		t.Fatal(err)
	} else if usrUuid1 != testUsrUuid {
		t.Error(usrUuid1)
	}

	usrUuid2, err := reg.UserUuid(testUsrName + "_d")
	if err != nil {
		t.Fatal(err)
	} else if usrUuid2 != "" {
		t.Error(usrUuid2)
	}
}
