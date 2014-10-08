package driver

import (
	"reflect"
	"testing"
	"time"
)

var testData = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func testRawDataStore(t *testing.T, reg RawDataStore) {
	// まだ無い。
	data1, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if data1 != nil {
		t.Error(data1)
	}

	// 入れる。
	if _, err := reg.Put(testKey, testData); err != nil {
		t.Fatal(err)
	}

	// ある。
	data2, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if data2 == nil || !reflect.DeepEqual(data2, testData) {
		t.Error(data2)
	}

	// 消す。
	if err := reg.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	data3, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if data3 != nil {
		t.Error(data3)
	}
}

func testRawDataStoreStamp(t *testing.T, reg RawDataStore) {
	// まだ無い。
	data1, stmp1, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if data1 != nil || stmp1 != nil {
		t.Error(data1, stmp1)
	}

	// 入れる。
	stmp2, err := reg.Put(testKey, testData)
	if err != nil {
		t.Fatal(err)
	}

	// ある。
	data3, stmp3, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if data3 == nil || !reflect.DeepEqual(data3, testData) || stmp3 == nil {
		t.Error(data3, stmp3)
	}

	// キャッシュと同じだから返らない。
	data4, stmp4, err := reg.Get(testKey, stmp2)
	if err != nil {
		t.Fatal(err)
	} else if data4 != nil || stmp4 == nil {
		t.Error(data4, stmp4)
	}

	// キャッシュが古いから返る。
	data5, stmp5, err := reg.Get(testKey, &Stamp{Date: stmp2.Date.Add(-time.Second), Digest: stmp2.Digest})
	if err != nil {
		t.Fatal(err)
	} else if data5 == nil || !reflect.DeepEqual(data5, testData) || stmp5 == nil {
		t.Error(data5, stmp5)
	}

	// ダイジェストが違うから返る。
	data6, stmp6, err := reg.Get(testKey, &Stamp{Date: stmp2.Date, Digest: stmp2.Digest + "a"})
	if err != nil {
		t.Fatal(err)
	} else if data6 == nil || !reflect.DeepEqual(data6, testData) || stmp6 == nil {
		t.Error(data6, stmp6)
	}

	// 消す。
	if err := reg.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	data7, stmp7, err := reg.Get(testKey, stmp2)
	if err != nil {
		t.Fatal(err)
	} else if data7 != nil || stmp7 != nil {
		t.Error(data7, stmp7)
	}
}
