package util

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	ca := NewCache(func(prio1, prio2 interface{}) bool {
		return prio1.(time.Time).Before(prio2.(time.Time))
	})

	ca.Put("key", "val", time.Unix(1, 0))
	ca.Put("key2", "val2", time.Unix(2, 0))

	if val := ca.Get("key"); val != "val" {
		t.Error(val)
	}

	ca.Update("key2", time.Unix(3, 0))
	ca.CleanLesser(time.Unix(2, 0))

	if val := ca.Get("key"); val != nil {
		t.Error(val)
	} else if val := ca.Get("key2"); val != "val2" {
		t.Error(val)
	}
}

func TestCacheManyElements(t *testing.T) {
	ca := NewCache(func(prio1, prio2 interface{}) bool {
		return prio1.(time.Time).Before(prio2.(time.Time))
	})

	prios := []time.Time{}
	for i := 0; i < 100; i++ {
		prios = append(prios, time.Unix(rand.Int63n(3000*365*24*60*60), 0))
		ca.Put("key"+strconv.Itoa(i), "val"+strconv.Itoa(i), prios[i])
	}

	thres := time.Unix(rand.Int63(), 0)
	ca.CleanLesser(thres)

	for i := 0; i < 100; i++ {
		val := ca.Get("key" + strconv.Itoa(i))
		if val == nil {
			continue
		}
		if prios[i].Before(thres) {
			t.Error(i, prios[i], thres)
		}
	}
}
