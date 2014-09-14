package driver

import (
	"reflect"
	"testing"
	"time"
)

const testJobId = "job-no-id"

var testJobRes = &JobResult{Status: 200, Body: "job-no-result"}

// 非キャッシュ用。
func testJobRegistry(t *testing.T, reg JobRegistry) {
	deadline := time.Now().Add(time.Second)

	res1, err := reg.Result(testJobId)
	if err != nil {
		t.Fatal(err)
	} else if res1 != nil {
		t.Error(res1)
	}

	if err := reg.AddResult(testJobId, testJobRes, deadline); err != nil {
		t.Fatal(err)
	}

	res2, err := reg.Result(testJobId)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(res2, testJobRes) {
		t.Error(res2)
	}
}

func testJobRegistryRemoveOld(t *testing.T, reg JobRegistry) {
	deadline := time.Now().Add(50 * time.Millisecond)

	if err := reg.AddResult(testJobId+"1", testJobRes, deadline); err != nil {
		t.Fatal(err)
	}
	if err := reg.AddResult(testJobId+"2", testJobRes, deadline.Add(time.Second)); err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	if err := reg.AddResult(testJobId+"3", testJobRes, deadline); err != nil {
		t.Fatal(err)
	}

	res1, err := reg.Result(testJobId + "1")
	if err != nil {
		t.Fatal(err)
	} else if res1 != nil {
		t.Error(res1)
	}
	res2, err := reg.Result(testJobId + "2")
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(res2, testJobRes) {
		t.Error(res2)
	}

}
