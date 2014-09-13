package driver

import (
	"testing"
)

// 事前に、ユーザー名 a_b-c、ユーザー UUID aaaa-bbbb-cccc で登録しとく。

// 非キャッシュ用。
func testUserNameIndex(t *testing.T, reg UserNameIndex) {
	usrUuid1, err := reg.UserUuid("a_b-c")
	if err != nil {
		t.Fatal(err)
	} else if usrUuid1 != "aaaa-bbbb-cccc" {
		t.Error(usrUuid1)
	}

	usrUuid2, err := reg.UserUuid("a_b-c_d")
	if err != nil {
		t.Fatal(err)
	} else if usrUuid2 != "" {
		t.Error(usrUuid2)
	}
}
