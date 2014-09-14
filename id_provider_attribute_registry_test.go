package driver

import (
	"reflect"
	"testing"
	"time"
)

// 要事前登録。

// 非キャッシュ用。
func testIdProviderAttributeRegistry(t *testing.T, reg IdProviderAttributeRegistry) {
	idpUuid := testIdpUuid
	attrName := testAttrName
	idpAttr := testAttr

	idpAttr1, err := reg.IdProviderAttribute(idpUuid, attrName)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(idpAttr1, idpAttr) {
		t.Error(idpAttr1)
	}

	idpAttr2, err := reg.IdProviderAttribute(idpUuid, attrName+"1")
	if err != nil {
		t.Fatal(err)
	} else if idpAttr2 != nil {
		t.Error(idpAttr2)
	}

	idpAttr3, err := reg.IdProviderAttribute(idpUuid+"1", attrName)
	if err != nil {
		t.Fatal(err)
	} else if idpAttr3 != nil {
		t.Error(idpAttr3)
	}
}

// キャッシュ用。
func testDatedIdProviderAttributeRegistry(t *testing.T, reg DatedIdProviderAttributeRegistry) {
	idpUuid := testIdpUuid
	attrName := testAttrName
	idpAttr := testAttr

	idpAttr1, stmp1, err := reg.StampedIdProviderAttribute(idpUuid, attrName, nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(idpAttr1, idpAttr) || stmp1 == nil {
		t.Error(idpAttr1, stmp1)
	}

	idpAttr2, stmp2, err := reg.StampedIdProviderAttribute(idpUuid, attrName+"1", nil)
	if err != nil {
		t.Fatal(err)
	} else if idpAttr2 != nil || stmp2 != nil {
		t.Error(idpAttr2, stmp2)
	}

	idpAttr3, stmp3, err := reg.StampedIdProviderAttribute(idpUuid+"1", attrName, nil)
	if err != nil {
		t.Fatal(err)
	} else if idpAttr3 != nil || stmp3 != nil {
		t.Error(idpAttr3, stmp3)
	}

	// キャッシュと同じだから返らない。
	idpAttr4, stmp4, err := reg.StampedIdProviderAttribute(idpUuid, attrName, stmp1)
	if err != nil {
		t.Fatal(err)
	} else if idpAttr4 != nil || stmp4 == nil {
		t.Error(idpAttr4, stmp4)
	}

	// キャッシュが古いから返る。
	idpAttr5, stmp5, err := reg.StampedIdProviderAttribute(idpUuid, attrName, &Stamp{Date: stmp1.Date.Add(-time.Second), Digest: stmp1.Digest})
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(idpAttr5, idpAttr) || stmp5 == nil {
		t.Error(idpAttr5, stmp5)
	}

	// ダイジェストが違うから返る。
	idpAttr6, stmp6, err := reg.StampedIdProviderAttribute(idpUuid, attrName, &Stamp{Date: stmp1.Date, Digest: stmp1.Digest + "a"})
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(idpAttr1, idpAttr) || stmp6 == nil {
		t.Error(idpAttr6, stmp6)
	}
}
