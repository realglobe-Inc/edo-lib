package driver

import (
	"reflect"
	"testing"
	"time"
)

const testJobId = "job-no-id"

var testJobRes = &JobResult{Status: 200, Body: "job-no-result"}

func testJobRegistry(t *testing.T, reg JobRegistry) {
	expiDur := 10 * time.Millisecond

	res1, _, err := reg.Result(testJobId, nil)
	if err != nil {
		t.Fatal(err)
	} else if res1 != nil {
		t.Error(res1)
	}

	if _, err := reg.AddResult(testJobId, testJobRes, time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	}

	res2, _, err := reg.Result(testJobId, nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(res2, testJobRes) {
		if !jsonEqual(res2, testJobRes) {
			t.Error(res2)
		}
	}
}

func testJobRegistryRemoveOld(t *testing.T, reg JobRegistry) {
	expiDur := 10 * time.Millisecond

	if _, err := reg.AddResult(testJobId+"1", testJobRes, time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	}
	if _, err := reg.AddResult(testJobId+"2", testJobRes, time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * expiDur)

	if _, err := reg.AddResult(testJobId+"3", testJobRes, time.Now().Add(-expiDur)); err != nil {
		t.Fatal(err)
	}

	res1, _, err := reg.Result(testJobId+"1", nil)
	if err != nil {
		t.Fatal(err)
	} else if res1 != nil {
		t.Error(res1)
	}
	res2, _, err := reg.Result(testJobId+"2", nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(res2, testJobRes) {
		if !jsonEqual(res2, testJobRes) {
			t.Error(res2)
		}
	}

}
