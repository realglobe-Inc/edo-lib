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

	ca.Put("key", "value", time.Unix(1, 0))
	ca.Put("key2", "value2", time.Unix(2, 0))

	if value, _ := ca.Get("key"); value != "value" {
		t.Error(value)
	}

	ca.Update("key2", time.Unix(3, 0))
	ca.CleanLower(time.Unix(2, 0))

	if value, _ := ca.Get("key"); value != nil {
		t.Error(value)
	} else if value, _ := ca.Get("key2"); value != "value2" {
		t.Error(value)
	}
}

func TestCacheManyElements(t *testing.T) {
	ca := NewCache(func(prio1, prio2 interface{}) bool {
		return prio1.(time.Time).Before(prio2.(time.Time))
	})

	prios := []time.Time{}
	for i := 0; i < 100; i++ {
		prios = append(prios, time.Unix(rand.Int63n(3000*365*24*60*60), 0))
		ca.Put("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), prios[i])
	}

	thres := time.Unix(rand.Int63(), 0)
	ca.CleanLower(thres)

	for i := 0; i < 100; i++ {
		value, _ := ca.Get("key" + strconv.Itoa(i))
		if value == nil {
			continue
		}
		if prios[i].Before(thres) {
			t.Error(i, prios[i], thres)
		}
	}
}
