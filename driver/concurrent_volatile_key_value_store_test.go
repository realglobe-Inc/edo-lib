// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package driver

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func testConcurrentVolatileKeyValueStore(t *testing.T, drv ConcurrentVolatileKeyValueStore) {
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
	if _, err := drv.Put(testKey, testVal, time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}

	// エントリが異なれば入れられない。
	delay := time.Minute
	if v, _, err := drv.GetAndSetEntry(testKey, nil, testKey+":entry", "0", time.Now().Add(delay)); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, testVal) {
		if !jsonEqual(v, testVal) {
			t.Error(v)
		}
	}

	exp := time.Now().Add(expiDur)
	if ok, _, err := drv.PutIfEntered(testKey, testVal, exp, testKey+":entry", "1"); err != nil {
		t.Fatal(err)
	} else if ok {
		t.Fatal("")
	}

	bef := time.Now()
	if ok, _, err := drv.PutIfEntered(testKey, testVal, exp, testKey+":entry", "0"); err != nil {
		t.Fatal(err)
	} else if !ok {
		t.Fatal("")
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

func testConcurrentVolatileKeyValueStoreConsistency(t *testing.T, drv ConcurrentVolatileKeyValueStore) {
	defer drv.Close()

	const n = 5
	const loop = 1000
	const expiDur = time.Minute

	if _, err := drv.Put(testKey, float64(0), time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	}
	defer drv.Remove(testKey)

	const eKey = testKey + ":entry"
	defer drv.Remove(eKey)

	resCh := make(chan int, n)
	errCh := make(chan error, n)
	for i := 0; i < n; i++ {
		go func(id int) {
			eVal := strconv.Itoa(id)
			res := 0
			defer func() { resCh <- res }()
			for j := 0; j < loop; j++ {
				v, _, err := drv.GetAndSetEntry(testKey, nil, eKey, eVal, time.Now().Add(expiDur))
				if err != nil {
					errCh <- err
					return
				}
				a, _ := v.(float64)
				a += 1
				ok, _, err := drv.PutIfEntered(testKey, a, time.Now().Add(expiDur), eKey, eVal)
				if err != nil {
					errCh <- err
					return
				}
				if ok {
					res++
				}
			}
		}(i)
	}

	res := 0
	for i := 0; i < n; i++ {
		res += <-resCh
	}

	v, _, err := drv.Get(testKey, nil)
	if err != nil {
		errCh <- err
		return
	}
	a, _ := v.(float64)
	if int(a) != res {
		t.Error(a, res)
	}

	for ok := false; !ok; {
		select {
		case err := <-errCh:
			t.Error(err)
		default:
			ok = true
		}
	}
}
