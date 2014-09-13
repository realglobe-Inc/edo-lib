package driver

import (
	"reflect"
	"testing"
	"time"
)

// 非キャッシュ用。
func testJobRegistry(t *testing.T, reg JobRegistry) {
	var jobId string = "a_b-c"
	deadline := time.Now().Add(time.Second)
	res := &JobResult{Status: 200, Body: "konnnann demashita"}

	res1, err := reg.Result(jobId)
	if err != nil {
		t.Fatal(err)
	} else if res1 != nil {
		t.Error(res1)
	}

	if err := reg.AddResult(jobId, res, deadline); err != nil {
		t.Fatal(err)
	}

	res2, err := reg.Result(jobId)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(res, res2) {
		t.Error(res2)
	}
}

func testJobRegistryRemoveOld(t *testing.T, reg JobRegistry) {
	jobId := "a_b-c"
	deadline := time.Now().Add(50 * time.Millisecond)
	res := &JobResult{Status: 200, Body: "konnnann demashita"}

	if err := reg.AddResult(jobId+"1", res, deadline); err != nil {
		t.Fatal(err)
	}
	if err := reg.AddResult(jobId+"2", res, deadline.Add(time.Second)); err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	if err := reg.AddResult(jobId+"3", res, deadline); err != nil {
		t.Fatal(err)
	}

	res1, err := reg.Result(jobId + "1")
	if err != nil {
		t.Fatal(err)
	} else if res1 != nil {
		t.Error(res1)
	}
	res2, err := reg.Result(jobId + "2")
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(res, res2) {
		t.Error(res, res2)
	}

}
