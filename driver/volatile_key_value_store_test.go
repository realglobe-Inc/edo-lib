package driver

import (
	"reflect"
	"testing"
	"time"
)

func testVolatileKeyValueStore(t *testing.T, drv VolatileKeyValueStore) {
	defer drv.Close()

	expiDur := 10 * time.Millisecond

	// まだ無い。
	if v, _, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Error(v)
	}

	// 入れる。
	if _, err := drv.Put(testKey, testVal, time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}

	// ある。
	if v, _, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, testVal) {
		if !jsonEqual(v, testVal) {
			t.Error(v)
		}
	}

	// 消す。
	if err := drv.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	if v, _, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Error(v)
	}

	// また入れる。
	exp := time.Now().Add(expiDur)
	bef := time.Now()
	if _, err := drv.Put(testKey, testVal, exp); err != nil {
		t.Fatal(err)
	}
	diff := int64(time.Since(bef) / time.Nanosecond)

	// 消えるかどうか。
	for {
		bef := time.Now()
		v, _, err := drv.Get(testKey, nil)
		if err != nil {
			t.Fatal(err)
		}
		aft := time.Now()

		// GC 等で時間が掛かることもあるため、aft > exp でも nil が返るとは限らない。
		// だが、aft <= exp であれば非 nil が返らなければならない。
		// 同様に、bef > exp であれば nil が返らなければならない。

		if aft.UnixNano() <= cutOff(exp.UnixNano(), 1e6)-diff { // redis の粒度がミリ秒のため。
			if v == nil {
				t.Error(aft)
				t.Error(exp)
				return
			}
		} else if bef.UnixNano() > cutOff(exp.UnixNano(), 1e6)+1e6+diff { // redis の粒度がミリ秒のため。
			if v != nil {
				t.Error(bef)
				t.Error(exp)
				return
			}
			// 消えた。
			return
		} else if v == nil { // bef <= exp < aft
			// 消えた。
			return
		}

		time.Sleep(time.Millisecond)
	}
}
