package driver

import (
	"reflect"
	"testing"
)

// 要事前登録。

// 非キャッシュ用。
func testUserAttributeRegistry(t *testing.T, reg UserAttributeRegistry) {
	usrUuid := testUsrUuid
	attrName := testAttrName
	usrAttr := testAttr

	usrAttr1, err := reg.UserAttribute(usrUuid, attrName)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(usrAttr1, usrAttr) {
		t.Error(usrAttr1)
	}

	usrAttr2, err := reg.UserAttribute(usrUuid, attrName+"1")
	if err != nil {
		t.Fatal(err)
	} else if usrAttr2 != nil {
		t.Error(usrAttr2)
	}

	usrAttr3, err := reg.UserAttribute(usrUuid+"1", attrName)
	if err != nil {
		t.Fatal(err)
	} else if usrAttr3 != nil {
		t.Error(usrAttr3)
	}
}
