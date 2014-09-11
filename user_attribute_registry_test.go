package driver

import (
	"testing"
)

// 事前に、ユーザー UUID a_b-c、属性名 attribute、属性値 abcd で登録しとく。

func testUserAttributeRegistry(t *testing.T, reg UserAttributeRegistry) {
	usrUuid := "a_b-c"
	attrName := "attribute"
	attr := "abcd"

	usrUuid1, err := reg.UserAttribute(usrUuid, attrName)
	if err != nil {
		t.Fatal(err)
	} else if usrUuid1 != attr {
		t.Error(usrUuid1)
	}

	usrUuid2, err := reg.UserAttribute(usrUuid, attrName+"1")
	if err != nil {
		t.Fatal(err)
	} else if usrUuid2 != nil && usrUuid2 != "" {
		t.Error(usrUuid2)
	}
}
