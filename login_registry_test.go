package driver

import (
	"testing"
)

// 非キャッシュ用。
// 事前に、
// abc-012 に a_b-c、
// を登録しとく。
func testLoginRegistry(t *testing.T, reg LoginRegistry) {
	usrUuid, err := reg.User(testAccToken)
	if err != nil {
		t.Fatal(err)
	} else if usrUuid != testUsrName {
		t.Error(usrUuid)
	}
}
