package driver

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func testConcurrentVolatileKeyValueStore(t *testing.T, reg ConcurrentVolatileKeyValueStore) {
	expiDur := 10 * time.Millisecond

	// まだ無い。
	if v, _, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Error(v)
	}

	// 入れる。
	if _, err := reg.Put(testKey, testVal, time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	}

	// ある。
	if v, _, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, testVal) {
		if !jsonEqual(v, testVal) {
			t.Error(v)
		}
	}

	// 消す。
	if err := reg.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	if v, _, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Error(v)
	}

	// また入れる。
	if _, err := reg.Put(testKey, testVal, time.Now().Add(2*expiDur)); err != nil {
		t.Fatal(err)
	}

	// エントリが異なれば入れられない。
	if v, _, err := reg.GetAndSetEntry(testKey, nil, testKey+":entry", "0", time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, testVal) {
		if !jsonEqual(v, testVal) {
			t.Error(v)
		}
	} else if ok, _, err := reg.PutIfEntered(testKey, testVal, time.Now().Add(expiDur), testKey+":entry", "1"); err != nil {
		t.Fatal(err)
	} else if ok {
		t.Fatal("")
	} else if ok, _, err := reg.PutIfEntered(testKey, testVal, time.Now().Add(expiDur), testKey+":entry", "0"); err != nil {
		t.Fatal(err)
	} else if !ok {
		t.Fatal("")
	}

	// 消えるまで待つ。
	time.Sleep(2 * expiDur)

	// もう無い。
	if v, _, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Error(v)
	}
}

func testConcurrentVolatileKeyValueStoreConsistency(t *testing.T, reg ConcurrentVolatileKeyValueStore) {
	const n = 5
	const loop = 1000
	const expiDur = time.Second

	if _, err := reg.Put(testKey, float64(0), time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	}

	const eKey = testKey + ":entry"
	resCh := make(chan int, n)
	errCh := make(chan error, n)
	for i := 0; i < n; i++ {
		go func(id int) {
			eVal := strconv.Itoa(id)
			res := 0
			defer func() { resCh <- res }()
			for j := 0; j < loop; j++ {
				v, _, err := reg.GetAndSetEntry(testKey, nil, eKey, eVal, time.Now().Add(expiDur))
				if err != nil {
					errCh <- err
					return
				}
				a, _ := v.(float64)
				a += 1
				ok, _, err := reg.PutIfEntered(testKey, a, time.Now().Add(expiDur), eKey, eVal)
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

	v, _, err := reg.Get(testKey, nil)
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
