package cache

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	ca := New(func(prio1, prio2 interface{}) bool {
		return prio1.(time.Time).Before(prio2.(time.Time))
	})

	ca.Put("key", "val", time.Unix(1, 0))
	ca.Put("key2", "val2", time.Unix(2, 0))

	if val, _ := ca.Get("key"); val != "val" {
		t.Error(val)
	}

	ca.Update("key2", time.Unix(3, 0))
	ca.CleanLower(time.Unix(2, 0))

	if val, _ := ca.Get("key"); val != nil {
		t.Error(val)
	} else if val, _ := ca.Get("key2"); val != "val2" {
		t.Error(val)
	}
}

func TestCacheManyElements(t *testing.T) {
	ca := New(func(prio1, prio2 interface{}) bool {
		return prio1.(time.Time).Before(prio2.(time.Time))
	})

	for j := 0; j < 100; j++ {
		prios := []time.Time{}
		for i := 0; i < 100; i++ {
			prio := time.Unix(rand.Int63n(3000*365*24*60*60), 0)
			prios = append(prios, prio)
			ca.Put("key"+strconv.Itoa(i), "val"+strconv.Itoa(i), prio)
		}

		thres := time.Unix(rand.Int63n(3000*365*24*60*60), 0)
		ca.CleanLower(thres)

		for i := 0; i < 100; i++ {
			val, prio := ca.Get("key" + strconv.Itoa(i))
			if !prios[i].Before(thres) {
				if val == nil {
					t.Error(i, val, prio, thres)
				} else if !prio.(time.Time).Equal(prios[i]) {
					t.Error(i, val, prio, thres)
				}
			} else {
				if val != nil {
					t.Error(i, val, prio, thres)
				}
			}
		}
	}
}

func TestCacheSameKeys(t *testing.T) {
	ca := New(func(prio1, prio2 interface{}) bool {
		return prio1.(int64) < prio2.(int64)
	})

	jMax := 100
	for j := 0; j <= jMax; j++ {
		for i := 0; i < 100; i++ {
			ca.Put("key"+strconv.Itoa(i), "val"+strconv.Itoa(i*j), rand.Int63())
			if len(ca.(*cache).prioQueue) > 100 || len(ca.(*cache).keyToIdx) > 100 {
				t.Error(len(ca.(*cache).prioQueue), len(ca.(*cache).keyToIdx))
			}
		}
	}

	for i := 0; i < 100; i++ {
		val, _ := ca.Get("key" + strconv.Itoa(i))
		if val == nil {
			t.Error(i, val)
		} else if val != "val"+strconv.Itoa(i*jMax) {
			t.Error(i, val)
		}
	}
}
