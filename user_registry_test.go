package driver

import (
	"encoding/json"
	"reflect"
	"testing"
)

func testUserRegistry(t *testing.T, reg UserRegistry) {
	attr1, _, err := reg.Attribute(testUsrUuid, testAttrName, nil)
	if err != nil {
		t.Fatal(err)
	} else if attr1 != nil {
		t.Error(attr1)
	}

	if _, err := reg.AddAttribute(testUsrUuid, testAttrName, testAttr); err != nil {
		t.Fatal(err)
	}

	attr2, _, err := reg.Attribute(testUsrUuid, testAttrName, nil)
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

	attr3, _, err := reg.Attribute(testUsrUuid, testAttrName, nil)
	if err != nil {
		t.Fatal(err)
	} else if attr3 != nil {
		t.Error(attr3)
	}
}
