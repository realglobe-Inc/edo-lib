package driver

import (
	"testing"
	"time"
)

func TestCachingRawDataStore(t *testing.T) {
	testRawDataStore(t, newCachingRawDataStore(newMemoryRawDataStore(0, 0)))
}

func TestCachingRawDataStoreStamp(t *testing.T) {
	testRawDataStoreStamp(t, newCachingRawDataStore(newMemoryRawDataStore(0, 0)))
}

func TestCachingRawDataStoreExpiration(t *testing.T) {
	staleDur := 10 * time.Millisecond
	expiDur := 50 * time.Millisecond
	reg := newCachingRawDataStore(newMemoryRawDataStore(staleDur, expiDur))

	// 入れる。
	if _, err := reg.Put(testKey, testData); err != nil {
		t.Fatal(err)
	}
	if _, err := reg.Put(testKey+"a", testData); err != nil {
		t.Fatal(err)
	}

	end := time.Now().Add(2 * expiDur)
	var caKeys map[string]bool
	var caStmp *Stamp
	for time.Now().Before(end) {
		keys, newCaStmp, err := reg.Keys(caStmp)
		if err != nil {
			t.Fatal(err)
		}
		if keys != nil {
			caKeys = keys
		}
		caStmp = newCaStmp
		if len(caKeys) != 2 {
			t.Error(caKeys)
		}
		time.Sleep(staleDur / 3)
	}
}
