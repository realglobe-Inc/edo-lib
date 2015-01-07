package util

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	type testType struct {
		D Duration
	}

	var a testType
	a.D = Duration(5*time.Hour + 3*time.Minute + 2*time.Second)

	buff, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}

	var b testType
	if err := json.Unmarshal(buff, &b); err != nil {
		t.Fatal(err)
	}

	if b != a {
		t.Error(b)
	}
}
