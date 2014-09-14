package driver

import (
	"encoding/json"
	"reflect"
	"testing"
)

// 非キャッシュ用。
func testUserRegistry(t *testing.T, reg UserRegistry) {
	attrs1, err := reg.Attributes(testUsrUuid)
	if err != nil {
		t.Fatal(err)
	} else if len(attrs1) != 0 {
		t.Error(attrs1)
	}
	attr1, err := reg.Attribute(testUsrUuid, testAttrName)
	if err != nil {
		t.Fatal(err)
	} else if attr1 != nil {
		t.Error(attr1)
	}

	if err := reg.AddAttribute(testUsrUuid, testAttrName, testAttr); err != nil {
		t.Fatal(err)
	}

	attrs2, err := reg.Attributes(testUsrUuid)
	if err != nil {
		t.Fatal(err)
	} else if attrs2 == nil {
		t.Error(attrs2)
	}
	attr2, err := reg.Attribute(testUsrUuid, testAttrName)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(testAttr, attr2) {
		// mgo で mongodb から取ってくると json の形式と違うことがあるけど、JSON 経由で同じなら許す。
		buff, _ := json.Marshal(attr2)
		var attr2_ interface{}
		json.Unmarshal(buff, &attr2_)
		if !reflect.DeepEqual(testAttr, attr2_) {
			t.Error(testAttr, attr2)
		}
	}

	if err = reg.RemoveAttribute(testUsrUuid, testAttrName); err != nil {
		t.Fatal(err)
	}

	attrs3, err := reg.Attributes(testUsrUuid)
	if err != nil {
		t.Fatal(err)
	} else if len(attrs3) != 0 {
		t.Error(attrs3)
	}
	attr3, err := reg.Attribute(testUsrUuid, testAttrName)
	if err != nil {
		t.Fatal(err)
	} else if attr3 != nil {
		t.Error(attr3)
	}
}
